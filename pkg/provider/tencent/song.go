package tencent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

const (
	// 跟vkey配合获取歌曲下载地址，可为任意字符串
	Guid = "0"
)

func (a *API) GetSong(ctx context.Context, songMid string) (*api.SongResponse, error) {
	resp, err := a.GetSongRaw(ctx, songMid)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Data[0]
	a.patchSongsURLV1(ctx, _song)
	a.patchSongsLyric(ctx, _song)
	songs := resolve(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(ctx context.Context, songMid string) (*SongResponse, error) {
	params := ghttp.Params{
		"songmid": songMid,
	}

	resp := new(SongResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSong,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song: %d", resp.Code)
	}

	return resp, nil
}

func (a *API) GetSongURLV1(ctx context.Context, songMid string, mediaMid string) (string, error) {
	resp, err := a.GetSongURLV1Raw(ctx, songMid, mediaMid)
	if err != nil {
		return "", err
	}
	if len(resp.Data.Items) == 0 {
		return "", errors.New("get song url: no data")
	}

	item := resp.Data.Items[0]
	if item.SubCode != 0 {
		return "", fmt.Errorf("get song url: %d", item.SubCode)
	}

	return fmt.Sprintf(SongURL, item.FileName, item.Vkey), nil
}

// 获取歌曲播放地址
func (a *API) GetSongURLV1Raw(ctx context.Context, songMid string, mediaMid string) (*SongURLResponseV1, error) {
	params := ghttp.Params{
		"songmid":  songMid,
		"filename": "M500" + mediaMid + ".mp3",
	}

	resp := new(SongURLResponseV1)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongURLV1,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %s", resp.ErrInfo)
	}

	return resp, nil
}

func (a *API) GetSongURLV2(ctx context.Context, songMid string) (string, error) {
	resp, err := a.GetSongsURLV2Raw(ctx, songMid)
	if err != nil {
		return "", err
	}
	if len(resp.Req0.Data.MidURLInfo) == 0 {
		return "", errors.New("get song url: no data")
	}

	n := len(resp.Req0.Data.Sip)
	if n == 0 {
		return "", errors.New("get song url: no sip")
	}

	// 随机获取一个sip
	sip := resp.Req0.Data.Sip[rand.Intn(n)]
	urlInfo := resp.Req0.Data.MidURLInfo[0]
	if urlInfo.PURL == "" {
		return "", errors.New("get song url: copyright protection")
	}

	return sip + urlInfo.PURL, nil
}

// 批量获取歌曲播放地址
func (a *API) GetSongsURLV2Raw(ctx context.Context, songMids ...string) (*SongURLResponseV2, error) {
	if len(songMids) > SongURLRequestLimit {
		songMids = songMids[:SongURLRequestLimit]
	}

	param := map[string]interface{}{
		"guid":      Guid,
		"loginflag": 1,
		"songmid":   songMids,
		"uin":       "0",
		"platform":  "20",
	}
	req0 := map[string]interface{}{
		"module": "vkey.GetVkeyServer",
		"method": "CgiGetVkey",
		"param":  param,
	}
	data := map[string]interface{}{
		"req0": req0,
	}

	enc, _ := json.Marshal(data)
	params := ghttp.Params{
		"data": string(enc),
	}
	resp := new(SongURLResponseV2)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongsURLV2,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %d", resp.Code)
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songMid string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songMid)
	if err != nil {
		return "", err
	}

	// lyric, err := base64.StdEncoding.DecodeString(resp.Lyric)
	// if err != nil {
	// 	return "", err
	// }

	return resp.Lyric, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songMid string) (*SongLyricResponse, error) {
	params := ghttp.Params{
		"songmid": songMid,
	}

	resp := new(SongLyricResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongLyric,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("get song lyric: %d", resp.Code)
	}

	return resp, nil
}
