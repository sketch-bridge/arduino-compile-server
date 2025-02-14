package web

import (
	"fmt"
	"net/http"
	"sketch-bridge/arduino-compile-server/common"
)

// ParseQueryParameters parses the query parameters. The query parameters are as follows:
//   - projectId: The project ID.
func ParseQueryParameters(r *http.Request) (*common.RequestParameters, error) {
	queryParams := r.URL.Query()
	projectId := queryParams.Get("projectId")
	if projectId == "" {
		return nil, fmt.Errorf("projectId is empty")
	}
	return &common.RequestParameters{
		ProjectId: projectId,
	}, nil
}
