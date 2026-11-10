// Package render генерирует текстовое и графическое представление графа.
package render

import "errors"

// Sentinel-ошибки рендера.
var (
	// ErrPlantUMLFailed возвращается, если plantuml завершился с ошибкой.
	ErrPlantUMLFailed = errors.New("render: plantuml завершился с ошибкой")
)
