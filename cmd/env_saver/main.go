package main

import (
	"fmt"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	"github.com/Kuksenok-i-s/env_saver/pkg/storage"
	"github.com/Kuksenok-i-s/env_saver/pkg/utils"
)

func main() {
	config := config.Config{
		WatchDir:   "/home/user/",
		RepoDir:    "./config-repo",
		RemoteRepo: "private-repo-url",
		Branch:     "main",
		SecretKey:  "----",
	}

	if err := storage.SaveSecret(config.SecretKey, "your-secret-value"); err != nil {
		fmt.Printf("Error saving secret: %v\n", err)
	}

	secret, err := storage.GetSecret(config.SecretKey)
	if err != nil {
		fmt.Printf("Error retrieving secret: %v\n", err)
	} else {
		fmt.Printf("Retrieved secret: %s\n", secret)
	}

	utils.WatchFileChanges(&config)
}
