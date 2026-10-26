// Package graph строит граф зависимостей и обходит его.
package graph

import "errors"

// Sentinel-ошибки построения графа.
var (
	// ErrPackageNotFound возвращается, если пакет отсутствует в репозитории.
	ErrPackageNotFound = errors.New("graph: пакет не найден")
	// ErrEmptyTestRepo возвращается, если тестовый файл пуст.
	ErrEmptyTestRepo = errors.New("graph: пустой тестовый репозиторий")
)
