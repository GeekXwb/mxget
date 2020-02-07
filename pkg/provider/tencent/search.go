package tencent

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

	n := len(resp.Data.Song.List)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.SongResponse, n)
	for i, s := range resp.Data.Song.List {
		artists := make([]string, len(s.Singer))
		for j, a := range s.Singer {
			artists[j] = strings.TrimSpace(a.Name)
		}
		songs[i] = &api.SongResponse{
			Id:     s.Mid,
			Name:   strings.TrimSpace(s.Title),
			Artist: strings.Join(artists, "/"),
			Album:  strings.TrimSpace(s.Album.Name),
		}
	}
	return &api.SearchSongsResponse{
		Keyword: keyword,
		Count:   uint32(n),
		Songs:   songs,
	}, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := ghttp.Params{
		"w": keyword,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(pageSize),
	}

	resp := new(SearchSongsResponse)
	err := a.SendRequest(ghttp.MethodGet, APISearch,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("search songs: %d", resp.Code)
	}

	return resp, nil
}
