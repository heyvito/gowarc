package warc

import (
	"github.com/heyvito/gowarc/cdxj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWarcIO(t *testing.T) {
	indexData, err := os.ReadFile("fixtures/autoindex.cdxj")
	require.NoError(t, err)
	index, err := cdxj.Parse(indexData)
	require.NoError(t, err)
	root := index.Find(func(i *cdxj.Item) bool {
		return i.URL == "https://vito.io/"
	})
	require.NotNil(t, root)

	archiveData, err := os.ReadFile("fixtures/data.warc.gz")
	require.NoError(t, err)
	warcio, err := NewIO(archiveData)
	require.NoError(t, err)
	parsed, err := warcio.ReadIndexedRecord(root)
	require.NoError(t, err)
	body, err := parsed.GetBody()
	require.NoError(t, err)
	assert.Contains(t, string(body), "Vito Sartori")
}

func TestWarcIOFile2(t *testing.T) {
	indexData, err := os.ReadFile("fixtures/autoindex.cdxj")
	require.NoError(t, err)
	index, err := cdxj.Parse(indexData)
	require.NoError(t, err)
	root := index.Find(func(i *cdxj.Item) bool {
		return i.URL == "https://vito.io/assets/css/main.css?cache=20250130T095925"
	})
	require.NotNil(t, root)
	archiveData, err := os.ReadFile("fixtures/data.warc.gz")
	require.NoError(t, err)
	warcio, err := NewIO(archiveData)
	require.NoError(t, err)
	parsed, err := warcio.ReadIndexedRecord(root)
	require.NoError(t, err)

	contentType := parsed.ContentType
	require.NotNil(t, contentType)
	assert.Equal(t, "text/css", *contentType)

	body, err := parsed.GetBody()
	require.NoError(t, err)
	assert.Contains(t, string(body), "*, *::before, *::after")
}
