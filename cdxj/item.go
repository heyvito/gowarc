package cdxj

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var httpSchemeRegexp = regexp.MustCompile(`^https?://`)

type intermediateItem struct {
	URL      string `json:"url"`
	MimeType string `json:"mime"`
	Status   string `json:"status"`
	Digest   string `json:"digest"`
	Length   string `json:"length"`
	Offset   string `json:"offset"`
	Filename string `json:"filename"`
}

type Item struct {
	Surt        string
	Timestamp   string
	URL         string
	MimeType    string
	Status      int
	Digest      string
	Length      int
	Offset      int
	Filename    string
	SymbolicURL string
}

func ItemFromString(line string) (*Item, error) {
	components := strings.SplitN(line, " ", 3)
	if len(components) != 3 {
		return nil, errors.New("invalid CDXJ item")
	}

	surt, timestamp, metadata := components[0], components[1], components[2]
	var decoded intermediateItem
	if err := json.Unmarshal([]byte(metadata), &decoded); err != nil {
		return nil, err
	}
	status, err := strconv.Atoi(decoded.Status)
	if err != nil {
		return nil, err
	}
	length, err := strconv.ParseInt(decoded.Length, 10, 64)
	if err != nil {
		return nil, err
	}
	offset, err := strconv.ParseInt(decoded.Offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Item{
		Surt:        surt,
		Timestamp:   timestamp,
		URL:         decoded.URL,
		MimeType:    decoded.MimeType,
		Status:      status,
		Digest:      decoded.Digest,
		Length:      int(length),
		Offset:      int(offset),
		Filename:    decoded.Filename,
		SymbolicURL: httpSchemeRegexp.ReplaceAllString(decoded.URL, ""),
	}, nil
}
