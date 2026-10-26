// Package app содержит основную логику приложения.
package app

import (
	"fmt"
	"io"

	"github.com/dmitdub/depviz/src/config"
	"github.com/dmitdub/depviz/src/graph"
)

// Run выводит все параметры конфигурации в формате key = value.
func Run(cfg *config.Config, w io.Writer) {
	fmt.Fprintln(w, "configuration:")
	for _, kv := range cfg.AsKeyValue() {
		fmt.Fprintf(w, "%-18s = %s\n", kv[0], kv[1])
	}
}

// RunGraph выводит граф зависимостей: узлы, глубину, зависимости.
func RunGraph(g *graph.Graph, w io.Writer) {
	if len(g.Order) == 0 {
		fmt.Fprintln(w, "graph: (empty)")
		return
	}
	fmt.Fprintf(w, "root: %s\n", g.Order[0])
	fmt.Fprintf(w, "nodes: %d\n", len(g.Nodes))
	fmt.Fprintln(w, "graph:")
	for _, name := range g.Order {
		node := g.Nodes[name]
		fmt.Fprintf(w, "  %s (depth=%d): %s\n",
			node.Name, node.Depth, joinDeps(node.Dependencies))
	}
}

// joinDeps склеивает зависимости через запятую.
func joinDeps(deps []string) string {
	if len(deps) == 0 {
		return "(none)"
	}
	out := deps[0]
	for _, d := range deps[1:] {
		out += ", " + d
	}
	return out
}
