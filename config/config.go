package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Config struct {
	SecretFile          string            `yaml:"secretFile"`
	TokenFile           string            `yaml:"tokenFile"`
	SlackUrl            string            `yaml:"slackUrl"`
	LabelChannels       map[string]string `yaml:"labelChannels"`
	ArchiveSentMessages bool              `yaml:"archiveSentMessages"`
}

var loadedCfg *Config
var ConfFile = ""

func Load() *Config {
	if loadedCfg == nil {
		if ConfFile == "" {
			// Look for config.yaml in the current working directory then the executable's directory
			if _, err := os.Stat("config.yaml"); err == nil {
				ConfFile = "config.yaml"
			} else {
				ConfFile = path.Join(path.Dir(os.Args[0]), "config.yaml")
			}
			envFile := os.Getenv("CONFIG_FILE")
			if envFile != "" {
				ConfFile = envFile
			}
		}
		contents, err := ioutil.ReadFile(ConfFile)
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
