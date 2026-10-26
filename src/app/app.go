// Package app содержит основную логику приложения.
package app

import (
	"fmt"
	"io"

	"github.com/dmitdub/depviz/src/apk"
	"github.com/dmitdub/depviz/src/config"
)

// Run выводит все параметры конфигурации в формате key = value.
func Run(cfg *config.Config, w io.Writer) {
	fmt.Fprintln(w, "configuration:")
	for _, kv := range cfg.AsKeyValue() {
		fmt.Fprintf(w, "%-18s = %s\n", kv[0], kv[1])
	}
}

// RunDeps выводит прямые зависимости пакета pkg.
func RunDeps(pkg *apk.Package, w io.Writer) {
	fmt.Fprintf(w, "package: %s %s\n", pkg.Name, pkg.Version)
	if len(pkg.Dependencies) == 0 {
		fmt.Fprintln(w, "dependencies: (none)")
		return
	}
	fmt.Fprintln(w, "dependencies:")
	for _, dep := range pkg.Dependencies {
		fmt.Fprintf(w, "  %s\n", dep)
	}
}
