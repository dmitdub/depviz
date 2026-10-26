package apk

// DownloadIndex скачивает APKINDEX.tar.gz по базовому url и возвращает содержимое APKINDEX.
func DownloadIndex(url string) ([]byte, error) {
	return downloadIndex(url)
}
