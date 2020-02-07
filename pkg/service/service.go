package service

import (
	"context"

	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/provider"
)

type MusicServerImpl struct{}

func (m *MusicServerImpl) SearchSongs(ctx context.Context, in *api.SearchSongghttpuest) (*api.SearchSongsResponse, error) {
	client, err := provider.GetClient(in.Platform)
	if err != nil {
		return nil, err
	}

	return client.SearchSongs(ctx, in.Keyword)
}

func (m *MusicServerImpl) GetSong(ctx context.Context, in *api.SongRequest) (*api.SongResponse, error) {
	client, err := provider.GetClient(in.Platform)
	if err != nil {
		return nil, err
	}

	return client.GetSong(ctx, in.Id)
}

func (m *MusicServerImpl) GetAlbum(ctx context.Context, in *api.AlbumRequest) (*api.AlbumResponse, error) {
	client, err := provider.GetClient(in.Platform)
	if err != nil {
		return nil, err
	}

	return client.GetAlbum(ctx, in.Id)
}

func (m *MusicServerImpl) GetPlaylist(ctx context.Context, in *api.PlaylistRequest) (*api.PlaylistResponse, error) {
	client, err := provider.GetClient(in.Platform)
	if err != nil {
		return nil, err
	}

	return client.GetPlaylist(ctx, in.Id)
}

func (m *MusicServerImpl) GetArtist(ctx context.Context, in *api.ArtistRequest) (*api.ArtistResponse, error) {
	client, err := provider.GetClient(in.Platform)
	if err != nil {
		return nil, err
	}

	return client.GetArtist(ctx, in.Id)
}
