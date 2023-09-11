package registry

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/go-digest"
	"imageRemover/logger"
	"net/http"
)

func getRepos(serverUrl string, headers *map[string]string) []string {
	catalogURL := fmt.Sprintf("%s/v2/_catalog", serverUrl)
	catalogResp, err := doRequest("GET", catalogURL, headers)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Can't get catalog from registry")
		panic(err)
	}

	var catalog struct {
		Repositories []string `json:"repositories"`
	}
	if err := json.NewDecoder(catalogResp.Body).Decode(&catalog); err != nil {
		logger.Log.Error().Err(err).Msg("Can't decode catalog from registry")
		panic(err)
	}
	repos := catalog.Repositories

	if err := catalogResp.Body.Close(); err != nil {
		logger.Log.Error().Err(err).Msg("Can't close catalog response")
	}

	return repos
}

func deleteImage(serverUrl, repo string, dig digest.Digest, dockerHeaders map[string]string, tag string) {
	digestUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", serverUrl, repo, dig)
	digResp, err := doRequest("DELETE", digestUrl, &dockerHeaders)
	if err != nil {
		logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
			Str("digest", dig.String()).
			Msg("Can't delete manifest from registry")
	} else if digResp.StatusCode == http.StatusAccepted {
		logger.Log.Info().Str("repo", repo).Str("tag", tag).
			Str("digest", dig.String()).
			Msg("Deleted manifest from registry")
	} else {
		logger.Log.Warn().Str("repo", repo).Str("tag", tag).
			Str("digest", dig.String()).Str("status", digResp.Status).
			Int("statusCode", digResp.StatusCode).
			Msg("Can't delete manifest from registry. Registry returned: ")
	}

	if err := digResp.Body.Close(); err != nil {
		logger.Log.Error().Err(err).Msg("Can't close digest response")
	}

}
