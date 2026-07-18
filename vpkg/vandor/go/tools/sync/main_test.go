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

func TestSyncRejectsTypeWithoutItsOwnContractAssertion(t *testing.T) {
	root := fixtureProject(t, "internal/transport/http/handler/example", `package example
import "example.test/project/internal/contract"
type FirstHandler struct{}
type SecondHandler struct{}
var _ contract.HTTPHandler = (*FirstHandler)(nil)
func NewFirstHandler() *FirstHandler { return &FirstHandler{} }
func NewSecondHandler() *SecondHandler { return &SecondHandler{} }
`)
	err := sync(root, "http")
	if err == nil || !strings.Contains(err.Error(), "SecondHandler") || !strings.Contains(err.Error(), "missing contract.HTTPHandler assertion") {
		t.Fatalf("sync() error = %v, want second handler contract error", err)
	}
}

func TestSyncUsesDistinctAliasesForSamePackageNames(t *testing.T) {
	root := fixtureProject(t, "internal/transport/http/handler/first", `package handler
import "example.test/project/internal/contract"
type FirstHandler struct{}
var _ contract.HTTPHandler = (*FirstHandler)(nil)
func NewFirstHandler() *FirstHandler { return &FirstHandler{} }
`)
	second := filepath.Join(root, "internal/transport/http/handler/second")
	if err := os.MkdirAll(second, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(second, "second.go"), []byte(`package handler
import "example.test/project/internal/contract"
type SecondHandler struct{}
var _ contract.HTTPHandler = (*SecondHandler)(nil)
func NewSecondHandler() *SecondHandler { return &SecondHandler{} }
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := sync(root, "http"); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(root, "internal/generated/fx/http_gen.go"))
	if err != nil || strings.Contains(string(output), "handler \"") || !strings.Contains(string(output), "handler_internal_transport_http_handler_first") || !strings.Contains(string(output), "handler_internal_transport_http_handler_second") {
		t.Fatalf("generated registry has colliding aliases: %s, %v", output, err)
	}
}

func TestSyncGeneratesWorkerAndSchedulerRegistrations(t *testing.T) {
	root := fixtureProject(t, "internal/worker/river", `package river
import "example.test/project/internal/contract"
type ExampleRegistrar struct{}
var _ contract.WorkerRegistrar = (*ExampleRegistrar)(nil)
func NewExampleRegistrar() *ExampleRegistrar { return &ExampleRegistrar{} }
`)
	scheduler := filepath.Join(root, "internal/scheduler/river")
	if err := os.MkdirAll(scheduler, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(scheduler, "example.go"), []byte(`package river
import "example.test/project/internal/contract"
type ExampleRegistrar struct{}
var _ contract.SchedulerRegistrar = (*ExampleRegistrar)(nil)
func NewExampleRegistrar() *ExampleRegistrar { return &ExampleRegistrar{} }
`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := sync(root, "worker"); err != nil {
		t.Fatal(err)
	}
	if err := sync(root, "scheduler"); err != nil {
		t.Fatal(err)
	}
	worker, err := os.ReadFile(filepath.Join(root, "internal/generated/fx/worker_gen.go"))
	if err != nil || !strings.Contains(string(worker), "registrar.RegisterWorkers(workers)") {
		t.Fatalf("worker registry = %s, %v", worker, err)
	}
	periodic, err := os.ReadFile(filepath.Join(root, "internal/generated/fx/scheduler_gen.go"))
	if err != nil || !strings.Contains(string(periodic), "registrar.RegisterPeriodicJobs()") || !strings.Contains(string(periodic), `group:"periodic_jobs,flatten"`) {
		t.Fatalf("scheduler registry = %s, %v", periodic, err)
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
