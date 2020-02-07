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

func (a *API) GetAlbum(ctx context.Context, albumId string) (*api.AlbumResponse, error) {
	resp, err := a.GetAlbumRaw(ctx, albumId, 1, 9999)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.MusicList)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsURL(ctx, SongDefaultBR, resp.Data.MusicList...)
	a.patchSongsLyric(ctx, resp.Data.MusicList...)
	songs := resolve(resp.Data.MusicList...)
	return &api.AlbumResponse{
		Id:     strconv.Itoa(resp.Data.AlbumId),
		Name:   strings.TrimSpace(resp.Data.Album),
		PicUrl: resp.Data.Pic,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取专辑，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetAlbumRaw(ctx context.Context, albumId string, page int, pageSize int) (*AlbumResponse, error) {
	params := ghttp.Params{
		"albumId": albumId,
		"pn":      strconv.Itoa(page),
		"rn":      strconv.Itoa(pageSize),
	}
	resp := new(AlbumResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetAlbum,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
