package tests

import (
	"errors"
	"strings"
	"testing"

	"github.com/dmitdub/depviz/src/apk"
)

// testIndex — минимальный APKINDEX для тестов.
const testIndex = `C:Q1aaa
P:curl
V:8.5.0-r0
A:x86_64
D:libcurl so:libc.musl-x86_64.so.1

C:Q1bbb
P:libcurl
V:8.5.0-r0
D:openssl zlib

C:Q1ccc
P:vim
V:9.0.0-r0
D:
`

// TestFindPackage проверяет поиск пакета по имени и версии.
func TestFindPackage(t *testing.T) {
	pkg, err := apk.FindPackage([]byte(testIndex), "curl", "8.5.0-r0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if pkg.Name != "curl" {
		t.Fatalf("want curl, got %s", pkg.Name)
	}
	if len(pkg.Dependencies) != 2 {
		t.Fatalf("want 2 deps, got %d", len(pkg.Dependencies))
	}
}

// TestFindPackageNotFound проверяет ошибку, если пакет не найден.
func TestFindPackageNotFound(t *testing.T) {
	_, err := apk.FindPackage([]byte(testIndex), "nope", "1.0")
	if !errors.Is(err, apk.ErrPackageNotFound) {
		t.Fatalf("want ErrPackageNotFound, got %v", err)
	}
}

// TestFindPackageEmptyIndex проверяет ошибку пустого индекса.
func TestFindPackageEmptyIndex(t *testing.T) {
	_, err := apk.FindPackage(nil, "curl", "8.5.0-r0")
	if !errors.Is(err, apk.ErrEmptyIndex) {
		t.Fatalf("want ErrEmptyIndex, got %v", err)
	}
}

// TestFindPackageNoDeps проверяет пакет без зависимостей.
func TestFindPackageNoDeps(t *testing.T) {
	pkg, err := apk.FindPackage([]byte(testIndex), "vim", "9.0.0-r0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(pkg.Dependencies) != 0 {
		t.Fatalf("want 0 deps, got %d", len(pkg.Dependencies))
	}
}

// TestFindPackageVersionMatch проверяет, что версия учитывается.
func TestFindPackageVersionMatch(t *testing.T) {
	_, err := apk.FindPackage([]byte(testIndex), "curl", "9.0.0-r0")
	if !errors.Is(err, apk.ErrPackageNotFound) {
		t.Fatalf("want ErrPackageNotFound, got %v", err)
	}
}

// TestDependenciesParsed проверяет, что D: разбивается по пробелам.
func TestDependenciesParsed(t *testing.T) {
	pkg, err := apk.FindPackage([]byte(testIndex), "libcurl", "8.5.0-r0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []string{"openssl", "zlib"}
	if len(pkg.Dependencies) != len(want) {
		t.Fatalf("want %d deps, got %d", len(want), len(pkg.Dependencies))
	}
	for i, w := range want {
		if pkg.Dependencies[i] != w {
			t.Fatalf("dep[%d]: want %s, got %s", i, w, pkg.Dependencies[i])
		}
	}
}

// TestDependenciesWithVersionConstraint проверяет зависимость с ограничением версии.
func TestDependenciesWithVersionConstraint(t *testing.T) {
	index := "P:foo\nV:1.0\nD:libc>=2.0 so:libc.so.1\n\n"
	pkg, err := apk.FindPackage([]byte(index), "foo", "1.0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(pkg.Dependencies) != 2 {
		t.Fatalf("want 2 deps, got %d", len(pkg.Dependencies))
	}
	if pkg.Dependencies[0] != "libc>=2.0" {
		t.Fatalf("want libc>=2.0, got %s", pkg.Dependencies[0])
	}
}

// TestDependenciesEmptyD проверяет, что пустое D: даёт пустой список.
func TestDependenciesEmptyD(t *testing.T) {
	index := "P:foo\nV:1.0\nD:\n\n"
	pkg, err := apk.FindPackage([]byte(index), "foo", "1.0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(pkg.Dependencies) != 0 {
		t.Fatalf("want 0 deps, got %d", len(pkg.Dependencies))
	}
}

// TestScannerNoTrailingNewline проверяет индекс без пустой строки в конце.
func TestScannerNoTrailingNewline(t *testing.T) {
	index := "P:foo\nV:1.0\nD:bar"
	pkg, err := apk.FindPackage([]byte(index), "foo", "1.0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(pkg.Dependencies) != 1 || pkg.Dependencies[0] != "bar" {
		t.Fatalf("want [bar], got %v", pkg.Dependencies)
	}
}

// TestStringsFieldsUsed проверяет, что лишние пробелы не ломают разбор.
func TestStringsFieldsUsed(t *testing.T) {
	index := "P:foo\nV:1.0\nD:  a   b  c  \n\n"
	pkg, err := apk.FindPackage([]byte(index), "foo", "1.0")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if strings.Join(pkg.Dependencies, ",") != "a,b,c" {
		t.Fatalf("want a,b,c, got %v", pkg.Dependencies)
	}
}
