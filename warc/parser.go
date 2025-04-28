package warc

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/heyvito/gowarc/httpext"
	"io"
	"strconv"
	"strings"
)

const version = "WARC/1.0"

func Parse(data []byte, length int) (*Record, error) {
	p, err := NewParser(data, length)
	if err != nil {
		return nil, err
	}
	return p.decode()
}

func NewParser(data []byte, length int) (*Parser, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("data too short")
	}

	if data[0] == 0x1F && data[1] == 0x8B {
		zr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		data, err = io.ReadAll(zr)
		_ = zr.Close()
		if err != nil {
			return nil, err
		}
	}

	return &Parser{
		data:   data,
		cursor: 0,
		length: length,
	}, nil
}

type Parser struct {
	data   []byte
	cursor int
	length int
}

func (p *Parser) decode() (*Record, error) {
	fields := map[string]string{}
	if err := p.readVersion(); err != nil {
		return nil, err
	}
	for {
		key, value, err := p.readHeaderField()
		if err != nil {
			return nil, err
		}
		fields[key] = value
		if p.isCRLF() {
			break
		}
	}

	if err := p.readCRLF(); err != nil {
		return nil, err
	}

	var data []byte
	if lenStr, ok := fields["Content-Length"]; ok {
		lenVal, err := strconv.Atoi(lenStr)
		if err != nil {
			return nil, err
		}
		data = p.data[p.cursor : p.cursor+lenVal]
	} else {
		data = p.data[p.cursor:]
	}

	httpParser := httpext.NewResponseParser(data)
	resp, err := httpParser.Parse()
	if err != nil {
		return nil, err
	}
	var contentLength *int
	if v, ok := resp.HeaderNamed("Content-Length"); ok {
		cLen, err := strconv.Atoi(v)
		if err == nil {
			contentLength = &cLen
		}
	}
	return &Record{
		Fields:          fields,
		Body:            resp.Body,
		ContentType:     resp.PtrHeaderNamed("Content-Type"),
		ContentLength:   contentLength,
		StatusCode:      &resp.StatusCode,
		StatusReason:    &resp.ReasonPhrase,
		ContentEncoding: resp.PtrHeaderNamed("Content-Encoding"),
		AllHeaders:      resp.Headers,
	}, nil
}

func (p *Parser) readCRLF() error {
	if p.cursor+1 >= p.length {
		return io.EOF
	}
	if p.data[p.cursor] != 0x0D || p.data[p.cursor+1] != 0x0A {
		return errors.New("expected CRLF")
	}
	p.cursor += 2
	return nil
}

func (p *Parser) eof() bool {
	return p.cursor >= p.length
}

func (p *Parser) readUntilCRLF() (string, error) {
	var str []rune
	var c byte
	for !p.eof() {
		c = p.data[p.cursor]
		if c == 0x0D && p.data[p.cursor+1] == 0x0A {
			break
		}
		str = append(str, rune(c))
		p.cursor++
	}
	return string(str), nil
}

func (p *Parser) readUntil(chr rune) string {
	var str []rune
	var c byte
	for !p.eof() {
		c = p.data[p.cursor]
		if rune(c) == chr {
			break
		}
		str = append(str, rune(c))
		p.cursor++
	}
	return string(str)
}

func (p *Parser) readVersion() error {
	v, err := p.readUntilCRLF()
	if err != nil {
		return err
	}
	if v != version {
		return fmt.Errorf("WARC version mismatch (expected %s, got %s)", version, v)
	}
	return p.readCRLF()
}

func (p *Parser) readHeaderField() (string, string, error) {
	name := p.readUntil(':')
	p.cursor += 1
	value, err := p.readUntilCRLF()
	if err != nil {
		return "", "", err
	}
	if err = p.readCRLF(); err != nil {
		return "", "", err
	}
	return strings.TrimSpace(name), strings.TrimSpace(value), nil
}

func (p *Parser) isCRLF() bool {
	if p.cursor+1 >= p.length {
		return true
	}
	return p.data[p.cursor] == 0x0D && p.data[p.cursor+1] == 0x0A
}
