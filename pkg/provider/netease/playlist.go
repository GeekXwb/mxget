package netease

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/utils"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.PlaylistResponse, error) {
	_playlistId, err := strconv.Atoi(playlistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetPlaylistRaw(ctx, _playlistId)
	if err != nil {
		return nil, err
	}

	n := resp.Playlist.Total
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	tracks := resp.Playlist.Tracks
	if n > SongRequestLimit {
		extra := n - SongRequestLimit
		trackIds := make([]int, extra)
		for i, j := SongRequestLimit, 0; i < n; i, j = i+1, j+1 {
			trackIds[j] = resp.Playlist.TrackIds[i].Id
		}

		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := 0; i < extra; i += SongRequestLimit {
			if ctx.Err() != nil {
				break
			}

			songIds := trackIds[i:utils.Min(i+SongRequestLimit, extra)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(ctx, songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Songs
			}()
		}

		go func() {
			for s := range queue {
				resp.Playlist.Tracks = append(tracks, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongsURL(ctx, SongDefaultBR, tracks...)
	a.patchSongsLyric(ctx, tracks...)
	songs := resolve(tracks...)
	return &api.PlaylistResponse{
		Id:     strconv.Itoa(resp.Playlist.Id),
		Name:   strings.TrimSpace(resp.Playlist.Name),
		PicUrl: resp.Playlist.PicURL,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId int) (*PlaylistResponse, error) {
	data := map[string]int{
		"id": playlistId,
		"n":  100000,
	}

	resp := new(PlaylistResponse)
	err := a.SendRequest(ghttp.MethodPost, APIGetPlaylist,
		ghttp.WithForm(weapi(data)),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
	}

	return resp, nil
}
