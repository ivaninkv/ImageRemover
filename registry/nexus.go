package registry

import (
	"encoding/json"
	"fmt"
	"strings"
)

func deleteAsset(nexusURL, user, pass, digest string) error {
	const repository = "docker-registry"

	split := strings.Split(digest, ":")
	getURL := fmt.Sprintf("%s/service/rest/v1/search/assets?repository=%s&%s=%s",
		nexusURL, repository, split[0], split[1])

	headers := makeHeaders(user, pass)
	getResp, err := doRequest("GET", getURL, &headers)
	defer getResp.Body.Close()
	if err != nil {
		return err
	}

	var assets struct {
		ID          string `json:"id"`
		DownloadURL string `json:"downloadUrl"`
	}
	if err := json.NewDecoder(getResp.Body).Decode(&assets); err != nil {
		return err
	}

	delURL := fmt.Sprintf("%s/service/rest/v1/assets/%s", nexusURL, assets.ID)
	delResp, err := doRequest("DELETE", delURL, &headers)
	defer delResp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
