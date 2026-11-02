package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildStubData(t *testing.T) {
	data := BuildStubData("github.com/example/app", "product_category")

	if data.ModuleName != "github.com/example/app" {
		t.Errorf("expected module name github.com/example/app, got %s", data.ModuleName)
	}
	if data.EntityPascal != "ProductCategory" {
		t.Errorf("expected ProductCategory, got %s", data.EntityPascal)
	}
	if data.EntityCamel != "productCategory" {
		t.Errorf("expected productCategory, got %s", data.EntityCamel)
	}
	if data.EntitySnake != "product_category" {
		t.Errorf("expected product_category, got %s", data.EntitySnake)
	}
	if data.EntityPluralSnake != "product_categories" {
		t.Errorf("expected product_categories, got %s", data.EntityPluralSnake)
	}
}

func TestGenerateModuleWithTypes(t *testing.T) {
	tmpDir := t.TempDir()

	// Write a dummy go.mod
	goModContent := "module github.com/example/testapp\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := GenerateModule(tmpDir, "Product"); err != nil {
		t.Fatalf("GenerateModule failed: %v", err)
	}

	expectedFiles := []string{
		"internal/domain/entity/product.go",
		"internal/types/product.go",
		"internal/domain/repository/product_repository.go",
		"internal/infrastructure/persistence/product_repository.go",
		"internal/usecase/interface/product_usecase.go",
		"internal/usecase/implementation/product_usecase.go",
		"internal/delivery/http/handlers/product_handler.go",
	}

	for _, relPath := range expectedFiles {
		fullPath := filepath.Join(tmpDir, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", relPath)
		}
	}
}

func TestGenerateSingleLayerTypes(t *testing.T) {
	tmpDir := t.TempDir()

	goModContent := "module github.com/example/testapp\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	if err := GenerateSingleLayer(tmpDir, "types", "Order"); err != nil {
		t.Fatalf("GenerateSingleLayer types failed: %v", err)
	}

	orderTypesPath := filepath.Join(tmpDir, "internal", "types", "order.go")
	if _, err := os.Stat(orderTypesPath); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", orderTypesPath)
	}
}
