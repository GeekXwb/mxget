package netease

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
	_albumId, err := strconv.Atoi(albumId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetAlbumRaw(ctx, _albumId)
	if err != nil {
		return nil, err
	}

	n := len(resp.Songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsURL(ctx, SongDefaultBR, resp.Songs...)
	a.patchSongsLyric(ctx, resp.Songs...)
	songs := resolve(resp.Songs...)
	return &api.AlbumResponse{
		Id:     strconv.Itoa(resp.Album.Id),
		Name:   strings.TrimSpace(resp.Album.Name),
		PicUrl: resp.Album.PicURL,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId int) (*AlbumResponse, error) {
	resp := new(AlbumResponse)
	err := a.SendRequest(ghttp.MethodPost, fmt.Sprintf(APIGetAlbum, albumId),
		ghttp.WithForm(weapi(struct{}{})),
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
