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
	// TODO add validation function
	err = yaml.Unmarshal(f, &data)

	if err != nil {
		log.Fatal(err)
	}

	return getConfigFromMap(data)
}

func getConfigFromMap(data map[string]interface{}) *Config {
	return &Config{
		WatchDir:         data["watch_dir"].(string),
		WatchedFileTypes: data["watched_file_types"].(string),
		RemoteRepo:       data["remote_repo"].(string),
		MakeRemoteBackup: data["make_remote_backup"].(bool),
		MakeTags:         data["make_tags"].(bool)}
}
