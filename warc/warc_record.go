package warc

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
)

type Record struct {
	Fields          map[string]string
	Body            []byte
	ContentType     *string
	ContentLength   *int
	StatusCode      *int
	StatusReason    *string
	ContentEncoding *string
	AllHeaders      [][]string
}

func (r *Record) IsGzipped() bool {
	if r.ContentEncoding == nil || *r.ContentEncoding == "" {
		return false
	}
	enc := strings.ToLower(*r.ContentEncoding)
	return strings.Contains(enc, "gzip")
}

func (r *Record) GetBody() ([]byte, error) {
	if r.IsGzipped() {
		zr, err := gzip.NewReader(bytes.NewReader(r.Body))
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = zr.Close()
		}()
		return io.ReadAll(zr)
	} else {
		return r.Body, nil
	}
}

func (r *Record) HeaderNamed(name string) (string, bool) {
	name = strings.ToLower(name)
	for _, header := range r.AllHeaders {
		if strings.ToLower(header[0]) == name {
			return header[1], true
		}
	}
	return "", false
}
