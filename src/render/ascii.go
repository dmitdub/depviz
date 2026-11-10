package render

import (
	"strings"

	"github.com/dmitdub/depviz/src/graph"
)

// ToASCII возвращает ASCII-дерево зависимостей от корня.
func ToASCII(g *graph.Graph) string {
	if len(g.Order) == 0 {
		return ""
	}
	var b strings.Builder
	root := g.Order[0]
	visited := make(map[string]bool)
	writeASCII(&b, g, root, "", true, true, visited)
	return b.String()
}

// writeASCII рекурсивно печатает узел и его зависимости с отступами.
// isRoot — флаг корневого узла (печатается без connector).
// visited предотвращает бесконечную рекурсию при циклах.
func writeASCII(b *strings.Builder, g *graph.Graph, name, prefix string, last, isRoot bool, visited map[string]bool) {
	if visited[name] {
		return
	}
	visited[name] = true

	node, ok := g.Nodes[name]
	if !ok {
		return
	}

	if isRoot {
		b.WriteString(name)
		b.WriteString("\n")
	} else {
		connector := "├── "
		if last {
			connector = "└── "
		}
		b.WriteString(prefix)
		b.WriteString(connector)
		b.WriteString(name)
		b.WriteString("\n")
	}

	childPrefix := prefix
	if !isRoot {
		if last {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, dep := range node.Dependencies {
		writeASCII(b, g, dep, childPrefix, i == len(node.Dependencies)-1, false, visited)
	}
}
