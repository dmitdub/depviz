package apk

import (
	"bufio"
	"bytes"
	"strings"
)

// ResolverFromIndex возвращает Resolver, работающий по APKINDEX.
func ResolverFromIndex(index []byte) func(name string) ([]string, error) {
	cache := make(map[string][]string)
	return func(name string) ([]string, error) {
		if deps, ok := cache[name]; ok {
			return deps, nil
		}
		deps, err := findDepsByName(index, name)
		if err != nil {
			return nil, err
		}
		cache[name] = deps
		return deps, nil
	}
}

// findDepsByName ищет пакет по имени и возвращает его D:.
func findDepsByName(index []byte, name string) ([]string, error) {
	if len(index) == 0 {
		return nil, ErrEmptyIndex
	}
	scanner := bufio.NewScanner(bytes.NewReader(index))
	var currentName string
	var currentDeps []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentName == name {
				return filterDeps(currentDeps), nil
			}
			currentName = ""
			currentDeps = nil
			continue
		}
		switch {
		case strings.HasPrefix(line, "P:"):
			currentName = line[2:]
		case strings.HasPrefix(line, "D:"):
			currentDeps = strings.Fields(line[2:])
		}
	}
	if currentName == name {
		return filterDeps(currentDeps), nil
	}
	return nil, ErrPackageNotFound
}

// filterDeps отбрасывает so:/cmd:-зависимости, абсолютные пути и отрезает версии.
func filterDeps(deps []string) []string {
	out := make([]string, 0, len(deps))
	for _, d := range deps {
		if strings.HasPrefix(d, "so:") ||
			strings.HasPrefix(d, "cmd:") ||
			strings.HasPrefix(d, "/") {
			continue
		}
		out = append(out, stripVersion(d))
	}
	return out
}

// stripVersion отрезает ограничение версии: libcurl=8.14.1-r2 -> libcurl.
func stripVersion(dep string) string {
	if i := strings.IndexAny(dep, "=<>"); i >= 0 {
		return dep[:i]
	}
	return dep
}
