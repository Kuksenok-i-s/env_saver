package utils

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	"github.com/Kuksenok-i-s/env_saver/pkg/storage"

	"github.com/fsnotify/fsnotify"
)

func RestoreFiles(repoDir string) error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore files: %v", err)
	}
	return nil
}

func WatchFileChanges(config *config.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(config.WatchDir)
	if err != nil {
		fmt.Printf("Error watching directory: %v\n", err)
		return
	}

	fmt.Printf("Watching changes in %s\n", config.WatchDir)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Printf("Detected change: %s\n", event.Name)
				commitMsg := fmt.Sprintf("Config changed on %s", time.Now().Format(time.RFC1123))
				if err := storage.GitCommit(config.RemoteRepo, commitMsg); err != nil {
					fmt.Printf("Error committing changes: %v\n", err)
				}
				if err := storage.GitPush(config.RemoteRepo, commitMsg); err != nil {
					fmt.Printf("Error pushing changes: %v\n", err)
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}
