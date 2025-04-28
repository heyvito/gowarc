package warc

import (
	"errors"
	"github.com/heyvito/gowarc/cdxj"
)

type IO struct {
	data []byte
}

func NewIO(data []byte) (*IO, error) {
	if len(data) < 2 {
		return nil, errors.New("data too short")
	}

	return &IO{data: data}, nil
}

func (p *IO) ReadRecord(offset, length int) (*Record, error) {
	return Parse(p.data[offset:offset+length], length)
}

func (p *IO) ReadIndexedRecord(record *cdxj.Item) (*Record, error) {
	return p.ReadRecord(record.Offset, record.Length)
}
