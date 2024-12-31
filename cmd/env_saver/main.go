package main

import (
	"fmt"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	"github.com/Kuksenok-i-s/env_saver/pkg/storage"
	"github.com/Kuksenok-i-s/env_saver/pkg/utils"
)

func main() {
	config := config.GetConfig()

	storage.InitStorage(config)
	fmt.Println("Watching for changes in directory:", config.WatchDir)
	utils.WatchFileChanges(config)

}
