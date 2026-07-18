package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestSyncRejectsConstructorWithoutContractAssertion(t *testing.T) {
	root := fixtureProject(t, "internal/usecase", `package usecase
	type Missing struct{}
type MissingInput struct{}
type MissingResult struct{}
func NewMissing() *Missing { return &Missing{} }
`)
	err := sync(root, "usecase")
	if err == nil || !strings.Contains(err.Error(), "missing contract.Command assertion") {
		t.Fatalf("sync() error = %v, want missing contract assertion", err)
	}
}

func TestSyncRejectsContractTypeWithoutConstructor(t *testing.T) {
	root := fixtureProject(t, "internal/usecase", `package usecase
import "example.test/project/internal/contract"
type Missing struct{}
type MissingInput struct{}
type MissingResult struct{}
var _ contract.Command[MissingInput, MissingResult] = (*Missing)(nil)
`)
	err := sync(root, "usecase")
	if err == nil || !strings.Contains(err.Error(), "missing NewMissing constructor") {
		t.Fatalf("sync() error = %v, want missing constructor", err)
	}
}

func TestSyncGeneratesHTTPRegistrationForDiscoveredHandler(t *testing.T) {
	root := fixtureProject(t, "internal/transport/http/handler/example", `package example
import "example.test/project/internal/contract"
type ExampleHandler struct{}
var _ contract.HTTPHandler = (*ExampleHandler)(nil)
func NewExampleHandler() *ExampleHandler { return &ExampleHandler{} }
`)
	if err := sync(root, "http"); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(root, "internal/generated/fx/http_gen.go"))
	if err != nil || !strings.Contains(string(output), "example.NewExampleHandler") || !strings.Contains(string(output), "handler.RegisterRoutes()") {
		t.Fatalf("generated HTTP registry = %s, %v", output, err)
	}
}

func fixtureProject(t *testing.T, directory, source string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.test/project\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, directory)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "example.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
