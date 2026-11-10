package apk

import (
	"bufio"
	"bytes"
	"strings"
)

// ReverseResolverFromIndex возвращает Resolver, который по имени пакета
// возвращает список пакетов, зависящих от него (обратные зависимости).
func ReverseResolverFromIndex(index []byte) func(name string) ([]string, error) {
	cache := make(map[string][]string)
	return func(name string) ([]string, error) {
		if deps, ok := cache[name]; ok {
			return deps, nil
		}
		deps, err := findReverseDeps(index, name)
		if err != nil {
			return nil, err
		}
		cache[name] = deps
		return deps, nil
	}
}

// findReverseDeps проходит весь индекс и собирает пакеты, у которых name в D:.
func findReverseDeps(index []byte, name string) ([]string, error) {
	if len(index) == 0 {
		return nil, ErrEmptyIndex
	}
	scanner := bufio.NewScanner(bytes.NewReader(index))
	var currentPkg string
	out := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			currentPkg = ""
			continue
		}
		switch {
		case strings.HasPrefix(line, "P:"):
			currentPkg = line[2:]
		case strings.HasPrefix(line, "D:"):
			for _, d := range strings.Fields(line[2:]) {
				if stripVersion(d) == name {
					out = append(out, currentPkg)
					break
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
