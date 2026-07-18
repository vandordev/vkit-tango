package main

import (
	"net/http"
	"testing"
)

func TestNewServerUsesConfiguredAddress(t *testing.T) {
	server := newServer("127.0.0.1:4101", http.NewServeMux())
	if server.Addr != "127.0.0.1:4101" {
		t.Fatalf("Addr = %q, want 127.0.0.1:4101", server.Addr)
	}
}
