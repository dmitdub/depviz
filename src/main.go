// Package main — точка входа depviz.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/depviz/src/apk"
	"github.com/dmitdub/depviz/src/app"
	"github.com/dmitdub/depviz/src/config"
	"github.com/dmitdub/depviz/src/graph"
)

// Коды возврата процесса.
const (
	exitOK    = 0
	exitError = 1
)

// defaultConfigPath — путь к конфигу по умолчанию.
const defaultConfigPath = "testdata/config.example.toml"

// main запускает приложение и завершает процесс с соответствующим кодом.
func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

// run выполняет разбор аргументов, загрузку конфигурации и построение графа.
func run(args []string, stderr *os.File) int {
	path := defaultConfigPath
	if len(args) > 0 {
		path = args[0]
	}

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка конфигурации: %v\n", err)
		return exitError
	}

	resolve, err := buildResolver(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка подготовки репозитория: %v\n", err)
		return exitError
	}

	g, err := graph.Build(cfg.PackageName, cfg.MaxDepth, cfg.FilterSubstring, resolve)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка построения графа: %v\n", err)
		return exitError
	}

	app.RunGraph(g, stderr)
	return exitOK
}

// buildResolver выбирает источник зависимостей: тестовый файл или APKINDEX.
func buildResolver(cfg *config.Config) (graph.Resolver, error) {
	if cfg.TestRepository {
		data, err := os.ReadFile(cfg.TestFilePath)
		if err != nil {
			return nil, err
		}
		repo, err := graph.ParseTestRepo(data)
		if err != nil {
			return nil, err
		}
		return graph.ResolverFromRepo(repo), nil
	}

	index, err := apk.DownloadIndex(cfg.RepositoryURL)
	if err != nil {
		return nil, err
	}
	return apk.ResolverFromIndex(index), nil
}
