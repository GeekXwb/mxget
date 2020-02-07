package kuwo

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
	APISearch         = "http://www.kuwo.cn/api/www/search/searchMusicBykeyWord"
	APIGetSong        = "http://www.kuwo.cn/api/www/music/musicInfo"
	APIGetSongURL     = "http://www.kuwo.cn/url?format=mp3&response=url&type=convert_url3"
	APIGetSongLyric   = "http://www.kuwo.cn/newh5/singles/songinfoandlrc"
	APIGetArtistInfo  = "http://www.kuwo.cn/api/www/artist/artist"
	APIGetArtistSongs = "http://www.kuwo.cn/api/www/artist/artistMusic"
	APIGetAlbum       = "http://www.kuwo.cn/api/www/album/albumInfo"
	APIGetPlaylist    = "http://www.kuwo.cn/api/www/playlist/playListInfo"

	SongDefaultBR = 128
)

var (
	std = New(request.DefaultClient)
)

type (
	CommonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Song struct {
		RId             int    `json:"rid"`
		Name            string `json:"name"`
		ArtistId        int    `json:"artistid"`
		Artist          string `json:"artist"`
		AlbumId         int    `json:"albumid"`
		Album           string `json:"album"`
		AlbumPic        string `json:"albumpic"`
		Track           int    `json:"track"`
		IsListenFee     bool   `json:"isListenFee"`
		SongTimeMinutes string `json:"songTimeMinutes"`
		Lyric           string `json:"-"`
		URL             string `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Total string  `json:"total"`
			List  []*Song `json:"list"`
		} `json:"data"`
	}

	SongResponse struct {
		CommonResponse
		Data Song `json:"data"`
	}

	SongURLResponse struct {
		CommonResponse
		URL string `json:"url"`
	}

	SongLyricResponse struct {
		Status int    `json:"status"`
		Msg    string `json:"msg,omitempty"`
		Data   struct {
			LrcList []struct {
				Time      string `json:"time"`
				LineLyric string `json:"lineLyric"`
			} `json:"lrclist"`
		} `json:"data"`
	}

	ArtistInfo struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Pic300 string `json:"pic300"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data ArtistInfo `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			List []*Song `json:"list"`
		} `json:"data"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			AlbumId   int     `json:"albumId"`
			Album     string  `json:"album"`
			Pic       string  `json:"pic"`
			MusicList []*Song `json:"musicList"`
		} `json:"data"`
	}

	PlaylistResponse struct {
		CommonResponse
		Data struct {
			Id        int     `json:"id"`
			Name      string  `json:"name"`
			Img700    string  `json:"img700"`
			MusicList []*Song `json:"musicList"`
		} `json:"data"`
	}

	API struct {
		Client *ghttp.Client
	}
)

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

func (s *SongResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongURLResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongLyricResponse) String() string {
	return utils.ToJSON(s, false)
}

func (a *ArtistInfoResponse) String() string {
	return utils.ToJSON(a, false)
}

func (a *ArtistSongsResponse) String() string {
	return utils.ToJSON(a, false)
}

func (a *AlbumResponse) String() string {
	return utils.ToJSON(a, false)
}

func (p *PlaylistResponse) String() string {
	return utils.ToJSON(p, false)
}

func (a *API) SendRequest(method string, url string, opts ...ghttp.RequestOption) *ghttp.Response {
	// csrf 必须跟 kw_token 保持一致
	csrf := "0"
	cookie, err := a.Client.FilterCookie(url, "kw_token")
	if err != nil {
		opts = append(opts, ghttp.WithCookies(ghttp.Cookies{
			"kw_token": csrf,
		}))
	} else {
		csrf = cookie.Value
	}

	headers := ghttp.Headers{
		"Origin":  "http://www.kuwo.cn",
		"Referer": "http://www.kuwo.cn",
		"csrf":    csrf,
	}
	opts = append(opts,
		ghttp.WithHeaders(headers),
	)
	return a.Client.Send(method, url, opts...)
}

func (a *API) patchSongsURL(ctx context.Context, br int, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(ctx, s.RId, br)
			if err == nil {
				s.URL = url
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(ctx, s.RId)
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
		songs[i] = &api.SongResponse{
			Id:       strconv.Itoa(s.RId),
			Name:     strings.TrimSpace(s.Name),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:    strings.TrimSpace(s.Album),
			PicUrl:   s.AlbumPic,
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			Url:      s.URL,
		}
	}
	return songs
}
