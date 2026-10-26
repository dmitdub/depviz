package tests

import (
	"errors"
	"testing"

	"github.com/dmitdub/depviz/src/graph"
)

const (
	// fullDepth — глубина, достаточная для полного обхода testRepo.
	fullDepth = 5
	// depthOne — глубина 1.
	depthOne = 1
	// deepDepth — заведомо большая глубина для проверки цикла.
	deepDepth = 10
	// totalNodes — сколько узлов в testRepo.
	totalNodes = 5
	// nodesAtDepthOne — сколько узлов достигается от A на глубине 1.
	nodesAtDepthOne = 3
	// repoEntries — сколько записей в разборе тестового файла.
	repoEntries = 3
	// aDeps — сколько зависимостей у A.
	aDeps = 2
)

// testRepo — тестовый репозиторий с циклом.
var testRepo = map[string][]string{
	"A": {"B", "C"},
	"B": {"D"},
	"C": {"D", "E"},
	"D": {"E"},
	"E": {"A"},
}

// TestBuildGraph проверяет базовый обход.
func TestBuildGraph(t *testing.T) {
	g, err := graph.Build("A", fullDepth, "", graph.ResolverFromRepo(testRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(g.Nodes) != totalNodes {
		t.Fatalf("want %d nodes, got %d", totalNodes, len(g.Nodes))
	}
}

// TestBuildGraphDepth проверяет ограничение глубины.
func TestBuildGraphDepth(t *testing.T) {
	g, err := graph.Build("A", depthOne, "", graph.ResolverFromRepo(testRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(g.Nodes) != nodesAtDepthOne {
		t.Fatalf("want %d nodes at depth %d, got %d",
			nodesAtDepthOne, depthOne, len(g.Nodes))
	}
}

// TestBuildGraphFilter проверяет фильтрацию по подстроке.
func TestBuildGraphFilter(t *testing.T) {
	g, err := graph.Build("A", fullDepth, "E", graph.ResolverFromRepo(testRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, ok := g.Nodes["E"]; ok {
		t.Fatalf("узел E должен быть отфильтрован")
	}
}

// TestBuildGraphCycle проверяет, что цикл не вызывает бесконечный обход.
func TestBuildGraphCycle(t *testing.T) {
	g, err := graph.Build("A", deepDepth, "", graph.ResolverFromRepo(testRepo))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(g.Nodes) != totalNodes {
		t.Fatalf("want %d nodes, got %d", totalNodes, len(g.Nodes))
	}
}

// TestBuildGraphNotFound проверяет, что отсутствующий пакет помечается как missing.
func TestBuildGraphNotFound(t *testing.T) {
	g, err := graph.Build("Z", fullDepth, "", graph.ResolverFromRepo(testRepo))
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

// TestParseTestRepo проверяет разбор тестового файла.
func TestParseTestRepo(t *testing.T) {
	data := []byte("A: B C\nB: D\nC:\n")
	repo, err := graph.ParseTestRepo(data)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(repo) != repoEntries {
		t.Fatalf("want %d entries, got %d", repoEntries, len(repo))
	}
	if len(repo["A"]) != aDeps {
		t.Fatalf("want %d deps for A, got %d", aDeps, len(repo["A"]))
	}
	if len(repo["C"]) != 0 {
		t.Fatalf("want 0 deps for C, got %d", len(repo["C"]))
	}
}

// TestParseTestRepoEmpty проверяет ошибку пустого файла.
func TestParseTestRepoEmpty(t *testing.T) {
	_, err := graph.ParseTestRepo(nil)
	if !errors.Is(err, graph.ErrEmptyTestRepo) {
		t.Fatalf("want ErrEmptyTestRepo, got %v", err)
	}
}
