package http

import (
	"bytes"
	"io"
	"net/http"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func GetDataset(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b bytes.Buffer
	r := io.TeeReader(resp.Body, &b)
	io.ReadAll(r)

	r = bytes.NewBuffer(b.Bytes())

	// transform encoding if shift_jis
	if !utf8.Valid(b.Bytes()) {
		r = transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}

	return r, nil
}
