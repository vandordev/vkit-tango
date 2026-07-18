package usecase

import "example.test/project/internal/contract"

type Example struct{}
type ExampleInput struct{}
type ExampleResult struct{}

var _ contract.Command[ExampleInput, ExampleResult] = (*Example)(nil)

func NewExample() *Example { return &Example{} }
