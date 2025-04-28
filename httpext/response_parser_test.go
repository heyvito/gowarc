package httpext

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPVersionStatusReasonParser(t *testing.T) {
	data := "HTTP/1.1 200 OK\r\n"
	p := NewResponseParser([]byte(data))
	version, status, reason, err := p.parseHTTPVersionAndStatus()
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1", version)
	assert.Equal(t, 200, status)
	assert.Equal(t, "OK", reason)
}

func TestHTTPHeaderParser(t *testing.T) {
	data := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 0\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n"
	p := NewResponseParser([]byte(data))
	_, _, _, err := p.parseHTTPVersionAndStatus()
	require.NoError(t, err)
	headers, err := p.parseHeaders()
	require.NoError(t, err)
	assert.Equal(t, 2, len(headers))
	assert.Equal(t, headers[0], []string{"Content-Length", "0"})
	assert.Equal(t, headers[1], []string{"Content-Type", "text/plain"})
}

func TestFullParse(t *testing.T) {
	data := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 0\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"This is the content"

	p := NewResponseParser([]byte(data))
	r, err := p.Parse()
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1", r.HTTPVersion)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, "OK", r.ReasonPhrase)
	assert.Equal(t, r.Headers[0], []string{"Content-Length", "0"})
	assert.Equal(t, r.Headers[1], []string{"Content-Type", "text/plain"})
	assert.Equal(t, []byte("This is the content"), r.Body)
}
