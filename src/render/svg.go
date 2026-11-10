package render

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// WriteSVG сохраняет SVG-изображение графа в файле path.
// Использует plantuml.jar через java -jar.
func WriteSVG(puml string, path string, plantumlJar string) error {
	tmpDir, err := os.MkdirTemp("", "depviz-render-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	pumlPath := filepath.Join(tmpDir, "graph.puml")
	if err := os.WriteFile(pumlPath, []byte(puml), 0o600); err != nil {
		return err
	}

	cmd := exec.Command("java", "-jar", plantumlJar, "-tsvg", pumlPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", ErrPlantUMLFailed, stderr.String())
	}

	svgPath := filepath.Join(tmpDir, "graph.svg")
	data, err := os.ReadFile(svgPath)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
