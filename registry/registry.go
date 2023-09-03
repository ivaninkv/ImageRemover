package registry

import (
	"github.com/heroku/docker-registry-client/registry"
	"imageRemover/config"
	"imageRemover/logger"
	"imageRemover/output"
	"strings"
)

func GetImages(cfg config.Config) map[string]bool {
	logger.Log.Debug().Msg("Getting images from registry")
	hub := getRegistry(cfg)

	// Мапа для хранения образов
	images := make(map[string]bool)

	repos, err := hub.Repositories()
	if err != nil {
		panic(err)
	}
	for _, repo := range repos {
		if strings.HasPrefix(repo, cfg.DockerRegistry.Folder) {
			logger.Log.Debug().Str("repo", repo).Msg("Repositories:")
			tags, err := hub.Tags(repo)
			if err != nil {
				panic(err)
			}
			for _, tag := range tags {
				images[repo+":"+tag] = true
			}
		}
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.RegistryFileName, images)
	}
	return images
}

func DeleteImages(cfg config.Config, images map[string]bool) {
	hub := getRegistry(cfg)
	for img := range images {
		split := strings.Split(img, ":")
		digest, err := hub.ManifestDigest(split[0], split[1])
		if err != nil {
			logger.Log.Warn().Err(err).Str("image", img).Msg("Can't delete image")
		} else {
			if cfg.DockerRegistry.DeleteImages {
				err := hub.DeleteManifest(split[0], digest)
				if err != nil {
					logger.Log.Warn().Err(err).
						Str("image", img).
						Str("digest", digest.String()).
						Msg("Can't delete image")
				}
				logger.Log.Info().Str("image", img).Msg("Image deleted")
			} else {
				logger.Log.Info().Str("image", img).Msg("Image will be delete")
			}
		}
	}
}

func getRegistry(cfg config.Config) *registry.Registry {
	// Создание клиента Docker Registry
	hub, err := registry.New(cfg.DockerRegistry.ServerUrl,
		cfg.DockerRegistry.User,
		cfg.DockerRegistry.Password)
	if err != nil {
		panic(err)
	}
	return hub
}
