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

func (a *API) GetArtist(ctx context.Context, artistId string) (*api.ArtistResponse, error) {
	_artistId, err := strconv.Atoi(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetArtistRaw(ctx, _artistId)
	if err != nil {
		return nil, err
	}

	n := len(resp.HotSongs)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongsURL(ctx, SongDefaultBR, resp.HotSongs...)
	a.patchSongsLyric(ctx, resp.HotSongs...)
	songs := resolve(resp.HotSongs...)
	return &api.ArtistResponse{
		Id:     strconv.Itoa(resp.Artist.Id),
		Name:   strings.TrimSpace(resp.Artist.Name),
		PicUrl: resp.Artist.PicURL,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(ctx context.Context, artistId int) (*ArtistResponse, error) {
	resp := new(ArtistResponse)
	err := a.SendRequest(ghttp.MethodPost, fmt.Sprintf(APIGetArtist, artistId),
		ghttp.WithForm(weapi(struct{}{})),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist: %s", resp.errorMessage())
	}

	return resp, nil
}
