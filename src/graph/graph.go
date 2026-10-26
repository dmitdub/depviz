package graph

import "strings"

// Resolver возвращает прямые зависимости пакета по имени.
type Resolver func(name string) ([]string, error)

// Node описывает узел графа зависимостей.
type Node struct {
	Name         string
	Dependencies []string
	Depth        int
	Missing      bool
}

// Graph — результат обхода зависимостей.
type Graph struct {
	Nodes    map[string]*Node
	Order    []string
	MaxDepth int
	Filter   string
}

// Build строит граф зависимостей от root с учётом maxDepth и фильтра.
func Build(root string, maxDepth int, filter string, resolve Resolver) (*Graph, error) {
	g := &Graph{
		Nodes:    make(map[string]*Node),
		MaxDepth: maxDepth,
		Filter:   filter,
	}
	if err := g.visit(root, 0, resolve); err != nil {
		return nil, err
	}
	return g, nil
}

// visit рекурсивно обходит зависимости узла name на глубине depth.
func (g *Graph) visit(name string, depth int, resolve Resolver) error {
	if depth > g.MaxDepth || g.seen(name) || g.isFiltered(name) {
		return nil
	}

	deps, err := resolve(name)
	if err != nil {
		g.Nodes[name] = &Node{Name: name, Depth: depth, Missing: true}
		g.Order = append(g.Order, name)
		return nil
	}

	g.Nodes[name] = &Node{Name: name, Dependencies: deps, Depth: depth}
	g.Order = append(g.Order, name)

	if depth == g.MaxDepth {
		return nil
	}
	for _, dep := range deps {
		if err := g.visit(dep, depth+1, resolve); err != nil {
			return err
		}
	}
	return nil
}

// seen возвращает true, если узел уже добавлен (обработка циклов).
func (g *Graph) seen(name string) bool {
	_, ok := g.Nodes[name]
	return ok
}

// isFiltered возвращает true, если имя содержит подстроку фильтра.
func (g *Graph) isFiltered(name string) bool {
	return g.Filter != "" && strings.Contains(name, g.Filter)
}
