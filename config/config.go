package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	SecretFile          string            `yaml:"secretFile"`
	TokenFile           string            `yaml:"tokenFile"`
	SlackUrl            string            `yaml:"slackUrl"`
	LabelChannels       map[string]string `yaml:"labelChannels"`
	ArchiveSentMessages bool              `yaml:"archiveSentMessages"`
}

var loadedCfg *Config

func Load() *Config {
	if loadedCfg == nil {
		confFile := os.Getenv("CONFIG_FILE")
		if confFile == "" {
			confFile = "config.yaml"
		}
		contents, err := ioutil.ReadFile(confFile)
		if err != nil {
			log.Fatal(err)
		}
		cfg := &Config{}
		err = yaml.Unmarshal(contents, cfg)
		if err != nil {
			log.Fatal(err)
		}
		loadedCfg = cfg
	}
	return loadedCfg
}
