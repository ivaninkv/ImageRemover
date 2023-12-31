package registry

import (
	"encoding/json"
	"fmt"
	"imageRemover/config"
	"imageRemover/logger"
	"imageRemover/output"
	"strings"
)

func GetImages(cfg config.Config) map[string]bool {
	logger.Log.Debug().Msg("Getting images from registry")

	images := make(map[string]bool)

	for _, registryConfig := range cfg.DockerRegistry {
		headers := makeDockerHeaders(registryConfig.User, registryConfig.Password)

		repos := getRepos(registryConfig.ServerUrl, &headers)

		for _, repo := range repos {
			if strings.Contains(repo, registryConfig.Folder) {
				logger.Log.Info().Str("repo", repo).Msg("Handling repo:")
				tagsURL := fmt.Sprintf("%s/v2/%s/tags/list", registryConfig.ServerUrl, repo)
				tagsResp, err := doRequest("GET", tagsURL, &headers)
				if err != nil {
					logger.Log.Error().Err(err).Msg("Can't get tags from registry")
				}

				var tags struct {
					Tags []string `json:"tags"`
				}
				if err := json.NewDecoder(tagsResp.Body).Decode(&tags); err != nil {
					logger.Log.Error().Err(err).Msg("Can't decode tags from registry")
				}
				for _, tag := range tags.Tags {
					if tag != "latest" {
						logger.Log.Info().Str("tag", tag).Msg("Handling tag:")
						images[fmt.Sprintf("%s:%s", repo, tag)] = true
					}
				}

				if err := tagsResp.Body.Close(); err != nil {
					logger.Log.Error().Err(err).Msg("Can't close tags response")
				}
			}
		}
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.RegistryFileName, images)
	}
	return images
}

func DeleteImages(cfg config.Config, images map[string]bool) {
	for _, registryConfig := range cfg.DockerRegistry {
		dockerHeaders := makeDockerHeaders(registryConfig.User, registryConfig.Password)
		for image := range images {
			parts := strings.Split(image, ":")
			repo, tag := parts[0], parts[1]
			dig, err := getDigest(registryConfig.ServerUrl,
				registryConfig.User, registryConfig.Password, repo, tag)
			if err != nil {
				logger.Log.Error().Err(err).Msg("Can't get digest from registry")
				continue
			}
			if registryConfig.DeleteImages {
				deleteImage(registryConfig.ServerUrl, repo, dig, dockerHeaders, tag)
				if registryConfig.Nexus.Url != "" {
					err := deleteAsset(registryConfig.Nexus.Url, registryConfig.Nexus.RepositoryName,
						registryConfig.User, registryConfig.Password, dig.String())
					if err != nil {
						logger.Log.Error().Err(err).Msg("Can't delete asset from nexus")
					}
				}
			} else {
				logger.Log.Info().Str("repo", repo).Str("tag", tag).
					Str("digest", dig.String()).
					Msg("Skipped deleting manifest from registry")
			}
		}
	}
}
