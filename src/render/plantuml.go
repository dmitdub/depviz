package render

import (
	"fmt"
	"strings"

	"github.com/dmitdub/depviz/src/graph"
)

// ToPlantUML возвращает текстовое представление графа на языке PlantUML.
// Узлы без зависимостей не выводятся отдельными строками — они и так
// появляются как цели стрелок.
func ToPlantUML(g *graph.Graph) string {
	var b strings.Builder
	b.WriteString("@startuml\n")
	for _, name := range g.Order {
		node := g.Nodes[name]
		for _, dep := range node.Dependencies {
			fmt.Fprintf(&b, "\"%s\" --> \"%s\"\n", node.Name, dep)
		}
	}
	b.WriteString("@enduml\n")
	return b.String()
}
