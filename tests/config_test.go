// Package tests содержит тесты для depviz.
package tests

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/dmitdub/depviz/src/config"
)

const (
	// testMaxDepth — глубина, используемая в тестовых конфигах.
	testMaxDepth = 3
	// expectedParams — сколько параметров возвращает Config.AsKeyValue.
	expectedParams = 8
)

// loadErrorCases — таблица кейсов для TestConfigLoadErrors.
var loadErrorCases = []struct {
	name string
	toml string
	want error
}{
	{
		"empty pkg",
		"repository_url = \"x\"\noutput_image = \"o\"\nmax_depth = 1\n",
		config.ErrEmptyPackageName,
	},
	{
		"empty repo",
		"package_name = \"p\"\noutput_image = \"o\"\nmax_depth = 1\n",
		config.ErrEmptyRepository,
	},
	{
		"bad depth",
		"package_name = \"p\"\nrepository_url = \"x\"\noutput_image = \"o\"\nmax_depth = 0\n",
		config.ErrInvalidDepth,
	},
	{
		"empty image",
		"package_name = \"p\"\nrepository_url = \"x\"\nmax_depth = 1\n",
		config.ErrEmptyOutputImage,
	},
	{
		"test no path",
		"package_name = \"p\"\noutput_image = \"o\"\nmax_depth = 1\ntest_repository = true\n",
		config.ErrInvalidTestFile,
	},
}

// writeTemp создаёт временный TOML-файл и возвращает его путь.
func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	return p
}

// TestConfigLoadValid проверяет успешный разбор корректного конфига.
func TestConfigLoadValid(t *testing.T) {
	p := writeTemp(t, `
package_name     = "curl"
package_version  = "8.5.0-r0"
repository_url   = "https://example.org/repo"
test_repository  = false
output_image     = "out.svg"
max_depth        = 3
filter_substring = "dev"
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.PackageName != "curl" || cfg.MaxDepth != testMaxDepth {
		t.Fatalf("bad parse: %+v", cfg)
	}
	if got := len(cfg.AsKeyValue()); got != expectedParams {
		t.Fatalf("ожидалось %d параметров, получено %d", expectedParams, got)
	}
}

// TestConfigLoadErrors проверяет, что каждая ошибка валидации возвращается.
func TestConfigLoadErrors(t *testing.T) {
	for _, tc := range loadErrorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := config.Load(writeTemp(t, tc.toml))
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}

// TestConfigLoadMissingFile проверяет ошибку отсутствия файла.
func TestConfigLoadMissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nope.toml"))
	if !errors.Is(err, config.ErrConfigFileNotFound) {
		t.Fatalf("want ErrConfigFileNotFound, got %v", err)
	}
}
