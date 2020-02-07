package xiami

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
	resp, err := a.GetAlbumRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	_songs := resp.Data.Data.AlbumDetail.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.AlbumResponse{
		Id:     resp.Data.Data.AlbumDetail.AlbumId,
		Name:   strings.TrimSpace(resp.Data.Data.AlbumDetail.AlbumName),
		PicUrl: resp.Data.Data.AlbumDetail.AlbumLogo,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	token, err := a.getToken(APIGetAlbum)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(albumId)
	if err != nil {
		model["albumStringId"] = albumId
	} else {
		model["albumId"] = albumId
	}
	params := signPayload(token, model)

	resp := new(AlbumResponse)
	err = a.SendRequest(ghttp.MethodGet, APIGetAlbum,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get album: %w", err)
	}

	return resp, nil
}
