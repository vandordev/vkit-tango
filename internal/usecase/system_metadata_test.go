package usecase_test

import (
	"testing"

	"github.com/vandordev/vkit-tango/internal/contract"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

func TestSetSystemMetadataInputIsIntentSpecific(t *testing.T) {
	input := usecase.SetSystemMetadataInput{Key: "maintenance", Value: map[string]any{"enabled": true}}
	if input.Key != "maintenance" || input.Value["enabled"] != true {
		t.Fatalf("unexpected input: %#v", input)
	}
}

func TestNewSetSystemMetadataImplementsCommand(t *testing.T) {
	command := usecase.NewSetSystemMetadata(usecase.Runner{})
	var _ contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult] = command
}
