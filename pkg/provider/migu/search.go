package migu

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	musics := resp.GetArray("musics")
	n := len(musics)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.SongResponse, n)
	for i, s := range musics {
		obj := s.Object()
		songs[i] = &api.SongResponse{
			Id:     obj.GetString("id"),
			Name:   obj.GetString("songName"),
			Artist: strings.ReplaceAll(obj.GetString("singerName"), ", ", "/"),
			Album:  obj.GetString("albumName"),
		}
	}
	return &api.SearchSongsResponse{
		Keyword: keyword,
		Count:   uint32(n),
		Songs:   songs,
	}, nil
}

// 搜索歌曲
// func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
// 	switchOption := map[string]int{
// 		"song":     1,
// 		"album":    0,
// 		"singer":   0,
// 		"tagSong":  0,
// 		"mvSong":   0,
// 		"songlist": 0,
// 		"bestShow": 0,
// 	}
// 	enc, _ := json.Marshal(switchOption)
// 	params := ghttp.Params{
// 		"searchSwitch": string(enc),
// 		"text":         keyword,
// 		"pageNo":       strconv.Itoa(page),
// 		"pageSize":     strconv.Itoa(pageSize),
// 	}
//
// 	resp := new(SearchSongsResponse)
// 	err := a.SendRequest(ghttp.MethodGet, APISearch,
// 		ghttp.WithQuery(params),
// 		ghttp.WithContext(ctx),
// 	).JSON(resp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.Code != "000000" {
// 		return nil, fmt.Errorf("search songs: %s", resp.errorMessage())
// 	}
//
// 	return resp, nil
// }

func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (ghttp.H, error) {
	params := ghttp.Params{
		"keyword": keyword,
		"pgc":     strconv.Itoa(page),
		"rows":    strconv.Itoa(pageSize),
	}

	data, err := a.SendRequest(ghttp.MethodGet, APISearchV2,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).H()
	if err != nil {
		return nil, err
	}

	if !data.GetBoolean("success") {
		return data, fmt.Errorf("search songs: %s", data.GetString("msg"))
	}

	return data, nil
}
