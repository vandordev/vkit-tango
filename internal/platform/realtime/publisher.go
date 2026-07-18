package realtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const ResourceUpdatedV1 = "resource.updated.v1"

type Event struct {
	Type        string `json:"type"`
	EventID     string `json:"event_id"`
	OccurredAt  string `json:"occurred_at"`
	ResourceID  string `json:"resource_id"`
	WorkspaceID string `json:"workspace_id"`
}

func (event Event) Validate() error {
	if event.Type != ResourceUpdatedV1 || event.EventID == "" || event.OccurredAt == "" || event.ResourceID == "" || event.WorkspaceID == "" {
		return fmt.Errorf("invalid realtime event")
	}
	return nil
}

type Publisher interface {
	Publish(context.Context, Event) error
}

type HTTPPublisher struct {
	BaseURL, APIKey string
	Client          *http.Client
}

func (publisher HTTPPublisher) Publish(ctx context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	baseURL := strings.TrimRight(publisher.BaseURL, "/")
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/internal/events", bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	request.Header.Set("x-realtime-api-key", publisher.APIKey)
	client := publisher.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("realtime publisher rejected event: %s", response.Status)
	}
	return nil
}
