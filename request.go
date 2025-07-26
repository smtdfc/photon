package photon

import (
	"net/http"
	"strings"
)

type HttpAdapterRequest struct {
	Original *http.Request
}

func (r *HttpAdapterRequest) GetHeader(name string) string {
	return r.Original.Header.Get(name)
}

func (r *HttpAdapterRequest) GetAllHeader() (map[string]string, error) {
	headers := make(map[string]string)
	for headerName, headerValue := range r.Original.Header {
		headers[headerName] = strings.Join(headerValue, ", ")
	}

	return headers, nil
}
