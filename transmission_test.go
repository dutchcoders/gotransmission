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

	r, err := wd.NewRequest("torrent-get", d)
	if err != nil {
		panic(err)
	}

	var resp TorrentGetResponse
	err = wd.Do(r, &resp)
	if err != nil {
		panic(err)
	}

	pretty.Print(resp)

}

func TestAdd(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentPutRequest{
		Filename: "magnet:?xt=urn:btih:674D163D2184353CE21F3DE5196B0A6D7C2F9FC2&dn=bbb_sunflower_1080p_60fps_stereo_abl.mp4&tr=udp%3a%2f%2ftracker.openbittorrent.com%3a80%2fannounce&tr=udp%3a%2f%2ftracker.publicbt.com%3a80%2fannounce&ws=http%3a%2f%2fdistribution.bbb3d.renderfarming.net%2fvideo%2fmp4%2fbbb_sunflower_1080p_60fps_stereo_abl.mp4",
	}

	r, err := wd.NewRequest("torrent-add", d)
	if err != nil {
		panic(err)
	}

	var resp TorrentPutResponse
	err = wd.Do(r, &resp)
	if err != nil {
		panic(err)
	}

	pretty.Print(resp)

}

func TestSet(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentSetRequest{
		IDs:           []int{3},
		FilesWanted:   []int{1, 3},
		FilesUnwanted: []int{},
	}

	r, err := wd.NewRequest("torrent-set", d)
	if err != nil {
		panic(err)
	}

	err = wd.Do(r, nil)
	if err != nil {
		panic(err)
	}

}

func TestStartNow(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentStartNowRequest{
		IDs: []int{3},
	}

	r, err := wd.NewRequest("torrent-start-now", d)
	if err != nil {
		panic(err)
	}

	err = wd.Do(r, nil)
	if err != nil {
		panic(err)
	}
}

func TestRemove(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentRemoveRequest{
		IDs:             []int{3},
		DeleteLocalData: false,
	}

	r, err := wd.NewRequest("torrent-remove", d)
	if err != nil {
		panic(err)
	}

	err = wd.Do(r, nil)
	if err != nil {
		panic(err)
	}

}
func TestStop(t *testing.T) {
	wd := New("http://localhost:9091/transmission/rpc")

	d := TorrentStopRequest{
		IDs: []int{3},
	}

	r, err := wd.NewRequest("torrent-stop", d)
	if err != nil {
		panic(err)
	}

	var resp interface{}
	err = wd.Do(r, &resp)
	if err != nil {
		panic(err)
	}

	pretty.Print(resp)

}
