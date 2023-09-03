package main

import (
	"imageRemover/config"
	"imageRemover/kube"
	"imageRemover/logger"
	"imageRemover/output"
	"imageRemover/registry"
	"io"
	"log"
)

func main() {
	log.SetOutput(io.Discard) // Отключение стандартного логгера

	logger.Log.Debug().Msg("Starting application")

	configFile := "config.yaml"
	cfg, err := config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}

	kubeImages := kube.GetImages(cfg)
	registryImages := registry.GetImages(cfg)

	diffImages := diffMap(cfg, kubeImages, registryImages)
	registry.DeleteImages(cfg, diffImages)
}

func diffMap(cfg config.Config, kubeImages, registryImages map[string]bool) map[string]bool {
	diffImages := make(map[string]bool)
	for image := range registryImages {
		if _, ok := kubeImages[image]; !ok {
			diffImages[image] = true
		}
	}

	if cfg.Output.WriteToTXT {
		output.WriteToTXT(cfg.Output.DiffFileName, diffImages)
	}
	return diffImages
}
