package service

import (
	"fmt"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	"github.com/Kuksenok-i-s/env_saver/pkg/storage"
	"github.com/Kuksenok-i-s/env_saver/pkg/utils"
)

// 1) Create repository
// 1.1) Add origitn if requested
// 2) Find configs in the dir
// 3) Make a copy of the configs
// 4) Add configs into repositoiry
// 5) Find secrets and save them into wallet
// 6) Make an initial commit
// 6.1) Push the commit if requested
// 7) Set trigger on the configs
// 8) Start watching for changes
// 8.1) If change detected, go to 3

func Start(config *config.Config) error {
	// TODO: refactor
	if config.RepositoryDir == "" {
		return fmt.Errorf("repository directory is not specified")
	}

	gitStorage := storage.NewGitStorage(config)
	localStorage := storage.NewLocalStorage()
	filesHandler := utils.NewFilesHandler(config)

	gitStorage.GitInit()
	localStorage.WriteConfig(config)
	filesHandler.SaveFiles(config.RepositoryDir)

	events := make(chan utils.FileUpdateEvent)
	errors := make(chan error)
	// TODO: Add Goroutine for GitTags
	go func() {
		filesHandler.WatchFileChanges(events, errors)
		for {
			select {
			case event := <-events:
				fmt.Println("Event:", event)
				filesHandler.SaveFiles(config.RepositoryDir)
				commitId, err := gitStorage.GitCommit(event.EventMessage)
				if err != nil {
					fmt.Println("Error:", err)
				}
				localStorage.SaveEvent(storage.Event{
					Title:       event.EventMessage,
					Description: event.EventMessage,
					Date:        event.Time.String(),
					CommitID:    commitId,
					CommitMsg:   event.EventMessage,
				})
				if config.MakeRemoteBackup {
					gitStorage.GitPush()
				}
			case err := <-errors:
				fmt.Println("Error:", err)
			}
		}
	}()
	return nil
}
