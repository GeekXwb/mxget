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

func (a *API) GetArtist(ctx context.Context, artistId string) (*api.ArtistResponse, error) {
	artistInfo, err := a.GetArtistInfoRaw(ctx, artistId)
	if err != nil {
		return nil, err
	}

	artistSongs, err := a.GetArtistSongsRaw(ctx, artistId, 1, 50)
	if err != nil {
		return nil, err
	}

	_songs := artistSongs.Data.Data.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.ArtistResponse{
		Id:     artistInfo.Data.Data.ArtistDetailVO.ArtistId,
		Name:   strings.TrimSpace(artistInfo.Data.Data.ArtistDetailVO.ArtistName),
		PicUrl: artistInfo.Data.Data.ArtistDetailVO.ArtistLogo,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(ctx context.Context, artistId string) (*ArtistInfoResponse, error) {
	token, err := a.getToken(APIGetArtistInfo)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := signPayload(token, model)

	resp := new(ArtistInfoResponse)
	err = a.SendRequest(ghttp.MethodGet, APIGetArtistInfo,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist info: %w", err)
	}

	return resp, nil
}

// 获取歌手歌曲
func (a *API) GetArtistSongsRaw(ctx context.Context, artistId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	token, err := a.getToken(APIGetArtistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := signPayload(token, model)

	resp := new(ArtistSongsResponse)
	err = a.SendRequest(ghttp.MethodGet, APIGetArtistSongs,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist songs: %w", err)
	}

	return resp, nil
}
