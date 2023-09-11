package registry

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func makeHeaders(user, password string) map[string]string {
	auth := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", user, password)))
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", auth),
	}
	return headers
}

func makeDockerHeaders(user, password string) map[string]string {
	headers := makeHeaders(user, password)
	headers["Accept"] = "application/vnd.docker.distribution.manifest.v2+json"

	return headers
}

func doRequest(method string, url string, headers *map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range *headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}
