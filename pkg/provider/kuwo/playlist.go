package kuwo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.PlaylistResponse, error) {
	resp, err := a.GetPlaylistRaw(ctx, playlistId, 1, 9999)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.MusicList)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	a.patchSongsURL(ctx, SongDefaultBR, resp.Data.MusicList...)
	a.patchSongsLyric(ctx, resp.Data.MusicList...)
	songs := resolve(resp.Data.MusicList...)
	return &api.PlaylistResponse{
		Id:     strconv.Itoa(resp.Data.Id),
		Name:   strings.TrimSpace(resp.Data.Name),
		PicUrl: resp.Data.Img700,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌单，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistResponse, error) {
	params := ghttp.Params{
		"pid": playlistId,
		"pn":  strconv.Itoa(page),
		"rn":  strconv.Itoa(pageSize),
	}

	resp := new(PlaylistResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetPlaylist,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
	}

	return resp, nil
}
