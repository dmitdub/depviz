// Package config отвечает за загрузку и валидацию TOML-конфигурации.
package config

import "errors"

// Sentinel-ошибки валидации конфигурации.
var (
	// ErrEmptyPackageName возвращается, если package_name пуст.
	ErrEmptyPackageName = errors.New("package_name: пустое значение")
	// ErrEmptyRepository возвращается, если repository_url пуст при test_repository=false.
	ErrEmptyRepository = errors.New("repository_url: пустое значение")
	// ErrInvalidDepth возвращается, если max_depth меньше minDepth.
	ErrInvalidDepth = errors.New("max_depth: должно быть >= 1")
	// ErrInvalidTestFile возвращается, если test_repository=true, а test_file_path пуст.
	ErrInvalidTestFile = errors.New("test_file_path: пустой при test_repository=true")
	// ErrEmptyOutputImage возвращается, если output_image пуст.
	ErrEmptyOutputImage = errors.New("output_image: пустое значение")
	// ErrConfigFileNotFound возвращается, если TOML-файл не найден.
	ErrConfigFileNotFound = errors.New("config: файл не найден")
)
