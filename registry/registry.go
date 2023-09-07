package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/opencontainers/go-digest"
	"imageRemover/config"
	"imageRemover/logger"
	"imageRemover/output"
	"net/http"
	"strings"
)

func GetImages(cfg config.Config) (images map[string]bool) {
	logger.Log.Debug().Msg("Getting images from registry")

	for _, registryConfig := range cfg.DockerRegistry {
		headers := makeHeaders(registryConfig.User, registryConfig.Password)

		repos := getRepos(registryConfig.ServerUrl, headers)

		for _, repo := range repos {
			if strings.Contains(repo, registryConfig.Folder) {
				tagsURL := fmt.Sprintf("%s/v2/%s/tags/list", registryConfig.ServerUrl, repo)
				tagsResp, err := doRequest("GET", tagsURL, headers)
				if err != nil {
					logger.Log.Error().Err(err).Msg("Can't get tags from registry")
				}
				defer tagsResp.Body.Close()

				var tags struct {
					Tags []string `json:"tags"`
				}
				if err := json.NewDecoder(tagsResp.Body).Decode(&tags); err != nil {
					logger.Log.Error().Err(err).Msg("Can't decode tags from registry")
				}
				for _, tag := range tags.Tags {
					images[fmt.Sprintf("%s:%s", repo, tag)] = true
				}
			}
		}
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.RegistryFileName, images)
	}
	return images
}

func getRepos(serverUrl string, headers map[string]string) []string {
	catalogURL := fmt.Sprintf("%s/v2/_catalog", serverUrl)
	catalogResp, err := doRequest("GET", catalogURL, headers)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Can't get catalog from registry")
		panic(err)
	}
	defer catalogResp.Body.Close()

	var catalog struct {
		Repositories []string `json:"repositories"`
	}
	if err := json.NewDecoder(catalogResp.Body).Decode(&catalog); err != nil {
		logger.Log.Error().Err(err).Msg("Can't decode catalog from registry")
		panic(err)
	}
	repos := catalog.Repositories
	return repos
}

func DeleteImages(cfg config.Config, images map[string]bool) {
	for _, registryConfig := range cfg.DockerRegistry {
		headers := makeHeaders(registryConfig.User, registryConfig.Password)
		for image := range images {
			parts := strings.Split(image, ":")
			repo, tag := parts[0], parts[1]
			tagUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", registryConfig.ServerUrl, repo, tag)
			tagResp, err := doRequest("GET", tagUrl, headers)
			if err != nil {
				logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
					Msg("Can't get manifest from registry")
			}
			defer tagResp.Body.Close()

			dig, err := digest.Parse(tagResp.Header.Get("Docker-Content-Digest"))
			if err != nil {
				logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
					Msg("Can't parse digest from registry")
			}
			digestUrl := fmt.Sprintf("%s/v2/%s/manifests/%s", registryConfig.ServerUrl, repo, dig)
			if registryConfig.DeleteImages {
				digResp, err := doRequest("DELETE", digestUrl, headers)
				if err != nil {
					logger.Log.Warn().Err(err).Str("repo", repo).Str("tag", tag).
						Str("digest", dig.String()).
						Msg("Can't delete manifest from registry")
				}
				defer digResp.Body.Close()
				if digResp.StatusCode == http.StatusAccepted {
					logger.Log.Info().Str("repo", repo).Str("tag", tag).
						Str("digest", dig.String()).
						Msg("Deleted manifest from registry")
				} else {
					logger.Log.Warn().Str("repo", repo).Str("tag", tag).
						Str("digest", dig.String()).Str("status", digResp.Status).
						Int("statusCode", digResp.StatusCode).
						Msg("Can't delete manifest from registry. Registry returned: ")
				}

			} else {
				logger.Log.Info().Str("repo", repo).Str("tag", tag).
					Str("digest", dig.String()).
					Msg("Skipped deleting manifest from registry")
			}
		}
	}
}

func makeHeaders(user, password string) map[string]string {
	auth := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", user, password)))
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", auth),
		"Accept":        "application/vnd.docker.distribution.manifest.v2+json",
	}
	return headers
}

func doRequest(method string, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}
