package web

import (
	"fmt"
	"net/http"
	"sketch-bridge/arduino-compile-server/common"
)

// ParseParameters parses the parameters in body. The parameters are as follows:
//   - projectId: The project ID.
func ParseParameters(r *http.Request) (*common.RequestParameters, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("failed to parse form: %v", err)
	}
	projectId := r.FormValue("projectId")
	if projectId == "" {
		return nil, fmt.Errorf("projectId is empty")
	}
	return &common.RequestParameters{
		ProjectId: projectId,
	}, nil
}
