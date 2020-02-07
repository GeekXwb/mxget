package baidu

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

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsURL(ctx, resp.SongList...)
	a.patchSongsLyric(ctx, resp.SongList...)
	songs := resolve(resp.SongList...)
	return &api.AlbumResponse{
		Id:     resp.AlbumInfo.AlbumId,
		Name:   strings.TrimSpace(resp.AlbumInfo.Title),
		PicUrl: strings.SplitN(resp.AlbumInfo.PicBig, "@", 2)[0],
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	params := ghttp.Params{
		"album_id": albumId,
	}

	resp := new(AlbumResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetAlbum,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 && resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
