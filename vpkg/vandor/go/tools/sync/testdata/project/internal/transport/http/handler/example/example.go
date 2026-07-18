package example

import "example.test/project/internal/contract"

type ExampleHandler struct{}

var _ contract.HTTPHandler = (*ExampleHandler)(nil)

func NewExampleHandler() *ExampleHandler { return &ExampleHandler{} }
