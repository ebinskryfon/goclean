package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ebinskryfon/goclean/pkg/strutil"
	"github.com/ebinskryfon/goclean/pkg/ui"
	"github.com/ebinskryfon/goclean/templates"
)

type ProjectOptions struct {
	ProjectName string
	ModuleName  string
	DBName      string
	TargetDir   string
}

// GenerateProject scaffolds a new Go Clean Architecture project from embedded templates.
func GenerateProject(opts ProjectOptions) error {
	if opts.DBName == "" {
		opts.DBName = strutil.ToSnakeCase(opts.ProjectName) + "_db"
	}
	if opts.ModuleName == "" {
		opts.ModuleName = opts.ProjectName
	}

	// Check if target directory exists and is non-empty
	if stat, err := os.Stat(opts.TargetDir); err == nil && stat.IsDir() {
		entries, err := os.ReadDir(opts.TargetDir)
		if err == nil && len(entries) > 0 {
			return fmt.Errorf("directory %s already exists and is not empty", opts.TargetDir)
		}
	}

	ui.Step(1, 4, "Scaffolding Clean Architecture files")

	err := fs.WalkDir(templates.FS, "project", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("project", path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		// Destination path (strip .tmpl suffix if present)
		cleanRelPath := strings.TrimSuffix(relPath, ".tmpl")
		destPath := filepath.Join(opts.TargetDir, cleanRelPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// Render template
		tmpl, err := template.New(relPath).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, opts); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", path, err)
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, buf.Bytes(), 0644)
	})

	if err != nil {
		return err
	}

	// Create .env from .env.example
	ui.Step(2, 4, "Configuring environment (.env)")
	envExamplePath := filepath.Join(opts.TargetDir, ".env.example")
	envPath := filepath.Join(opts.TargetDir, ".env")
	if envBytes, err := os.ReadFile(envExamplePath); err == nil {
		_ = os.WriteFile(envPath, envBytes, 0644)
	}

	// Initialize git repository
	ui.Step(3, 4, "Initializing Git repository")
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = opts.TargetDir
	_ = gitCmd.Run()

	// Run go mod tidy
	ui.Step(4, 4, "Downloading and tidying dependencies (go mod tidy)")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = opts.TargetDir
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		ui.Warn("go mod tidy warning (you can run 'go mod tidy' later manually): %s", string(tidyOutput))
	}

	return nil
}
