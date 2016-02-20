package transmission

// https://trac.transmissionbt.com/browser/trunk/extras/rpc-spec.txt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	StatusStopped = iota
	StatusCheckWait
	StatusCheck
	StatusDownloadWait
	StatusDownload
	StatusSeedWait
	StatusSeed
)

type TorrentGetRequest struct {
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
	FilesWanted   []int `json:"files-wanted,omitempty"`
	FilesUnwanted []int `json:"files-unwanted,omitempty"`
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
	Id           int64         `json:"id"`
	IsFinished   bool          `json:"isFinished"`
	Name         string        `json:"name"`
	Peers        []interface{} `json:"peers"`
	PercentDone  float64       `json:"percentDone"`
	RateDownload int64         `json:"rateDownload"`
	RateUpload   int64         `json:"rateUpload"`
	Status       int64         `json:"status"`
	TotalSize    int64         `json:"totalSize"`
}

type TorrentPutRequest struct {
	Filename string `json:"filename,omitempty"`
}

type TorrentPutResponse struct {
	TorrentAdded struct {
		Hash string  `json:"hashString"`
		ID   float64 `json:"id"`
		Name string  `json:"name"`
	} `json:"torrent-added"`
}

func (c *transmission) NewRequest(method string, arguments interface{}) (*http.Request, error) {
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

type transmission struct {
	client         *http.Client
	url            string
	transmissionId string
}

func New(url string) *transmission {
	return &transmission{
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

func (t *transmission) do(req *http.Request) (*transmissionResponse, error) {
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

		r := io.TeeReader(resp.Body, os.Stderr)

		err = json.NewDecoder(r).Decode(&dr)
		if err != nil {
			return nil, err
		}

		if dr.Result != "success" {
			fmt.Printf("%#v", dr)
			return &dr, &Error{
				Result: dr.Result,
			}
		}

		return &dr, nil
	}
}

func (t *transmission) Do(req *http.Request, v interface{}) error {
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
