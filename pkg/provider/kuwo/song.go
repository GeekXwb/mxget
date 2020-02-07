package kuwo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, mid string) (*api.SongResponse, error) {
	resp, err := a.GetSongRaw(ctx, mid)
	if err != nil {
		return nil, err
	}

	a.patchSongsURL(ctx, SongDefaultBR, &resp.Data)
	a.patchSongsLyric(ctx, &resp.Data)
	songs := resolve(&resp.Data)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(ctx context.Context, mid string) (*SongResponse, error) {
	params := ghttp.Params{
		"mid": mid,
	}

	resp := new(SongResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSong,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongURL(ctx context.Context, mid int, br int) (string, error) {
	resp, err := a.GetSongURLRaw(ctx, mid, br)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

// 获取歌曲播放地址，br: 比特率，128/192/320 可选
func (a *API) GetSongURLRaw(ctx context.Context, mid int, br int) (*SongURLResponse, error) {
	var _br int
	switch br {
	case 128, 192, 320:
		_br = br
	default:
		_br = 320
	}
	params := ghttp.Params{
		"rid": strconv.Itoa(mid),
		"br":  fmt.Sprintf("%dkmp3", _br),
	}

	resp := new(SongURLResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongURL,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get song url: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, mid int) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, mid)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, i := range resp.Data.LrcList {
		t, err := strconv.ParseFloat(i.Time, 64)
		if err != nil {
			return "", err
		}
		m := int(t / 60)
		s := int(t - float64(m*60))
		ms := int((t - float64(m*60) - float64(s)) * 100)
		sb.WriteString(fmt.Sprintf("[%02d:%02d:%02d]%s\n", m, s, ms, i.LineLyric))
	}

	return sb.String(), nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, mid int) (*SongLyricResponse, error) {
	params := ghttp.Params{
		"musicId": strconv.Itoa(mid),
	}

	resp := new(SongLyricResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongLyric,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("get song lyric: %s", resp.Msg)
	}

	return resp, nil
}
