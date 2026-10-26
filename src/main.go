// Package main — точка входа depviz.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/depviz/src/apk"
	"github.com/dmitdub/depviz/src/app"
	"github.com/dmitdub/depviz/src/config"
)

const (
	exitOK    = 0
	exitError = 1
)

const defaultConfigPath = "testdata/config.example.toml"

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

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

	if cfg.TestRepository {
		fmt.Fprintln(stderr, "режим тестового репозитория пока не реализован")
		return exitError
	}

	index, err := apk.DownloadIndex(cfg.RepositoryURL)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка загрузки APKINDEX: %v\n", err)
		return exitError
	}

	pkg, err := apk.FindPackage(index, cfg.PackageName, cfg.PackageVersion)
	if err != nil {
		fmt.Fprintf(stderr, "Ошибка поиска пакета: %v\n", err)
		return exitError
	}

	app.RunDeps(pkg, stderr)
	return exitOK
}
