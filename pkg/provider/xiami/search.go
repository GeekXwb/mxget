package xiami

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Data.Songs)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.SongResponse, n)
	for i, s := range resp.Data.Data.Songs {
		songs[i] = &api.SongResponse{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Singers, " / ", "/")),
			Album:  strings.TrimSpace(s.AlbumName),
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
	token, err := a.getToken(APISearch)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"key": keyword,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

	resp := new(SearchSongsResponse)
	err = a.SendRequest(ghttp.MethodGet, APISearch,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("search songs: %w", err)
	}

	return resp, nil
}
