package apk

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// indexEntryName — имя файла внутри tar.gz.
const indexEntryName = "APKINDEX"

// indexFileName — имя файла индекса в репозитории.
const indexFileName = "APKINDEX.tar.gz"

// downloadIndex скачивает APKINDEX.tar.gz по базовому url и возвращает содержимое APKINDEX.
func downloadIndex(baseURL string) ([]byte, error) {
	fullURL := baseURL
	if !strings.HasSuffix(fullURL, indexFileName) {
		if !strings.HasSuffix(fullURL, "/") {
			fullURL += "/"
		}
		fullURL += indexFileName
	}

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDownloadFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP %d", ErrDownloadFailed, resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnpackFailed, err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrUnpackFailed, err)
		}
		if hdr.Name == indexEntryName {
			return io.ReadAll(tr)
		}
	}
	return nil, ErrIndexEntryNotFound
}
