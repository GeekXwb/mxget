package tencent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetArtist(ctx context.Context, singerMid string) (*api.ArtistResponse, error) {
	resp, err := a.GetArtistRaw(ctx, singerMid, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	_songs := make([]*Song, n)
	for i, v := range resp.Data.List {
		_songs[i] = v.MusicData
	}

	a.patchSongsURLV1(ctx, _songs...)
	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.ArtistResponse{
		Id:     resp.Data.SingerMid,
		Name:   strings.TrimSpace(resp.Data.SingerName),
		PicUrl: fmt.Sprintf(ArtistPicURL, resp.Data.SingerMid),
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(ctx context.Context, singerMid string, page int, pageSize int) (*ArtistResponse, error) {
	params := ghttp.Params{
		"singermid": singerMid,
		"begin":     strconv.Itoa(page),
		"num":       strconv.Itoa(pageSize),
	}

	resp := new(ArtistResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetArtist,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get artist: %d", resp.Code)
	}

	return resp, nil
}
