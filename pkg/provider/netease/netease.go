package netease

import (
	"context"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/request"
	"github.com/winterssy/mxget/pkg/utils"
)

const (
	APILinux          = "https://music.163.com/api/linux/forward"
	APISearch         = "https://music.163.com/weapi/search/get"
	APIGetSongs       = "https://music.163.com/weapi/v3/song/detail"
	APIGetSongsURL    = "https://music.163.com/weapi/song/enhance/player/url"
	APIGetArtist      = "https://music.163.com/weapi/v1/artist/%d"
	APIGetAlbum       = "https://music.163.com/weapi/v1/album/%d"
	APIGetPlaylist    = "https://music.163.com/weapi/v3/playlist/detail"
	APIEmailLogin     = "https://music.163.com/weapi/login"
	APICellphoneLogin = "https://music.163.com/weapi/login/cellphone"
	APIRefreshLogin   = "https://music.163.com/weapi/login/token/refresh"
	APILogout         = "https://music.163.com/weapi/logout"

	SongRequestLimit = 1000
	SongDefaultBR    = 128
)

var (
	std = New(request.DefaultClient)

	cookies ghttp.Cookies

	defaultHeaders = ghttp.Headers{
		"Origin":  "https://music.163.com",
		"Referer": "https://music.163.com",
	}
)

type (
	CommonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Song struct {
		Id      int      `json:"id"`
		Name    string   `json:"name"`
		Artists []Artist `json:"ar"`
		Album   Album    `json:"al"`
		Track   int      `json:"no"`
		Lyric   string   `json:"-"`
		URL     string   `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Result struct {
			Songs []*struct {
				Id      int    `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"album"`
			} `json:"songs"`
			SongCount int `json:"songCount"`
		} `json:"result"`
	}

	SongURL struct {
		Id   int    `json:"id"`
		URL  string `json:"url"`
		BR   int    `json:"br"`
		Code int    `json:"code"`
	}

	SongsResponse struct {
		CommonResponse
		Songs []*Song `json:"songs"`
	}

	SongURLResponse struct {
		CommonResponse
		Data []struct {
			Code int    `json:"code"`
			Id   int    `json:"id"`
			BR   int    `json:"br"`
			URL  string `json:"url"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Lrc struct {
			Lyric string `json:"lyric"`
		} `json:"lrc"`
		TLyric struct {
			Lyric string `json:"lyric"`
		} `json:"tlyric"`
	}

	Artist struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	ArtistResponse struct {
		CommonResponse
		Artist struct {
			Artist
		} `json:"artist"`
		HotSongs []*Song `json:"hotSongs"`
	}

	Album struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	AlbumResponse struct {
		CommonResponse
		Album Album   `json:"album"`
		Songs []*Song `json:"songs"`
	}

	PlaylistResponse struct {
		CommonResponse
		Playlist struct {
			Id       int     `json:"id"`
			Name     string  `json:"name"`
			PicURL   string  `json:"coverImgUrl"`
			Tracks   []*Song `json:"tracks"`
			TrackIds []struct {
				Id int `json:"id"`
			} `json:"trackIds"`
			Total int `json:"trackCount"`
		} `json:"playlist"`
	}

	LoginResponse struct {
		CommonResponse
		LoginType int `json:"loginType"`
		Account   struct {
			Id       int    `json:"id"`
			UserName string `json:"userName"`
		} `json:"account"`
	}

	API struct {
		Client *ghttp.Client
	}
)

func init() {
	cookies, _ = createCookie()
}

func New(client *ghttp.Client) *API {
	return &API{
		Client: client,
	}
}

func Client() *API {
	return std
}

func (c *CommonResponse) errorMessage() string {
	if c.Msg == "" {
		return strconv.Itoa(c.Code)
	}
	return c.Msg
}

func (s *SearchSongsResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongsResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongURLResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongLyricResponse) String() string {
	return utils.ToJSON(s, false)
}

func (a *ArtistResponse) String() string {
	return utils.ToJSON(a, false)
}

func (a *AlbumResponse) String() string {
	return utils.ToJSON(a, false)
}

func (p *PlaylistResponse) String() string {
	return utils.ToJSON(p, false)
}

func (e *LoginResponse) String() string {
	return utils.ToJSON(e, false)
}

func (a *API) SendRequest(method string, url string, opts ...ghttp.RequestOption) *ghttp.Response {
	opts = append(opts,
		ghttp.WithHeaders(defaultHeaders),
	)

	// 如果已经登录，不需要额外设置cookies，cookie jar会自动管理
	_, err := a.Client.FilterCookie(url, "MUSIC_U")
	if err != nil {
		opts = append(opts, ghttp.WithCookies(cookies))
	}

	return a.Client.Send(method, url, opts...)
}

func (a *API) patchSongsURL(ctx context.Context, br int, songs ...*Song) {
	ids := make([]int, len(songs))
	for i, s := range songs {
		ids[i] = s.Id
	}

	resp, err := a.GetSongsURLRaw(ctx, br, ids...)
	if err == nil && len(resp.Data) != 0 {
		m := make(map[int]string, len(resp.Data))
		for _, i := range resp.Data {
			m[i.Id] = i.URL
		}
		for _, s := range songs {
			s.URL = m[s.Id]
		}
	}
}

func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(ctx, s.Id)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func resolve(src ...*Song) []*api.SongResponse {
	songs := make([]*api.SongResponse, len(src))
	for i, s := range src {
		artists := make([]string, len(s.Artists))
		for j, a := range s.Artists {
			artists[j] = strings.TrimSpace(a.Name)
		}
		songs[i] = &api.SongResponse{
			Id:       strconv.Itoa(s.Id),
			Name:     strings.TrimSpace(s.Name),
			Artist:   strings.Join(artists, "/"),
			Album:    strings.TrimSpace(s.Album.Name),
			PicUrl:   s.Album.PicURL,
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			Url:      s.URL,
		}
	}
	return songs
}
