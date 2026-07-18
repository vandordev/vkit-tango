package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncUsecaseWritesCommandProvider(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.test/project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "internal/usecase")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	source := `package usecase
import "example.test/project/internal/contract"
type Example struct{}
type ExampleInput struct{}
type ExampleResult struct{}
var _ contract.Command[ExampleInput, ExampleResult] = (*Example)(nil)
func NewExample() *Example { return &Example{} }
`
	if err := os.WriteFile(filepath.Join(dir, "example.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := sync(root, "usecase"); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(root, "internal/generated/fx/usecases_gen.go"))
	if err != nil || string(output) == "" {
		t.Fatalf("generated output = %q, %v", output, err)
	}
}
