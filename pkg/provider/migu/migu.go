package migu

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/request"
	"github.com/winterssy/mxget/pkg/utils"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	APISearch         = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/search_all.do?isCopyright=1&isCorrect=1"
	APISearchV2       = "http://m.music.migu.cn/migu/remoting/scr_search_tag?type=2"
	APIGetSongId      = "http://music.migu.cn/v3/api/music/audioPlayer/songs?type=1"
	APIGetSong        = "https://app.c.nf.migu.cn/MIGUM2.0/v2.0/content/querySongBySongId.do?contentId=0"
	APIGetSongURL     = "https://app.c.nf.migu.cn/MIGUM2.0/v2.0/content/listen-url?copyrightId=0&netType=01&toneFlag=HQ"
	APIGetSongLyric   = "http://music.migu.cn/v3/api/music/audioPlayer/getLyric"
	APIGetSongPic     = "http://music.migu.cn/v3/api/music/audioPlayer/getSongPic"
	APIGetArtistInfo  = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2002"
	APIGetArtistSongs = "https://app.c.nf.migu.cn/MIGUM3.0/v1.0/template/singerSongs/release?templateVersion=2"
	APIGetAlbum       = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2003"
	APIGetPlaylist    = "https://app.c.nf.migu.cn/MIGUM2.0/v1.0/content/resourceinfo.do?needSimple=01&resourceType=2021"

	SongURL = "https://app.pd.nf.migu.cn/MIGUM2.0/v1.0/content/sub/listenSong.do?contentId=%s&copyrightId=0&netType=01&resourceType=%s&toneFlag=%s&channel=0"

	SongDefaultBR = 128
)

var (
	std = New(request.DefaultClient)

	codeRate = map[int]string{
		64:  "LQ",
		128: "PQ",
		320: "HQ",
		999: "SQ",
	}

	defaultHeaders = ghttp.Headers{
		"channel": "0",
		"Origin":  "http://music.migu.cn/v3",
		"Referer": "http://music.migu.cn/v3",
	}
)

type (
	CommonResponse struct {
		Code string `json:"code"`
		Info string `json:"info,omitempty"`
	}

	SearchSongsResult struct {
		ResourceType string `json:"resourceType"`
		ContentId    string `json:"contentId"`
		CopyrightId  string `json:"copyrightId"`
		Id           string `json:"id"`
		Name         string `json:"name"`
		Singers      []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"singers"`
		Albums []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"albums"`
	}

	SearchSongsResponse struct {
		CommonResponse
		SongResultData struct {
			TotalCount string               `json:"totalCount"`
			Result     []*SearchSongsResult `json:"result"`
		} `json:"songResultData"`
	}

	ImgItem struct {
		ImgSizeType string `json:"imgSizeType"`
		Img         string `json:"img"`
	}

	SongIdResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg,omitempty"`
		Items      []struct {
			SongId string `json:"songId"`
		} `json:"items"`
	}

	Song struct {
		ResourceType string    `json:"resourceType"`
		ContentId    string    `json:"contentId"`
		CopyrightId  string    `json:"copyrightId"`
		SongId       string    `json:"songId"`
		SongName     string    `json:"songName"`
		SingerId     string    `json:"singerId"`
		Singer       string    `json:"singer"`
		AlbumId      string    `json:"albumId"`
		Album        string    `json:"album"`
		AlbumImgs    []ImgItem `json:"albumImgs"`
		LrcURL       string    `json:"lrcUrl"`
		PicURL       string    `json:"-"`
		Lyric        string    `json:"-"`
		URL          string    `json:"-"`
	}

	SongResponse struct {
		CommonResponse
		Resource []*Song `json:"resource"`
	}

	SongURLResponse struct {
		CommonResponse
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg"`
		Lyric      string `json:"lyric"`
	}

	SongPicResponse struct {
		ReturnCode string `json:"returnCode"`
		Msg        string `json:"msg"`
		SmallPic   string `json:"smallPic"`
		MediumPic  string `json:"mediumPic"`
		LargePic   string `json:"largePic"`
	}

	ArtistInfo struct {
		ResourceType string    `json:"resourceType"`
		SingerId     string    `json:"singerId"`
		Singer       string    `json:"singer"`
		Imgs         []ImgItem `json:"imgs"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Resource []ArtistInfo `json:"resource"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			ContentItemList []struct {
				ItemList []struct {
					Song Song `json:"song"`
				} `json:"itemList"`
			} `json:"contentItemList"`
		} `json:"data"`
	}

	Album struct {
		ResourceType string    `json:"resourceType"`
		AlbumId      string    `json:"albumId"`
		Title        string    `json:"title"`
		ImgItems     []ImgItem `json:"imgItems"`
		SongItems    []*Song   `json:"songItems"`
	}

	AlbumResponse struct {
		CommonResponse
		Resource []Album `json:"resource"`
	}

	Playlist struct {
		ResourceType string `json:"resourceType"`
		MusicListId  string `json:"musicListId"`
		Title        string `json:"title"`
		ImgItem      struct {
			Img string `json:"img"`
		} `json:"imgItem"`
		SongItems []*Song `json:"songItems"`
	}

	PlaylistResponse struct {
		CommonResponse
		Resource []Playlist `json:"resource"`
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
	if c.Info == "" {
		return c.Code
	}
	return c.Info
}

func (s *SearchSongsResponse) String() string {
	return utils.ToJSON(s, false)
}

func (s *SongIdResponse) String() string {
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

func (s *SongPicResponse) String() string {
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
	opts = append(opts,
		ghttp.WithHeaders(defaultHeaders),
	)
	return a.Client.Send(method, url, opts...)
}

func (a *API) patchSongsInfo(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			picURL, err := a.GetSongPic(ctx, s.SongId)
			if err == nil {
				if !strings.HasPrefix(picURL, "http:") {
					picURL = "http:" + picURL
				}
				s.PicURL = picURL
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongsURL(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(ctx, s.ContentId, s.ResourceType)
			if err == nil {
				s.URL = url
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

// 部分歌词文本文件由于不是utf-8编码，可能会乱码，目前无解
func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			if s.LrcURL != "" {
				lrcBytes, err := a.SendRequest(ghttp.MethodGet, s.LrcURL,
					ghttp.WithContext(ctx),
				).Content()
				if err == nil {
					if utf8.Valid(lrcBytes) {
						s.Lyric = string(lrcBytes)
					} else {
						// GBK 编码
						lrcBytes, err = simplifiedchinese.GB18030.NewDecoder().Bytes(lrcBytes)
						if err == nil {
							s.Lyric = string(lrcBytes)
						}
					}
				}
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func picURL(imgs []ImgItem) string {
	for _, i := range imgs {
		if i.ImgSizeType == "03" {
			return i.Img
		}
	}
	return ""
}

func songURL(contentId string, br int) string {
	var _br int
	switch br {
	case 64, 128, 320, 999:
		_br = br
	default:
		_br = 320
	}
	return fmt.Sprintf(SongURL, contentId, "E", codeRate[_br])
}

func resolve(src ...*Song) []*api.SongResponse {
	songs := make([]*api.SongResponse, len(src))
	for i, s := range src {
		url := songURL(s.ContentId, SongDefaultBR)
		songs[i] = &api.SongResponse{
			Id:       s.SongId,
			Name:     strings.TrimSpace(s.SongName),
			Artist:   strings.TrimSpace(strings.ReplaceAll(s.Singer, "|", "/")),
			Album:    strings.TrimSpace(s.Album),
			PicUrl:   picURL(s.AlbumImgs),
			Lyric:    s.Lyric,
			Playable: url != "",
			Url:      url,
		}
	}
	return songs
}
