package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func GetConfig() *Config {
	// lets save all in one file for now
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}

	err = yaml.Unmarshal(f, &data)

	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		WatchDir:   data["watch_dir"].(string),
		RepoDir:    data["repo_dir"].(string),
		RemoteRepo: data["remote_repo"].(string),
		Branch:     data["branch"].(string),
		SecretKey:  data["secret_key"].(string),
	}
}
