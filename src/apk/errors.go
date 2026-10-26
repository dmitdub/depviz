// Package apk содержит парсинг индекса пакетов Alpine Linux (APKINDEX).
package apk

import "errors"

// Sentinel-ошибки парсинга APKINDEX.
var (
	// ErrPackageNotFound возвращается, если пакет с заданным именем и версией не найден.
	ErrPackageNotFound = errors.New("apk: пакет не найден")
	// ErrEmptyIndex возвращается, если APKINDEX пуст.
	ErrEmptyIndex = errors.New("apk: пустой APKINDEX")
	// ErrDownloadFailed возвращается, если не удалось скачать APKINDEX.tar.gz.
	ErrDownloadFailed = errors.New("apk: не удалось скачать APKINDEX.tar.gz")
	// ErrUnpackFailed возвращается, если не удалось распаковать tar.gz.
	ErrUnpackFailed = errors.New("apk: не удалось распаковать APKINDEX.tar.gz")
	// ErrIndexEntryNotFound возвращается, если внутри архива нет файла APKINDEX.
	ErrIndexEntryNotFound = errors.New("apk: файл APKINDEX не найден в архиве")
)
