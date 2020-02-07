package provider

import (
	"context"
	"errors"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/provider/baidu"
	"github.com/winterssy/mxget/pkg/provider/kugou"
	"github.com/winterssy/mxget/pkg/provider/kuwo"
	"github.com/winterssy/mxget/pkg/provider/migu"
	"github.com/winterssy/mxget/pkg/provider/netease"
	"github.com/winterssy/mxget/pkg/provider/tencent"
	"github.com/winterssy/mxget/pkg/provider/xiami"
)

type (
	API interface {
		// 搜索歌曲
		SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error)
		// 获取单曲
		GetSong(ctx context.Context, songId string) (*api.SongResponse, error)
		// 获取歌手
		GetArtist(ctx context.Context, artistId string) (*api.ArtistResponse, error)
		// 获取专辑
		GetAlbum(ctx context.Context, albumId string) (*api.AlbumResponse, error)
		// 获取歌单
		GetPlaylist(ctx context.Context, playlistId string) (*api.PlaylistResponse, error)
		// 网络请求
		SendRequest(method string, url string, opts ...ghttp.RequestOption) *ghttp.Response
	}
)

func GetClient(platform string) (API, error) {
	switch platform {
	case "netease", "nc":
		return netease.Client(), nil
	case "tencent", "qq":
		return tencent.Client(), nil
	case "migu", "mg":
		return migu.Client(), nil
	case "kugou", "kg":
		return kugou.Client(), nil
	case "kuwo", "kw":
		return kuwo.Client(), nil
	case "xiami", "xm":
		return xiami.Client(), nil
	case "qianqian", "baidu", "bd":
		return baidu.Client(), nil
	default:
		return nil, errors.New("unexpected music platform")
	}
}

func GetDesc(platform string) string {
	switch platform {
	case "netease", "nc":
		return "netease cloud music"
	case "tencent", "qq":
		return "qq music"
	case "migu", "mg":
		return "migu music"
	case "kugou", "kg":
		return "kugou music"
	case "kuwo", "kw":
		return "kuwo music"
	case "xiami", "xm":
		return "xiami music"
	case "qianqian", "baidu", "bd":
		return "qianqian music"
	default:
		return "unknown"
	}
}
