package kugou

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetPlaylist(ctx context.Context, specialId string) (*api.PlaylistResponse, error) {
	playlistInfo, err := a.GetPlaylistInfoRaw(ctx, specialId)
	if err != nil {
		return nil, err
	}

	playlistSongs, err := a.GetPlaylistSongsRaw(ctx, specialId, 1, -1)
	if err != nil {
		return nil, err
	}

	n := len(playlistSongs.Data.Info)
	if n == 0 {
		return nil, errors.New("get playlist songs: no data")
	}

	a.patchSongInfo(ctx, playlistSongs.Data.Info...)
	a.patchSongsInfo(ctx, playlistSongs.Data.Info...)
	a.patchSongsLyric(ctx, playlistSongs.Data.Info...)
	songs := resolve(playlistSongs.Data.Info...)
	return &api.PlaylistResponse{
		Id:     strconv.Itoa(playlistInfo.Data.SpecialId),
		Name:   strings.TrimSpace(playlistInfo.Data.SpecialName),
		PicUrl: strings.ReplaceAll(playlistInfo.Data.ImgURL, "{size}", "480"),
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌单信息
func (a *API) GetPlaylistInfoRaw(ctx context.Context, specialId string) (*PlaylistInfoResponse, error) {
	params := ghttp.Params{
		"specialid": specialId,
	}

	resp := new(PlaylistInfoResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetPlaylistInfo,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get playlist info: %s", resp.errorMessage())
	}

	return resp, nil
}

// 获取歌单歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetPlaylistSongsRaw(ctx context.Context, specialId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	params := ghttp.Params{
		"specialid": specialId,
		"page":      strconv.Itoa(page),
		"pagesize":  strconv.Itoa(pageSize),
	}

	resp := new(PlaylistSongsResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetPlaylistSongs,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get playlist songs: %s", resp.errorMessage())
	}

	return resp, nil
}
