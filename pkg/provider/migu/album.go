package migu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetAlbum(ctx context.Context, albumId string) (*api.AlbumResponse, error) {
	resp, err := a.GetAlbumRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsLyric(ctx, resp.Resource[0].SongItems...)
	songs := resolve(resp.Resource[0].SongItems...)
	return &api.AlbumResponse{
		Id:     resp.Resource[0].AlbumId,
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicUrl: picURL(resp.Resource[0].ImgItems),
		Count:  uint32(len(songs)),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	params := ghttp.Params{
		"resourceId": albumId,
	}

	resp := new(AlbumResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetAlbum,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "000000" {
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
