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

func parsePathParams(path string, pattern string) PathMatch {
	u, err := url.Parse(path)
	if err != nil {
		return PathMatch{Match: false}
	}

	path = strings.Trim(u.Path, "/")
	query := u.Query()

	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")

	params := map[string]string{}
	i := 0

	for i < len(patternParts) {
		if i >= len(pathParts) {
			if patternParts[i] == "*" {
				params["wildcard"] = ""
				break
			}
			return PathMatch{Match: false}
		}

		pp := patternParts[i]
		p := pathParts[i]

		switch {
		case strings.HasPrefix(pp, ":"):
			paramName := strings.TrimPrefix(pp, ":")
			params[paramName] = p

		case pp == "*":
			params["wildcard"] = strings.Join(pathParts[i:], "/")
			return PathMatch{Params: params, Query: query, Match: true}

		case pp != p:
			return PathMatch{Match: false}
		}
		i++
	}

	if len(pathParts) != len(patternParts) && (i == len(patternParts)) {
		// còn dư path nhưng hết pattern
		return PathMatch{Match: false}
	}

	return PathMatch{
		Params: params,
		Query:  query,
		Match:  true,
	}
}
