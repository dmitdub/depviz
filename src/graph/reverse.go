package graph

// Reverse строит обратный граф: кто зависит от root.
// Обход рекурсивный, с учётом maxDepth, filter и обработкой циклов.
func Reverse(root string, maxDepth int, filter string, reverseResolve Resolver) (*Graph, error) {
	g := &Graph{
		Nodes:    make(map[string]*Node),
		MaxDepth: maxDepth,
		Filter:   filter,
	}
	if err := g.visitReverse(root, 0, reverseResolve); err != nil {
		return nil, err
	}
	return g, nil
}

// visitReverse рекурсивно обходит обратные зависимости узла name на глубине depth.
func (g *Graph) visitReverse(name string, depth int, resolve Resolver) error {
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
		if err := g.visitReverse(dep, depth+1, resolve); err != nil {
			return err
		}
	}
	return nil
}
