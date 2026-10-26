package apk

import (
	"bufio"
	"bytes"
	"strings"
)

// Package описывает запись APKINDEX.
type Package struct {
	Name         string
	Version      string
	Dependencies []string
}

// FindPackage ищет пакет по имени и версии в APKINDEX и возвращает его запись.
func FindPackage(index []byte, name, version string) (*Package, error) {
	if len(index) == 0 {
		return nil, ErrEmptyIndex
	}

	scanner := bufio.NewScanner(bytes.NewReader(index))
	var current *Package
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if matches(current, name, version) {
				return current, nil
			}
			current = nil
			continue
		}

		if current == nil {
			current = &Package{}
		}
		applyField(current, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if matches(current, name, version) {
		return current, nil
	}
	return nil, ErrPackageNotFound
}

// matches проверяет, что запись соответствует имени и версии.
func matches(p *Package, name, version string) bool {
	return p != nil && p.Name == name && p.Version == version
}

// applyField разбирает одну строку APKINDEX и обновляет запись.
func applyField(p *Package, line string) {
	switch {
	case strings.HasPrefix(line, "P:"):
		p.Name = line[2:]
	case strings.HasPrefix(line, "V:"):
		p.Version = line[2:]
	case strings.HasPrefix(line, "D:"):
		p.Dependencies = strings.Fields(line[2:])
	}
}
