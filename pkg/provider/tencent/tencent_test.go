package tencent_test

import (
	"context"
	"os"
	"testing"

	"github.com/winterssy/mxget/pkg/provider/tencent"
)

var (
	client *tencent.API
	ctx    context.Context
)

func setup() {
	client = tencent.New(nil)
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs(ctx, "Alan Walker")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong(ctx, "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(song)
}

func TestAPI_GetSongURLV1(t *testing.T) {
	url, err := client.GetSongURLV1(ctx, "002Zkt5S2z8JZx", "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongURLV2(t *testing.T) {
	url, err := client.GetSongURLV2(ctx, "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(ctx, "002Zkt5S2z8JZx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lyric)
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist(ctx, "000Sp0Bz4JXH0o")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(artist)
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum(ctx, "002fRO0N4FftzY")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(album)
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist(ctx, "5474239760")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(playlist)
}
