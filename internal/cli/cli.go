package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bogem/id3v2"
	"github.com/winterssy/ghttp"
	"github.com/winterssy/glog"
	"github.com/winterssy/mxget/internal/settings"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/utils"
)

func ConcurrentDownload(ctx context.Context, client provider.API, savePath string, songs ...*api.SongResponse) {
	savePath = filepath.Join(settings.Cfg.Dir, utils.TrimInvalidFilePathChars(savePath))
	if err := os.MkdirAll(savePath, 0755); err != nil {
		glog.Fatal(err)
	}

	var limit int
	switch {
	case settings.Limit < 1:
		limit = runtime.NumCPU()
	case settings.Limit > 32:
		limit = 32
	default:
		limit = settings.Limit
	}

	c := concurrency.New(limit)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *api.SongResponse) {
			defer c.Done()
			songInfo := fmt.Sprintf("%s - %s", s.Artist, s.Name)
			if !s.Playable {
				glog.Errorf("Download [%s] failed: song unavailable", songInfo)
				return
			}

			filePath := filepath.Join(savePath, utils.TrimInvalidFilePathChars(songInfo))
			glog.Infof("Start download [%s]", songInfo)
			mp3FilePath := filePath + ".mp3"
			if !settings.Force {
				_, err := os.Stat(mp3FilePath)
				if err == nil {
					glog.Infof("Song already downloaded: [%s]", songInfo)
					return
				}
			}

			err := client.
				SendRequest(ghttp.MethodGet, s.Url,
					ghttp.WithContext(ctx),
				).
				Save(mp3FilePath, 0664)
			if err != nil {
				glog.Errorf("Download [%s] failed: %v", songInfo, err)
				_ = os.Remove(mp3FilePath)
				return
			}
			glog.Infof("Download [%s] complete", songInfo)

			if settings.Tag {
				glog.Infof("Update music metadata: [%s]", songInfo)
				writeTag(ctx, client, mp3FilePath, s)
			}

			if settings.Lyric && s.Lyric != "" {
				glog.Infof("Save lyric: [%s]", songInfo)
				lrcFilePath := filePath + ".lrc"
				saveLyric(lrcFilePath, s.Lyric)
			}
		}(s)
	}
	c.Wait()
}

func saveLyric(filePath string, lyric string) {
	err := ioutil.WriteFile(filePath, []byte(lyric), 0644)
	if err != nil {
		_ = os.Remove(filePath)
	}
}

func writeTag(ctx context.Context, client provider.API, filePath string, song *api.SongResponse) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(song.Name)
	tag.SetArtist(song.Artist)
	tag.SetAlbum(song.Album)

	if song.Lyric != "" {
		uslt := id3v2.UnsynchronisedLyricsFrame{
			Encoding:          id3v2.EncodingUTF8,
			Language:          "eng",
			ContentDescriptor: song.Name,
			Lyrics:            song.Lyric,
		}
		tag.AddUnsynchronisedLyricsFrame(uslt)
	}

	if song.PicUrl != "" {
		pic, err := client.SendRequest(ghttp.MethodGet, song.PicUrl,
			ghttp.WithContext(ctx),
		).Content()
		if err == nil {
			picFrame := id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    "image/jpeg",
				PictureType: id3v2.PTFrontCover,
				Description: "Front cover",
				Picture:     pic,
			}
			tag.AddAttachedPicture(picFrame)
		}
	}

	_ = tag.Save()
}
