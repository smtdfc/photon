package photon

import (
	"net/url"
	"strings"
)

type PathMatch struct {
	Params map[string]string
	Query  url.Values
	Match  bool
}

func ParsePathParams(path string, pattern string) PathMatch {
	var rawPath string
	var query url.Values
	if idx := strings.Index(path, "?"); idx != -1 {
		rawPath = path[:idx]
		queryPart := path[idx+1:]
		query, _ = url.ParseQuery(queryPart)
	} else {
		rawPath = path
		query = url.Values{}
	}

	pathParts := strings.Split(strings.Trim(rawPath, "/"), "/")
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")

	if len(pathParts) != len(patternParts) {
		return PathMatch{Match: false}
	}

	params := map[string]string{}

	for i := 0; i < len(pathParts); i++ {
		pp := patternParts[i]
		p := pathParts[i]

		if strings.HasPrefix(pp, ":") {
			paramName := strings.TrimPrefix(pp, ":")
			params[paramName] = p
		} else if pp != p {
			return PathMatch{Match: false}
		}
	}

	return PathMatch{
		Params: params,
		Query:  query,
		Match:  true,
	}
}
