package method

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

func Tags(path string) []string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "v1" {
		panic(fmt.Sprintf("method.Tags requires an /api/v1 resource path: %q", path))
	}
	for _, segment := range segments[2:] {
		if segment != "" && !(strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")) {
			return []string{segment}
		}
	}
	panic(fmt.Sprintf("method.Tags requires an /api/v1 resource path: %q", path))
}

func GET[I, O any](api huma.API, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	register(api, http.MethodGet, path, summary, tags, handler)
}

func POST[I, O any](api huma.API, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	register(api, http.MethodPost, path, summary, tags, handler)
}

func PUT[I, O any](api huma.API, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	register(api, http.MethodPut, path, summary, tags, handler)
}

func PATCH[I, O any](api huma.API, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	register(api, http.MethodPatch, path, summary, tags, handler)
}

func DELETE[I, O any](api huma.API, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	register(api, http.MethodDelete, path, summary, tags, handler)
}

func register[I, O any](api huma.API, httpMethod, path, summary string, tags []string, handler func(context.Context, *I) (*O, error)) {
	operationID := strings.ToLower(httpMethod) + "_" + strings.NewReplacer("/", "_", "{", "", "}", "").Replace(strings.TrimPrefix(path, "/"))
	huma.Register(api, huma.Operation{
		OperationID: operationID,
		Method:      httpMethod,
		Path:        path,
		Summary:     summary,
		Tags:        tags,
	}, handler)
}
