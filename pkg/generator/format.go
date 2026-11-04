package generator

import (
	"fmt"
	"go/format"
	"path/filepath"
)

// formatGo runs gofmt over rendered content when path is a Go file, so that
// generated code is always gofmt-clean regardless of template whitespace.
// Non-Go files are returned unchanged.
func formatGo(path string, content []byte) ([]byte, error) {
	if filepath.Ext(path) != ".go" {
		return content, nil
	}
	out, err := format.Source(content)
	if err != nil {
		return nil, fmt.Errorf("failed to format %s: %w", path, err)
	}
	return out, nil
}
