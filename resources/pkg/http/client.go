package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetDataset(url string) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b bytes.Buffer
	r := io.TeeReader(resp.Body, &b)
	ioutil.ReadAll(r)

	r = bytes.NewBuffer(b.Bytes())

	// transform encoding if shift_jis
	if !utf8.Valid(b.Bytes()) {
		r = transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}

	return r, nil
}
