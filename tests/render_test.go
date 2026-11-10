package tests

import (
	"strings"
	"testing"

	"github.com/dmitdub/depviz/src/graph"
	"github.com/dmitdub/depviz/src/render"
)

const (
	// renderNodes — сколько узлов в тестовом графе.
	renderNodes = 3
)

// buildRenderGraph строит небольшой граф для тестов рендера.
func buildRenderGraph(t *testing.T) *graph.Graph {
	t.Helper()
	repo := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}
	g, err := graph.Build("A", 5, "", graph.ResolverFromRepo(repo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	return g
}

// TestToPlantUML проверяет, что PlantUML-текст содержит заголовок и стрелки.
func TestToPlantUML(t *testing.T) {
	g := buildRenderGraph(t)
	puml := render.ToPlantUML(g)

	if !strings.Contains(puml, "@startuml") {
		t.Fatalf("нет @startuml:\n%s", puml)
	}
	if !strings.Contains(puml, "@enduml") {
		t.Fatalf("нет @enduml:\n%s", puml)
	}
	if !strings.Contains(puml, "\"A\" --> \"B\"") {
		t.Fatalf("нет стрелки A --> B:\n%s", puml)
	}
	if !strings.Contains(puml, "\"B\" --> \"C\"") {
		t.Fatalf("нет стрелки B --> C:\n%s", puml)
	}
}

// TestToASCII проверяет, что ASCII-дерево содержит корень и потомков.
func TestToASCII(t *testing.T) {
	g := buildRenderGraph(t)
	tree := render.ToASCII(g)

	if !strings.HasPrefix(tree, "A") {
		t.Fatalf("дерево не начинается с A:\n%s", tree)
	}
	if !strings.Contains(tree, "B") {
		t.Fatalf("нет B:\n%s", tree)
	}
	if !strings.Contains(tree, "C") {
		t.Fatalf("нет C:\n%s", tree)
	}
}

// TestToASCIIEmpty проверяет пустой граф.
func TestToASCIIEmpty(t *testing.T) {
	g := &graph.Graph{Nodes: map[string]*graph.Node{}, Order: []string{}}
	if render.ToASCII(g) != "" {
		t.Fatalf("ожидалась пустая строка")
	}
}
