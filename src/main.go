// Package main — точка входа depviz.
package main

import (
	"fmt"
	"os"

	"github.com/dmitdub/depviz/src/app"
	"github.com/dmitdub/depviz/src/config"
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

// run выполняет разбор аргументов, загрузку конфигурации и вывод.
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

	app.Run(cfg, stderr)
	return exitOK
}
