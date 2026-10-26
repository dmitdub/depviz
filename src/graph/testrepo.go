package graph

import (
	"bufio"
	"bytes"
	"strings"
)

// ParseTestRepo разбирает тестовый репозиторий из байтов.
// Формат строки: "ИМЯ: ЗАВ1 ЗАВ2 ...".
func ParseTestRepo(data []byte) (map[string][]string, error) {
	if len(data) == 0 {
		return nil, ErrEmptyTestRepo
	}

	repo := make(map[string][]string)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, deps, ok := parseLine(line)
		if !ok {
			continue
		}
		repo[name] = deps
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return repo, nil
}

// parseLine разбирает одну строку "ИМЯ: ЗАВ1 ЗАВ2".
func parseLine(line string) (string, []string, bool) {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return "", nil, false
	}
	name := strings.TrimSpace(line[:idx])
	if name == "" {
		return "", nil, false
	}
	rest := strings.TrimSpace(line[idx+1:])
	if rest == "" {
		return name, nil, true
	}
	return name, strings.Fields(rest), true
}

// ResolverFromRepo возвращает Resolver, работающий по карте тестового репозитория.
func ResolverFromRepo(repo map[string][]string) Resolver {
	return func(name string) ([]string, error) {
		deps, ok := repo[name]
		if !ok {
			return nil, ErrPackageNotFound
		}
		return deps, nil
	}
}
