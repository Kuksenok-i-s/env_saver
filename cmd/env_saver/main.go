package main

import (
	"github.com/Kuksenok-i-s/env_saver/internal/service"
	"github.com/Kuksenok-i-s/env_saver/pkg/config"
)

func main() {
	config := config.GetConfig()
	service.Start(config)

}
