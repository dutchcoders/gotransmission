package transmission

import (
	"testing"

	"github.com/kr/pretty"
)

func TestGet(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentGetRequest{
		Fields: []string{"id", "name", "percentDone", "totalSize", "rateDownload", "rateUpload", "files", "isFinished", "status", "error", "haveValid", "errorString", "peers"},
	}

	if resp, err := wd.Get(d); err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", pretty.Formatter(resp))
	}

}

func TestAdd(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentAddRequest{
		Filename: "magnet:?xt=urn:btih:674D163D2184353CE21F3DE5196B0A6D7C2F9FC2&dn=bbb_sunflower_1080p_60fps_stereo_abl.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_60fps_stereo_abl.mp4",
	}

	if resp, err := wd.Add(d); err != nil {
		t.Error(err)
	} else {
		t.Logf("%#v", pretty.Formatter(resp))
	}
}

func TestSet(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentSetRequest{
		IDs:           []int{3},
		FilesWanted:   []int{1, 3},
		FilesUnwanted: []int{},
	}

	if err := wd.Set(d); err != nil {
		t.Error(err)
	}

}

func TestStartNow(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")
	if err := wd.StartNow(3); err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentRemoveRequest{
		IDs:             []int{3},
		DeleteLocalData: false,
	}

	if err := wd.Remove(d); err != nil {
		t.Error(err)
	}

}
func TestStop(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	if err := wd.Stop(3); err != nil {
		t.Error(err)
	}
}
