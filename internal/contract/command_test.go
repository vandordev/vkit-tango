package contract_test

import (
	"context"
	"testing"

	"github.com/vandordev/vkit-tango/internal/contract"
)

type testCommand struct{}

func (testCommand) Execute(context.Context, string) (int, error) { return 1, nil }

var _ contract.Command[string, int] = testCommand{}

func TestCommandExecutesTypedInput(t *testing.T) {
	got, err := testCommand{}.Execute(context.Background(), "input")
	if err != nil || got != 1 {
		t.Fatalf("Execute() = (%d, %v), want (1, nil)", got, err)
	}
}
