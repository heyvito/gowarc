package httpext

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
)

type HTTPResponse struct {
	HTTPVersion  string
	StatusCode   int
	ReasonPhrase string
	Headers      [][]string
	Body         []byte
}

func (r *HTTPResponse) HeaderNamed(name string) (string, bool) {
	name = strings.ToLower(name)
	for _, header := range r.Headers {
		if strings.ToLower(header[0]) == name {
			return header[1], true
		}
	}
	return "", false
}

func (r *HTTPResponse) HeadersNamed(name string) []string {
	var values []string
	name = strings.ToLower(name)
	for _, header := range r.Headers {
		if strings.ToLower(header[0]) == name {
			values = append(values, header[1])
		}
	}
	return values
}

func (r *HTTPResponse) IsGzipped() bool {
	for _, v := range r.HeadersNamed("Content-Encoding") {
		if strings.Contains(strings.ToLower(v), "gzip") {
			return true
		}
	}
	return false
}

func (r *HTTPResponse) Decompress() ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = zr.Close() }()
	return io.ReadAll(zr)
}

func (r *HTTPResponse) PtrHeaderNamed(name string) *string {
	v, ok := r.HeaderNamed(name)
	if ok {
		return &v
	}
	return nil
}
