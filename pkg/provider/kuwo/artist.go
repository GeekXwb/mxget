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

func (a *API) GetArtist(ctx context.Context, artistId string) (*api.ArtistResponse, error) {
	artistInfo, err := a.GetArtistInfoRaw(ctx, artistId)
	if err != nil {
		return nil, err
	}

	artistSongs, err := a.GetArtistSongsRaw(ctx, artistId, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(artistSongs.Data.List)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongsURL(ctx, SongDefaultBR, artistSongs.Data.List...)
	a.patchSongsLyric(ctx, artistSongs.Data.List...)
	songs := resolve(artistSongs.Data.List...)
	return &api.ArtistResponse{
		Id:     strconv.Itoa(artistInfo.Data.Id),
		Name:   strings.TrimSpace(artistInfo.Data.Name),
		PicUrl: artistInfo.Data.Pic300,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(ctx context.Context, artistId string) (*ArtistInfoResponse, error) {
	params := ghttp.Params{
		"artistid": artistId,
	}

	resp := new(ArtistInfoResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetArtistInfo,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist info: %s", resp.errorMessage())
	}

	return resp, nil
}

// 获取歌手歌曲，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetArtistSongsRaw(ctx context.Context, artistId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	params := ghttp.Params{
		"artistid": artistId,
		"pn":       strconv.Itoa(page),
		"rn":       strconv.Itoa(pageSize),
	}

	resp := new(ArtistSongsResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetArtistSongs,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get artist songs: %s", resp.errorMessage())
	}

	return resp, nil
}
