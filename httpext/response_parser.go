package httpext

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func NewResponseParser(data []byte) *ResponseParser {
	return &ResponseParser{
		data:   data,
		cursor: 0,
	}
}

type ResponseParser struct {
	data   []byte
	cursor int
}

func (p *ResponseParser) Parse() (*HTTPResponse, error) {
	version, status, reason, err := p.parseHTTPVersionAndStatus()
	if err != nil {
		return nil, err
	}
	headers, err := p.parseHeaders()
	return &HTTPResponse{
		HTTPVersion:  version,
		StatusCode:   status,
		ReasonPhrase: reason,
		Headers:      headers,
		Body:         p.data[p.cursor:],
	}, nil
}

func (p *ResponseParser) parseHTTPVersionAndStatus() (string, int, string, error) {
	mark := p.cursor
	if !p.Match("HTTP/1.") {
		return "", 0, "", fmt.Errorf("unsupported HTTP version")
	}

	if err := p.AdvanceUntil(' '); err != nil {
		return "", 0, "", err
	}

	version := string(p.data[mark:p.cursor])
	p.Advance(1)

	mark = p.cursor
	if err := p.AdvanceUntil(' '); err != nil {
		return "", 0, "", err
	}
	rawStatus := string(p.data[mark:p.cursor])
	status, err := strconv.Atoi(rawStatus)
	if err != nil {
		return "", 0, "", err
	}
	p.Advance(1)

	mark = p.cursor
	p.AdvanceToCRLF()
	reason := string(p.data[mark:p.cursor])
	if err = p.ConsumeCRLF(); err != nil {
		return "", 0, "", err
	}

	return version, status, reason, nil
}

func (p *ResponseParser) parseHeaders() ([][]string, error) {
	var headers [][]string
	for {
		key, value, err := p.parseHeader()
		if err != nil {
			return nil, err
		}
		headers = append(headers, []string{key, value})
		if p.IsCRLF() {
			if err = p.ConsumeCRLF(); err != nil {
				return nil, err
			}
			break
		}
	}
	return headers, nil
}

func (p *ResponseParser) parseHeader() (string, string, error) {
	mark := p.cursor
	if err := p.AdvanceUntil(':'); err != nil {
		return "", "", err
	}
	key := string(p.data[mark:p.cursor])
	if !p.Match(":") {
		return "", "", fmt.Errorf("expected a colon after header key")
	}
	mark = p.cursor
	p.AdvanceToCRLF()
	value := p.data[mark:p.cursor]
	if err := p.ConsumeCRLF(); err != nil {
		return "", "", err
	}
	return key, strings.TrimSpace(string(value)), nil
}

func (p *ResponseParser) Peek(n int) (uint8, error) {
	if p.cursor+n > len(p.data) {
		return 0, io.EOF
	}
	return p.data[p.cursor+n], nil
}

func (p *ResponseParser) Match(str string) bool {
	matched := true
	for i, c := range []rune(str) {
		if rune(p.data[p.cursor+i]) != c {
			matched = false
			break
		}
	}
	if matched {
		p.Advance(len(str))
	}
	return matched
}

func (p *ResponseParser) Advance(n int) {
	p.cursor += n
}

func (p *ResponseParser) AdvanceUntil(chr byte) error {
	for {
		peek, err := p.Peek(0)
		if err != nil {
			return err
		}
		if peek != chr {
			p.Advance(1)
		} else {
			return nil
		}
	}
}

func (p *ResponseParser) AdvanceToCRLF() {
	for {
		if p.IsCRLF() {
			break
		}
		p.Advance(1)
	}
}

func (p *ResponseParser) ConsumeCRLF() error {
	if p.IsCRLF() {
		p.Advance(2)
		return nil
	}

	return fmt.Errorf("expected CRLF")
}

func (p *ResponseParser) IsCRLF() bool {
	p1, err := p.Peek(0)
	if err != nil {
		return false
	}
	p2, err := p.Peek(1)
	if err != nil {
		return false
	}

	return p1 == 13 && p2 == 10
}

func (p *ResponseParser) IsDoubleCRLF() bool {
	p1, err := p.Peek(0)
	if err != nil {
		return false
	}
	p2, err := p.Peek(1)
	if err != nil {
		return false
	}
	p3, err := p.Peek(2)
	if err != nil {
		return false
	}
	p4, err := p.Peek(3)
	if err != nil {
		return false
	}

	return p1 == 13 && p2 == 10 && p3 == 13 && p4 == 10
}
