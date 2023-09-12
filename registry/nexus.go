package registry

import (
	"encoding/json"
	"fmt"
	"strings"
)

func deleteAsset(nexusURL, repository, user, pass, digest string) error {
	split := strings.Split(digest, ":")
	getURL := fmt.Sprintf("%s/service/rest/v1/search/assets?repository=%s&%s=%s",
		nexusURL, repository, split[0], split[1])

	headers := makeHeaders(user, pass)
	getResp, err := doRequest("GET", getURL, &headers)
	defer getResp.Body.Close()
	if err != nil {
		return err
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(getResp.Body).Decode(&searchResponse); err != nil {
		return err
	}

	delURL := fmt.Sprintf("%s/service/rest/v1/assets/%s", nexusURL, searchResponse.Items[0].ID)
	delResp, err := doRequest("DELETE", delURL, &headers)
	defer delResp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

type Asset struct {
	ID          string `json:"id"`
	Repository  string `json:"repository"`
	Format      string `json:"format"`
	DownloadURL string `json:"downloadUrl"`
}

type SearchResponse struct {
	Items []Asset `json:"items"`
}
