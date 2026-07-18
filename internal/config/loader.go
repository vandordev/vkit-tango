package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var moduleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var interpolationPattern = regexp.MustCompile(`\$\{([A-Z][A-Z0-9_]*)(:-([^${}]*))?\}`)

type Loader struct {
	Directory   string
	Environment map[string]string
}

func (loader Loader) Load(modules ...string) (map[string]any, error) {
	config := map[string]any{}
	seen := map[string]struct{}{}

	for _, module := range modules {
		if !moduleNamePattern.MatchString(module) {
			return nil, fmt.Errorf("invalid configuration module name %q", module)
		}
		if _, ok := seen[module]; ok {
			return nil, fmt.Errorf("configuration module %q was selected more than once", module)
		}
		seen[module] = struct{}{}

		contents, err := os.ReadFile(filepath.Join(loader.Directory, module+".yaml"))
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("configuration module %q was not found", module)
			}
			return nil, fmt.Errorf("read configuration module %q: %w", module, err)
		}

		parsed := map[string]any{}
		if err := yaml.Unmarshal(contents, &parsed); err != nil {
			return nil, fmt.Errorf("parse configuration module %q: %w", module, err)
		}

		resolved, err := interpolate(parsed, loader.Environment, module)
		if err != nil {
			return nil, err
		}
		merge(config, resolved.(map[string]any))
	}

	return config, nil
}

func interpolate(value any, environment map[string]string, module string) (any, error) {
	switch typed := value.(type) {
	case string:
		if !strings.Contains(typed, "${") {
			return typed, nil
		}
		firstStart := strings.Index(typed, "${")
		nextStart := strings.Index(typed[firstStart+2:], "${")
		firstEnd := strings.Index(typed[firstStart:], "}")
		if nextStart >= 0 && firstEnd >= 0 && firstStart+2+nextStart < firstStart+firstEnd {
			return nil, fmt.Errorf("invalid configuration interpolation in module %q", module)
		}
		matches := interpolationPattern.FindAllStringSubmatchIndex(typed, -1)
		if len(matches) == 0 || strings.Contains(interpolationPattern.ReplaceAllString(typed, ""), "${") {
			return nil, fmt.Errorf("invalid configuration interpolation in module %q", module)
		}

		return interpolationPattern.ReplaceAllStringFunc(typed, func(match string) string {
			parts := interpolationPattern.FindStringSubmatch(match)
			name := parts[1]
			if resolved := environment[name]; resolved != "" {
				return resolved
			}
			if parts[2] != "" {
				return parts[3]
			}
			return "\x00missing:" + name
		}), nil
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, nested := range typed {
			resolved, err := interpolate(nested, environment, module)
			if err != nil {
				return nil, err
			}
			if stringValue, ok := resolved.(string); ok && strings.HasPrefix(stringValue, "\x00missing:") {
				return nil, fmt.Errorf("missing configuration environment variable %q in module %q", strings.TrimPrefix(stringValue, "\x00missing:"), module)
			}
			result[key] = resolved
		}
		return result, nil
	case []any:
		result := make([]any, len(typed))
		for index, nested := range typed {
			resolved, err := interpolate(nested, environment, module)
			if err != nil {
				return nil, err
			}
			result[index] = resolved
		}
		return result, nil
	default:
		return value, nil
	}
}

func merge(destination, source map[string]any) {
	for key, sourceValue := range source {
		if destinationValue, ok := destination[key].(map[string]any); ok {
			if sourceObject, ok := sourceValue.(map[string]any); ok {
				merge(destinationValue, sourceObject)
				continue
			}
		}
		destination[key] = sourceValue
	}
}
