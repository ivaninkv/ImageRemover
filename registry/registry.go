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
	// Мапа для хранения образов
	images := make(map[string]bool)

	for registry := range cfg.DockerRegistry {
		hub := getRegistry(cfg.DockerRegistry[registry].ServerUrl,
			cfg.DockerRegistry[registry].User,
			cfg.DockerRegistry[registry].Password)
		repos, err := hub.Repositories()
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			if strings.HasPrefix(repo, cfg.DockerRegistry[registry].Folder) {
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
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.RegistryFileName, images)
	}
	return images
}

func DeleteImages(cfg config.Config, images map[string]bool) {
	for registry := range cfg.DockerRegistry {
		hub := getRegistry(cfg.DockerRegistry[registry].ServerUrl,
			cfg.DockerRegistry[registry].User,
			cfg.DockerRegistry[registry].Password)
		for img := range images {
			split := strings.Split(img, ":")
			digest, err := hub.ManifestDigest(split[0], split[1])
			if err != nil {
				logger.Log.Warn().Err(err).Str("image", img).Msg("Can't delete image")
			} else {
				if cfg.DockerRegistry[registry].DeleteImages {
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

}

func getRegistry(ServerUrl, User, Password string) *registry.Registry {
	// Создание клиента Docker Registry
	hub, err := registry.New(ServerUrl, User, Password)
	if err != nil {
		panic(err)
	}
	return hub
}
