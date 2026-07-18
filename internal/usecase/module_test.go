package usecase_test

import (
	"testing"

	"github.com/vandordev/vkit-tango/internal/usecase"
	"go.uber.org/fx"
)

func TestModuleProvidesRunner(t *testing.T) {
	var _ fx.Option = usecase.Module
}
