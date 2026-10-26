package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// minDepth — минимально допустимая глубина анализа.
const minDepth = 1

// Config описывает настраиваемые пользователем параметры.
type Config struct {
	PackageName     string `toml:"package_name"`
	PackageVersion  string `toml:"package_version"`
	RepositoryURL   string `toml:"repository_url"`
	TestRepository  bool   `toml:"test_repository"`
	TestFilePath    string `toml:"test_file_path"`
	OutputImage     string `toml:"output_image"`
	MaxDepth        int    `toml:"max_depth"`
	FilterSubstring string `toml:"filter_substring"`
}

// Load читает TOML-файл по пути path и возвращает провалидированную конфигурацию.
func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("config: ошибка разбора TOML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate проверяет все параметры и возвращает первую найденную ошибку.
func (c *Config) Validate() error {
	if c.PackageName == "" {
		return ErrEmptyPackageName
	}
	if c.RepositoryURL == "" && !c.TestRepository {
		return ErrEmptyRepository
	}
	if c.MaxDepth < minDepth {
		return ErrInvalidDepth
	}
	if c.OutputImage == "" {
		return ErrEmptyOutputImage
	}
	if c.TestRepository && c.TestFilePath == "" {
		return ErrInvalidTestFile
	}
	if c.TestRepository {
		c.TestFilePath = filepath.Clean(c.TestFilePath)
	}
	return nil
}

// AsKeyValue возвращает параметры в виде пар "ключ=значение".
func (c *Config) AsKeyValue() [][2]string {
	return [][2]string{
		{"package.name", c.PackageName},
		{"package.version", c.PackageVersion},
		{"repository.url", c.RepositoryURL},
		{"test.repository", fmt.Sprintf("%v", c.TestRepository)},
		{"test.file.path", c.TestFilePath},
		{"output.image", c.OutputImage},
		{"max.depth", fmt.Sprintf("%d", c.MaxDepth)},
		{"filter.substring", c.FilterSubstring},
	}
}
