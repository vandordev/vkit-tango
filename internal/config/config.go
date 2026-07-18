package config

import (
	"fmt"
	"strconv"
)

type Database struct {
	URL string
}

type HTTPAPI struct {
	Host          string
	Port          int
	PublicBaseURL string
}

type Realtime struct {
	Host           string
	Port           int
	PublicURL      string
	TicketSecret   string
	InternalAPIKey string
}

type API struct {
	Database Database
	HTTPAPI  HTTPAPI
	Realtime Realtime
}

func LoadAPI(loader Loader) (API, error) {
	loaded, err := loader.Load("app", "database", "http_api", "realtime", "observability")
	if err != nil {
		return API{}, err
	}

	database, err := object(loaded, "database")
	if err != nil {
		return API{}, err
	}
	httpAPI, err := object(loaded, "http_api")
	if err != nil {
		return API{}, err
	}
	realtime, err := object(loaded, "realtime")
	if err != nil {
		return API{}, err
	}

	port, err := integer(httpAPI, "port")
	if err != nil {
		return API{}, err
	}
	realtimePort, err := integer(realtime, "port")
	if err != nil {
		return API{}, err
	}

	return API{
		Database: Database{URL: stringValue(database, "url")},
		HTTPAPI: HTTPAPI{
			Host:          stringValue(httpAPI, "host"),
			Port:          port,
			PublicBaseURL: stringValue(httpAPI, "public_base_url"),
		},
		Realtime: Realtime{
			Host:           stringValue(realtime, "host"),
			Port:           realtimePort,
			PublicURL:      stringValue(realtime, "public_url"),
			TicketSecret:   stringValue(realtime, "ticket_secret"),
			InternalAPIKey: stringValue(realtime, "internal_api_key"),
		},
	}, nil
}

func object(values map[string]any, key string) (map[string]any, error) {
	value, ok := values[key].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("configuration %q must be an object", key)
	}
	return value, nil
}

func stringValue(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}

func integer(values map[string]any, key string) (int, error) {
	value := stringValue(values, key)
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("configuration %q must be a positive integer", key)
	}
	return parsed, nil
}
