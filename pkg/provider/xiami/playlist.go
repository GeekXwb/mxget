package xiami

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
	resp, err := a.GetPlaylistDetailRaw(ctx, playlistId, 1, SongRequestLimit)
	if err != nil {
		return nil, err
	}

	n, _ := strconv.Atoi(resp.Data.Data.CollectDetail.SongCount)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	_songs := resp.Data.Data.CollectDetail.Songs
	if n > SongRequestLimit {
		allSongs := resp.Data.Data.CollectDetail.AllSongs
		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := SongRequestLimit; i < n; i += SongRequestLimit {
			if ctx.Err() != nil {
				break
			}

			songIds := allSongs[i:utils.Min(i+SongRequestLimit, n)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(ctx, songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Data.Data.Songs
			}()
		}

		go func() {
			for s := range queue {
				_songs = append(_songs, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := resolve(_songs...)
	return &api.PlaylistResponse{
		Id:     resp.Data.Data.CollectDetail.ListId,
		Name:   strings.TrimSpace(resp.Data.Data.CollectDetail.CollectName),
		PicUrl: resp.Data.Data.CollectDetail.CollectLogo,
		Count:  uint32(n),
		Songs:  songs,
	}, nil
}

// 获取歌单详情，包含歌单信息跟歌曲
func (a *API) GetPlaylistDetailRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistDetailResponse, error) {
	token, err := a.getToken(APIGetPlaylistDetail)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

	resp := new(PlaylistDetailResponse)
	err = a.SendRequest(ghttp.MethodGet, APIGetPlaylistDetail,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get playlist detail: %w", err)
	}

	return resp, nil
}

// 获取歌单歌曲，不包含歌单信息
func (a *API) GetPlaylistSongsRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	token, err := a.getToken(APIGetPlaylistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

	resp := new(PlaylistSongsResponse)
	err = a.SendRequest(ghttp.MethodGet, APIGetPlaylistSongs,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get playlist songs: %w", err)
	}

	return resp, nil
}
