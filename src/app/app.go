// Package app содержит основную логику приложения.
package app

import (
	"fmt"
	"io"

	"github.com/dmitdub/depviz/src/config"
)

// Run выводит все параметры конфигурации в формате key = value.
func Run(cfg *config.Config, w io.Writer) {
	fmt.Fprintln(w, "configuration:")
	for _, kv := range cfg.AsKeyValue() {
		fmt.Fprintf(w, "%-18s = %s\n", kv[0], kv[1])
	}
}
