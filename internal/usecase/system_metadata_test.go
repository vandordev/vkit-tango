package usecase

import "testing"

func TestSetSystemMetadataInputIsIntentSpecific(t *testing.T) {
	input := SetSystemMetadataInput{Key: "maintenance", Value: map[string]any{"enabled": true}}
	if input.Key != "maintenance" || input.Value["enabled"] != true {
		t.Fatalf("unexpected input: %#v", input)
	}
}
