package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ebinskryfon/goclean/pkg/strutil"
	"github.com/ebinskryfon/goclean/pkg/ui"
	"github.com/ebinskryfon/goclean/templates"
)

type StubData struct {
	ModuleName         string
	EntityPascal       string // e.g. "Product"
	EntityCamel        string // e.g. "product"
	EntitySnake        string // e.g. "product"
	EntityKebab        string // e.g. "product"
	EntityPluralSnake  string // e.g. "products"
	EntityPluralPascal string // e.g. "Products"
}

// FindGoModuleName locates the nearest go.mod and extracts the module path.
func FindGoModuleName(startDir string) (string, string, error) {
	curr := startDir
	for {
		goModPath := filepath.Join(curr, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			file, err := os.Open(goModPath)
			if err != nil {
				return "", "", err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "module ") {
					modName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
					return modName, curr, nil
				}
			}
			return "", "", fmt.Errorf("module line not found in %s", goModPath)
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", "", fmt.Errorf("go.mod not found in %s or any parent directory", startDir)
}

// BuildStubData prepares naming variations for code templates.
func BuildStubData(moduleName, rawName string) StubData {
	pascal := strutil.ToPascalCase(rawName)
	camel := strutil.ToCamelCase(rawName)
	snake := strutil.ToSnakeCase(rawName)
	kebab := strutil.ToKebabCase(rawName)
	pluralSnake := strutil.Pluralize(snake)
	pluralPascal := strutil.Pluralize(pascal)

	return StubData{
		ModuleName:         moduleName,
		EntityPascal:       pascal,
		EntityCamel:        camel,
		EntitySnake:        snake,
		EntityKebab:        kebab,
		EntityPluralSnake:  pluralSnake,
		EntityPluralPascal: pluralPascal,
	}
}

// renderStub renders a stub from templates/stubs and writes it to targetPath.
func renderStub(stubName, targetPath string, data StubData) error {
	if _, err := os.Stat(targetPath); err == nil {
		ui.Warn("Skipped existing file: %s", targetPath)
		return nil
	}

	content, err := templates.FS.ReadFile("stubs/" + stubName)
	if err != nil {
		return fmt.Errorf("failed to read stub %s: %w", stubName, err)
	}

	tmpl, err := template.New(stubName).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse stub %s: %w", stubName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute stub %s: %w", stubName, err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(targetPath, buf.Bytes(), 0644); err != nil {
		return err
	}

	ui.Success("Created: %s", targetPath)
	return nil
}

// GenerateModule creates all layers for a Clean Architecture feature.
func GenerateModule(rootPath, rawName string) error {
	moduleName, projectRoot, err := FindGoModuleName(rootPath)
	if err != nil {
		return err
	}

	data := BuildStubData(moduleName, rawName)

	ui.Header(fmt.Sprintf("Generating Clean Architecture Module: %s", data.EntityPascal))

	files := []struct {
		stub string
		dest string
	}{
		{"entity.tmpl", filepath.Join(projectRoot, "internal", "domain", "entity", data.EntitySnake+".go")},
		{"repository.tmpl", filepath.Join(projectRoot, "internal", "domain", "repository", data.EntitySnake+"_repository.go")},
		{"persistence.tmpl", filepath.Join(projectRoot, "internal", "infrastructure", "persistence", data.EntitySnake+"_repository.go")},
		{"usecase_interface.tmpl", filepath.Join(projectRoot, "internal", "usecase", "interface", data.EntitySnake+"_usecase.go")},
		{"usecase_implementation.tmpl", filepath.Join(projectRoot, "internal", "usecase", "implementation", data.EntitySnake+"_usecase.go")},
		{"handler.tmpl", filepath.Join(projectRoot, "internal", "delivery", "http", "handlers", data.EntitySnake+"_handler.go")},
	}

	for _, f := range files {
		if err := renderStub(f.stub, f.dest, data); err != nil {
			return err
		}
	}

	// Generate migration
	timestamp := time.Now().Format("20060102150405")
	upMigration := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, data.EntityPluralSnake))
	downMigration := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, data.EntityPluralSnake))

	if err := renderStub("migration.tmpl", upMigration, data); err != nil {
		return err
	}

	downSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", data.EntityPluralSnake)
	_ = os.WriteFile(downMigration, []byte(downSQL), 0644)
	ui.Success("Created: %s", downMigration)

	fmt.Println()
	ui.Info("Next steps to register the module:")
	fmt.Printf("  1. In %s/internal/infrastructure/container/container.go:\n", projectRoot)
	fmt.Printf("     %sRepo := persistence.New%sRepository(db)\n", data.EntityCamel, data.EntityPascal)
	fmt.Printf("     %sUsecase := usecaseimpl.New%sUsecase(%sRepo)\n", data.EntityCamel, data.EntityPascal, data.EntityCamel)
	fmt.Printf("     %sHandler := handlers.New%sHandler(%sUsecase)\n\n", data.EntityCamel, data.EntityPascal, data.EntityCamel)
	fmt.Printf("  2. In %s/internal/delivery/http/routes.go:\n", projectRoot)
	fmt.Printf("     %s := v1.Group(\"/%s\")\n", data.EntityPluralSnake, data.EntityPluralSnake)
	fmt.Printf("     %s.POST(\"\", c.%sHandler.Create)\n", data.EntityPluralSnake, data.EntityPascal)
	fmt.Printf("     %s.GET(\"\", c.%sHandler.List)\n", data.EntityPluralSnake, data.EntityPascal)
	fmt.Printf("     %s.GET(\"/:id\", c.%sHandler.GetByID)\n", data.EntityPluralSnake, data.EntityPascal)
	fmt.Printf("     %s.DELETE(\"/:id\", c.%sHandler.Delete)\n\n", data.EntityPluralSnake, data.EntityPascal)

	return nil
}

// GenerateSingleLayer generates an isolated layer (entity, repo, usecase, handler, or migration).
func GenerateSingleLayer(rootPath, layerType, rawName string) error {
	moduleName, projectRoot, err := FindGoModuleName(rootPath)
	if err != nil {
		return err
	}

	data := BuildStubData(moduleName, rawName)

	switch strings.ToLower(layerType) {
	case "entity":
		dest := filepath.Join(projectRoot, "internal", "domain", "entity", data.EntitySnake+".go")
		return renderStub("entity.tmpl", dest, data)
	case "repo", "repository":
		dest1 := filepath.Join(projectRoot, "internal", "domain", "repository", data.EntitySnake+"_repository.go")
		dest2 := filepath.Join(projectRoot, "internal", "infrastructure", "persistence", data.EntitySnake+"_repository.go")
		if err := renderStub("repository.tmpl", dest1, data); err != nil {
			return err
		}
		return renderStub("persistence.tmpl", dest2, data)
	case "usecase":
		dest1 := filepath.Join(projectRoot, "internal", "usecase", "interface", data.EntitySnake+"_usecase.go")
		dest2 := filepath.Join(projectRoot, "internal", "usecase", "implementation", data.EntitySnake+"_usecase.go")
		if err := renderStub("usecase_interface.tmpl", dest1, data); err != nil {
			return err
		}
		return renderStub("usecase_implementation.tmpl", dest2, data)
	case "handler":
		dest := filepath.Join(projectRoot, "internal", "delivery", "http", "handlers", data.EntitySnake+"_handler.go")
		return renderStub("handler.tmpl", dest, data)
	case "migration":
		timestamp := time.Now().Format("20060102150405")
		upDest := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, data.EntityPluralSnake))
		downDest := filepath.Join(projectRoot, "migrations", fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, data.EntityPluralSnake))
		if err := renderStub("migration.tmpl", upDest, data); err != nil {
			return err
		}
		downSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", data.EntityPluralSnake)
		_ = os.WriteFile(downDest, []byte(downSQL), 0644)
		ui.Success("Created: %s", downDest)
		return nil
	default:
		return fmt.Errorf("unknown layer type: %s", layerType)
	}
}
