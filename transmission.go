package transmission

// https://trac.transmissionbt.com/browser/trunk/extras/rpc-spec.txt

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kr/pretty"
)

type Status uint64

const (
	StatusStopped Status = iota
	StatusCheckWait
	StatusCheck
	StatusDownloadWait
	StatusDownload
	StatusSeedWait
	StatusSeed
)

type TorrentGetRequest struct {
	IDs    []int    `json:"ids,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

type TorrentGetResponse struct {
	Torrents []Torrent `json:"torrents"`
}

type TorrentStartNowRequest struct {
	IDs []int `json:"ids,omitempty"`
}

type TorrentStopRequest struct {
	IDs []int `json:"ids,omitempty"`
}

type TorrentSetRequest struct {
	IDs           []int `json:"ids,omitempty"`
	FilesWanted   []int `json:"files-wanted"`
	FilesUnwanted []int `json:"files-unwanted"`
}

type TorrentRemoveRequest struct {
	IDs             []int `json:"ids,omitempty"`
	DeleteLocalData bool  `json:"delete-local-data"`
}

type Torrent struct {
	Error       int64  `json:"error"`
	ErrorString string `json:"errorString"`
	Files       []struct {
		BytesCompleted int64  `json:"bytesCompleted"`
		Length         int64  `json:"length"`
		Name           string `json:"name"`
	} `json:"files"`
	HaveValid    int64         `json:"haveValid"`
	ID           int64         `json:"id"`
	IsFinished   bool          `json:"isFinished"`
	Name         string        `json:"name"`
	Peers        []interface{} `json:"peers"`
	PercentDone  float64       `json:"percentDone"`
	RateDownload int64         `json:"rateDownload"`
	RateUpload   int64         `json:"rateUpload"`
	Status       Status        `json:"status"`
	TotalSize    int64         `json:"totalSize"`
}

type TorrentAddRequest struct {
	Filename string `json:"filename,omitempty"`
}

type TorrentAddResponse struct {
	TorrentAdded struct {
		Hash string `json:"hashString"`
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"torrent-added"`
}

func (t *Transmission) Add(d TorrentAddRequest) (*TorrentAddResponse, error) {
	r, err := t.NewRequest("torrent-add", d)
	if err != nil {
		return nil, err
	}

	var resp TorrentAddResponse
	err = t.Do(r, &resp)
	return &resp, err
}

func (t *Transmission) Get(d TorrentGetRequest) (*TorrentGetResponse, error) {
	r, err := t.NewRequest("torrent-get", d)
	if err != nil {
		return nil, err
	}

	var resp TorrentGetResponse
	err = t.Do(r, &resp)
	return &resp, err
}

func (t *Transmission) Set(d TorrentSetRequest) error {
	r, err := t.NewRequest("torrent-set", d)
	if err != nil {
		return err
	}

	var v interface{}
	defer func() {
		pretty.Print(v)
	}()

	return t.Do(r, &v)
}

func (t *Transmission) Remove(d TorrentRemoveRequest) error {
	r, err := t.NewRequest("torrent-remove", d)
	if err != nil {
		return err
	}

	return t.Do(r, nil)
}

func (t *Transmission) Stop(ids ...int) error {
	d := TorrentStopRequest{
		IDs: ids,
	}

	r, err := t.NewRequest("torrent-stop", d)
	if err != nil {
		return err
	}

	return t.Do(r, nil)
}

func (t *Transmission) StartNow(ids ...int) error {
	d := TorrentStartNowRequest{
		IDs: ids,
	}

	r, err := t.NewRequest("torrent-start-now", d)
	if err != nil {
		return err
	}

	return t.Do(r, nil)
}

func (c *Transmission) NewRequest(method string, arguments interface{}) (*http.Request, error) {
	body := struct {
		Arguments interface{} `json:"arguments"`
		Method    string      `json:"method"`
	}{
		Arguments: arguments,
		Method:    method,
	}

	var buf io.ReadWriter = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/json; charset=UTF-8")
	req.Header.Add("Accept", "text/json")
	return req, nil
}

type Transmission struct {
	client         *http.Client
	url            string
	transmissionId string
}

func New(url string) *Transmission {
	return &Transmission{
		client: http.DefaultClient,
		url:    url,
	}
}

type Error struct {
	Result string
}

func (e *Error) Error() string {
	return e.Result
}

type transmissionResponse struct {
	Result    string           `json:"result"`
	Arguments *json.RawMessage `json:"arguments"`
}

func (t *Transmission) do(req *http.Request) (*transmissionResponse, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	for {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

		req.Header.Set("X-Transmission-Session-Id", t.transmissionId)

		resp, err := t.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode == 409 {
			t.transmissionId = resp.Header.Get("X-Transmission-Session-Id")
			continue
		}

		defer resp.Body.Close()

		dr := transmissionResponse{}

		r := resp.Body // io.TeeReader(resp.Body, os.Stderr)

		err = json.NewDecoder(r).Decode(&dr)
		if err != nil {
			return nil, err
		}

		if dr.Result != "success" {
			return &dr, &Error{
				Result: dr.Result,
			}
		}

		return &dr, nil
	}
}

func (t *Transmission) Do(req *http.Request, v interface{}) error {
	dr, err := t.do(req)
	if err != nil {
		return err
	}

	switch v := v.(type) {
	case io.Writer:
		value := ""
		if err = json.Unmarshal(*dr.Arguments, &value); err != nil {
			return err
		}

		v.Write([]byte(value))
	case interface{}:
		return json.Unmarshal(*dr.Arguments, &v)
	}

	return nil
}
