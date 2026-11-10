// Package main — точка входа depviz.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/depviz/src/apk"
	"github.com/dmitdub/depviz/src/app"
	"github.com/dmitdub/depviz/src/config"
	"github.com/dmitdub/depviz/src/graph"
	"github.com/dmitdub/depviz/src/render"
)

// Коды возврата процесса.
const (
	exitOK    = 0
	exitError = 1
)

// defaultConfigPath — путь к конфигу по умолчанию.
const defaultConfigPath = "testdata/config.example.toml"

// defaultPlantUMLJar — путь к plantuml.jar по умолчанию.
const defaultPlantUMLJar = "bin/plantuml.jar"

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

	g, err := buildGraph(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка построения графа: %v\n", err)
		return exitError
	}

	app.RunGraph(g, stderr)

	if err := runRender(cfg, g, stderr); err != nil {
		fmt.Fprintf(stderr, "Ошибка визуализации: %v\n", err)
		return exitError
	}
	return exitOK
}

// buildGraph строит прямой или обратный граф в зависимости от cfg.Reverse.
func buildGraph(cfg *config.Config) (*graph.Graph, error) {
	if cfg.Reverse {
		resolve, err := buildReverseResolver(cfg)
		if err != nil {
			return nil, err
		}
		return graph.Reverse(cfg.PackageName, cfg.MaxDepth, cfg.FilterSubstring, resolve)
	}

	resolve, err := buildResolver(cfg)
	if err != nil {
		return nil, err
	}
	return graph.Build(cfg.PackageName, cfg.MaxDepth, cfg.FilterSubstring, resolve)
}

// buildResolver выбирает источник прямых зависимостей: тестовый файл или APKINDEX.
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

// buildReverseResolver выбирает источник обратных зависимостей.
func buildReverseResolver(cfg *config.Config) (graph.Resolver, error) {
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
	return apk.ReverseResolverFromIndex(index), nil
}

// runRender генерирует PlantUML-текст, SVG-файл и ASCII-дерево.
func runRender(cfg *config.Config, g *graph.Graph, stderr *os.File) error {
	puml := render.ToPlantUML(g)
	fmt.Fprintln(stderr, puml)

	if err := render.WriteSVG(puml, cfg.OutputImage, defaultPlantUMLJar); err != nil {
		return err
	}
	fmt.Fprintf(stderr, "SVG сохранён: %s\n", cfg.OutputImage)

	if cfg.ASCIITree {
		fmt.Fprint(stderr, render.ToASCII(g))
	}
	return nil
}
