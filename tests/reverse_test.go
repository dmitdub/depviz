package tests

import (
	"testing"

	"github.com/dmitdub/depviz/src/graph"
)

const (
	// reverseDepth — глубина для проверки обратного обхода.
	reverseDepth = 1
	// reverseNodes — сколько узлов в обратном графе от E на глубине 1.
	reverseNodes = 3
	// reverseSingleNode — один узел (для TestReverseGraphNoDeps).
	reverseSingleNode = 1
)

// reverseTestRepo — тестовый репозиторий для обратных зависимостей.
var reverseTestRepo = map[string][]string{
	"A": {"B", "C"},
	"B": {"D"},
	"C": {"D", "E"},
	"D": {"E"},
	"E": {"A"},
}

// reverseResolver возвращает Resolver, работающий по карте обратных зависимостей.
func reverseResolver(repo map[string][]string) graph.Resolver {
	return func(name string) ([]string, error) {
		deps, ok := repo[name]
		if !ok {
			return nil, graph.ErrPackageNotFound
		}
		return deps, nil
	}
}

// TestReverseGraph проверяет обратный обход на тестовом репозитории.
func TestReverseGraph(t *testing.T) {
	reverseRepo := map[string][]string{
		"A": {"E"},
		"B": {"A"},
		"C": {"A"},
		"D": {"B", "C"},
		"E": {"C", "D"},
	}
	g, err := graph.Reverse("E", reverseDepth, "", reverseResolver(reverseRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(g.Nodes) != reverseNodes {
		t.Fatalf("want %d nodes, got %d", reverseNodes, len(g.Nodes))
	}
}

// TestReverseGraphFilter проверяет фильтр в обратном обходе.
func TestReverseGraphFilter(t *testing.T) {
	reverseRepo := map[string][]string{
		"A": {"E"},
		"B": {"A"},
	}
	g, err := graph.Reverse("E", reverseDepth, "B", reverseResolver(reverseRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, ok := g.Nodes["B"]; ok {
		t.Fatalf("узел B должен быть отфильтрован")
	}
}

// TestReverseGraphMissing проверяет, что отсутствующий пакет помечается как missing.
func TestReverseGraphMissing(t *testing.T) {
	g, err := graph.Reverse("Z", reverseDepth, "", reverseResolver(reverseTestRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	node, ok := g.Nodes["Z"]
	if !ok {
		t.Fatalf("want node Z")
	}
	if !node.Missing {
		t.Fatalf("want Z marked as missing")
	}
}

// TestReverseGraphDepth проверяет ограничение глубины в обратном обходе.
func TestReverseGraphDepth(t *testing.T) {
	reverseRepo := map[string][]string{
		"A": {"E"},
		"B": {"A"},
		"C": {"B"},
	}
	g, err := graph.Reverse("E", reverseDepth, "", reverseResolver(reverseRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, ok := g.Nodes["C"]; ok {
		t.Fatalf("узел C не должен попасть на глубину %d", reverseDepth)
	}
}

// TestReverseGraphNoDeps проверяет пакет без обратных зависимостей.
func TestReverseGraphNoDeps(t *testing.T) {
	reverseRepo := map[string][]string{
		"A": {},
	}
	g, err := graph.Reverse("A", reverseDepth, "", reverseResolver(reverseRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(g.Nodes) != reverseSingleNode {
		t.Fatalf("want %d node, got %d", reverseSingleNode, len(g.Nodes))
	}
}
