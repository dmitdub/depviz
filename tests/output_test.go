package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dmitdub/depviz/src/app"
	"github.com/dmitdub/depviz/src/config"
)

// TestConfigOutput проверяет, что все параметры конфигурации
// попадают в вывод в формате key = value.
func TestConfigOutput(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.toml")
	content := `
package_name     = "curl"
repository_url   = "https://example.org/repo"
test_repository  = false
output_image     = "out.svg"
max_depth        = 3
filter_substring = "dev"
`
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	var buf bytes.Buffer
	app.Run(cfg, &buf)

	out := buf.String()
	keys := []string{
		"package.name", "repository.url", "test.repository",
		"test.file.path", "output.image", "max.depth", "filter.substring",
	}
	for _, key := range keys {
		if !strings.Contains(out, key) {
			t.Errorf("вывод не содержит ключ %q:\n%s", key, out)
		}
	}
	if !strings.Contains(out, "configuration:") {
		t.Errorf("вывод не содержит заголовок configuration:\n%s", out)
	}
}
