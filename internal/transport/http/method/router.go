package method

import (
	"fmt"
	"regexp"
	"strings"
)

var versionSegment = regexp.MustCompile(`(^|/)v[0-9]+(/|$)`)

type Router struct {
	version string
}

func New(version string) Router {
	return Router{version: version}
}

func (router Router) Path(relativePath string) (string, error) {
	if !strings.HasPrefix(relativePath, "/") {
		return "", fmt.Errorf("route path %q must begin with /", relativePath)
	}
	if strings.Contains(relativePath, "/api") {
		return "", fmt.Errorf("route path %q must not contain /api", relativePath)
	}
	if versionSegment.MatchString(relativePath) {
		return "", fmt.Errorf("route path %q must not contain an API version", relativePath)
	}

	return "/api/" + router.version + relativePath, nil
}
