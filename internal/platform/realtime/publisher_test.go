package realtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func event() Event {
	return Event{Type: ResourceUpdatedV1, EventID: "b7fa9ad5-9c93-4cce-a83d-8d0438abef12", OccurredAt: "2026-07-18T00:00:00Z", ResourceID: "resource-1", WorkspaceID: "workspace-1"}
}

func TestHTTPPublisherPublishesVersionedEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/internal/events" || r.Header.Get("x-realtime-api-key") != "secret" || r.Header.Get("content-type") != "application/json" {
			t.Fatalf("unexpected request: %s", r.URL)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	if err := (HTTPPublisher{BaseURL: server.URL, APIKey: "secret"}).Publish(context.Background(), event()); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
}

func TestHTTPPublisherRejectsNonSuccessResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
	defer server.Close()
	if err := (HTTPPublisher{BaseURL: server.URL}).Publish(context.Background(), event()); err == nil {
		t.Fatal("Publish() error = nil")
	}
}
