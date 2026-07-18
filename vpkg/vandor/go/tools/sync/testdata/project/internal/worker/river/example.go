package river

import "example.test/project/internal/contract"

type ExampleRegistrar struct{}

var _ contract.WorkerRegistrar = (*ExampleRegistrar)(nil)

func NewExampleRegistrar() *ExampleRegistrar { return &ExampleRegistrar{} }
