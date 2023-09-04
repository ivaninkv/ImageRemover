package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"imageRemover/logger"
	"os"
	"reflect"
)

type Config struct {
	KubeCluster []struct {
		ServerUrl string `yaml:"ServerUrl"`
		Namespace string `yaml:"Namespace"`
		Token     string `yaml:"Token"`
	} `yaml:"KubeCluster"`
	DockerRegistry []struct {
		ServerUrl    string `yaml:"ServerUrl"`
		Folder       string `yaml:"Folder"`
		User         string `yaml:"User"`
		Password     string `yaml:"Password"`
		DeleteImages bool   `yaml:"DeleteImages"`
	} `yaml:"DockerRegistry"`
	Output struct {
		WriteToTXT       bool   `yaml:"WriteToTXT"`
		RegistryFileName string `yaml:"RegistryFileName"`
		KubeFileName     string `yaml:"KubeFileName"`
		DiffFileName     string `yaml:"DiffFileName"`
	} `yaml:"Output"`
}

func ReadConfig(filePath string) (Config, error) {
	logger.Log.Debug().Msg("Reading config")

	config := Config{}
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	if reflect.DeepEqual(config, Config{}) {
		return config, errors.New("config is empty")
	}

	return config, nil
}
