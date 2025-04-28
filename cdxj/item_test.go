package cdxj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
import "github.com/stretchr/testify/require"

func TestItemFromString(t *testing.T) {
	line := `io,vito)/assets/css/main.css?cache=20250130t095925 20250402021331 {"url":"https://vito.io/assets/css/main.css?cache=20250130T095925","mime":"text/css","status":"200","digest":"OLQHSNRAKBAHCDGUFOWPKV5EPFFJLVVL","length":"4726","offset":"10229","filename":"data.warc.gz"}`
	item, err := ItemFromString(line)
	require.NoError(t, err)
	assert.Equal(t, "https://vito.io/assets/css/main.css?cache=20250130T095925", item.URL)
	assert.Equal(t, "text/css", item.MimeType)
	assert.Equal(t, 200, item.Status)
	assert.Equal(t, "OLQHSNRAKBAHCDGUFOWPKV5EPFFJLVVL", item.Digest)
	assert.Equal(t, 4726, item.Length)
	assert.Equal(t, 10229, item.Offset)
	assert.Equal(t, "data.warc.gz", item.Filename)
	assert.Equal(t, "20250402021331", item.Timestamp)
	assert.Equal(t, "io,vito)/assets/css/main.css?cache=20250130t095925", item.Surt)
}

func TestParseAppleDevStreaming(t *testing.T) {
	line := "com,apple,devstreaming-cdn)/videos/wwdc/2021/10061/4/d12f25a4-d409-4dda-9dcf-72c97e9875c3/cmaf.m3u8 20250407234442 {\"metadata\":\"{\\\"adaptive_max_resolution\\\": 921600, \\\"adaptive_max_bandwidth\\\": 2000000}\",\"url\":\"https://devstreaming-cdn.apple.com/videos/wwdc/2021/10061/4/D12F25A4-D409-4DDA-9DCF-72C97E9875C3/cmaf.m3u8\",\"mime\":\"application/vnd.apple.mpegurl\",\"status\":\"200\",\"digest\":\"GZDM6NSUOPMUZ4PG5AVNFFWG5NHFIEFO\",\"length\":\"1859\",\"offset\":\"1277614\",\"filename\":\"data.warc.gz\"}"
	item, err := ItemFromString(line)
	require.NoError(t, err)

	assert.Equal(t, "https://devstreaming-cdn.apple.com/videos/wwdc/2021/10061/4/D12F25A4-D409-4DDA-9DCF-72C97E9875C3/cmaf.m3u8", item.URL)
	assert.Equal(t, "application/vnd.apple.mpegurl", item.MimeType)
	assert.Equal(t, 200, item.Status)
	assert.Equal(t, "GZDM6NSUOPMUZ4PG5AVNFFWG5NHFIEFO", item.Digest)
	assert.Equal(t, 1859, item.Length)
	assert.Equal(t, 1277614, item.Offset)
	assert.Equal(t, "data.warc.gz", item.Filename)
	assert.Equal(t, "20250407234442", item.Timestamp)
	assert.Equal(t, "com,apple,devstreaming-cdn)/videos/wwdc/2021/10061/4/d12f25a4-d409-4dda-9dcf-72c97e9875c3/cmaf.m3u8", item.Surt)
}
