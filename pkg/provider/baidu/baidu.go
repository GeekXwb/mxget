package baidu

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
	APISearch       = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.search.merge&from=android&version=8.1.4.0&format=json&type=-1&isNew=1"
	APIGetSong      = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.song.getInfos&format=json&from=android&version=8.1.4.0"
	APIGetSongs     = "http://music.taihe.com/data/music/fmlink"
	APIGetSongLyric = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.song.lry&format=json&from=android&version=8.1.4.0"
	APIGetArtist    = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.artist.getSongList&from=android&version=8.1.4.0&format=json&order=2"
	APIGetAlbum     = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.album.getAlbumInfo&from=android&version=8.1.4.0&format=json"
	APIGetPlaylist  = "http://musicapi.qianqian.com/v1/restserver/ting?method=baidu.ting.ugcdiy.getBaseInfo&from=android&version=8.1.4.0"
)

var (
	std = New(request.DefaultClient)

	defaultHeaders = ghttp.Headers{
		"Origin":  "http://music.taihe.com",
		"Referer": "http://music.taihe.com",
	}
)

type (
	CommonResponse struct {
		ErrorCode    int    `json:"error_code,omitempty"`
		ErrorMessage string `json:"error_message,omitempty"`
	}

	Song struct {
		SongId     string `json:"song_id"`
		Title      string `json:"title"`
		TingUid    string `json:"ting_uid"`
		Author     string `json:"author"`
		AlbumId    string `json:"album_id"`
		AlbumTitle string `json:"album_title"`
		PicBig     string `json:"pic_big"`
		LrcLink    string `json:"lrclink"`
		CopyType   string `json:"copy_type"`
		URL        string `json:"-"`
		Lyric      string `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Result struct {
			SongInfo struct {
				SongList []*Song `json:"song_list"`
			} `json:"song_info"`
		} `json:"result"`
	}

	URL struct {
		ShowLink    string `json:"show_link"`
		FileFormat  string `json:"file_format"`
		FileBitrate int    `json:"file_bitrate"`
		FileLink    string `json:"file_link"`
	}

	SongResponse struct {
		CommonResponse
		SongInfo Song `json:"songinfo"`
		SongURL  struct {
			URL []URL `json:"url"`
		} `json:"songurl"`
	}

	SongsResponse struct {
		ErrorCode int `json:"errorCode"`
		Data      struct {
			SongList []*struct {
				SongId     int    `json:"songId"`
				SongName   string `json:"songName"`
				ArtistId   string `json:"artistId"`
				ArtistName string `json:"artistName"`
				AlbumId    int    `json:"albumId"`
				AlbumName  string `json:"albumName"`
				SongPicBig string `json:"songPicBig"`
				LrcLink    string `json:"lrcLink"`
				CopyType   int    `json:"copyType"`
				SongLink   string `json:"songLink"`
				ShowLink   string `json:"showLink"`
				Format     string `json:"format"`
				Rate       int    `json:"rate"`
			} `json:"songList"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Title      string `json:"title"`
		LrcContent string `json:"lrcContent"`
	}

	ArtistResponse struct {
		CommonResponse
		ArtistInfo struct {
			TingUid   string `json:"ting_uid"`
			Name      string `json:"name"`
			AvatarBig string `json:"avatar_big"`
		} `json:"artistinfo"`
		SongList []*Song `json:"songlist"`
	}

	// 百度的这个接口设计比较坑爹，无数据时albumInfo字段为数组类型，导致json反序列化失败
	AlbumResponse struct {
		CommonResponse
		AlbumInfo struct {
			AlbumId string `json:"album_id"`
			Title   string `json:"title"`
			PicBig  string `json:"pic_big"`
		} `json:"albumInfo"`
		SongList []*Song `json:"songlist"`
	}

	PlaylistResponse struct {
		CommonResponse
		Result struct {
			Info struct {
				ListId    string `json:"list_id"`
				ListTitle string `json:"list_title"`
				ListPic   string `json:"list_pic"`
			} `json:"info"`
			SongList []*Song `json:"songlist"`
		} `json:"result"`
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
	if c.ErrorMessage == "" {
		return strconv.Itoa(c.ErrorCode)
	}
	return c.ErrorMessage
}

func (s *SearchSongsResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongsResponse) String() string {
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

func (a *API) SendRequest(method string, url string, opts ...ghttp.RequestOption) *ghttp.Response {
	opts = append(opts,
		ghttp.WithHeaders(defaultHeaders),
	)
	return a.Client.Send(method, url, opts...)
}

func songURL(urls []URL) string {
	for _, i := range urls {
		if i.FileFormat == "mp3" {
			return i.ShowLink
		}
	}
	return ""
}

func (a *API) patchSongsURL(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			resp, err := a.GetSongRaw(ctx, s.SongId)
			if err == nil {
				s.URL = songURL(resp.SongURL.URL)
				if s.LrcLink == "" {
					s.LrcLink = resp.SongInfo.LrcLink
				}
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
			lyric, err := a.SendRequest(ghttp.MethodGet, s.LrcLink,
				ghttp.WithContext(ctx),
			).Text()
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
			Id:       s.SongId,
			Name:     strings.TrimSpace(s.Title),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Author, ",", "/")),
			Album:    strings.TrimSpace(s.AlbumTitle),
			PicUrl:   strings.SplitN(s.PicBig, "@", 2)[0],
			Lyric:    s.Lyric,
			Playable: s.URL != "",
			Url:      s.URL,
		}
	}
	return songs
}
