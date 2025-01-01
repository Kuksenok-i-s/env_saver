package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	"github.com/Kuksenok-i-s/env_saver/pkg/storage"

	"github.com/fsnotify/fsnotify"
)

type FilesHandler interface {
	HandleFileChanges(config *config.Config)
	RestoreFiles(repoDir string) error
	WatchFileChanges(config *config.Config)
	GetFilesByType(config *config.Config) ([]string, error)
	GenEnvsByFile(fileNames []string) (map[string]map[string]string, error)
}

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

func getFileTypes(config *config.Config) []string {
	if config.WatchedFileTypes == "" {
		defaultTypes := []string{
			".env",
			".env.development",
			".env.production",
			".env.test",
			".bashrc",
			".bash_profile",
			".zshrc",
			".profile",
			".vimrc",
			".gitconfig",
			"Makefile",
			// TODO add add handlers and parsers for other file types
			// ".yaml",
			// ".yml",
			// ".json",
			// ".toml",
			// ".ini",
			// ".npmrc",
			// ".dockerignore",
		}
		return defaultTypes
	}

	types := strings.Split(config.WatchedFileTypes, ",")

	for i := range types {
		types[i] = strings.TrimSpace(types[i])
	}
	return types
}

func GetFilesByType(config *config.Config) ([]string, error) {
	files, err := os.ReadDir(config.WatchDir)
	if err != nil {
		return nil, err
	}
	var matchingFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		for _, fileType := range getFileTypes(config) {
			if strings.HasSuffix(file.Name(), fileType) {
				matchingFiles = append(matchingFiles, file.Name())
			}
		}
	}
	if len(matchingFiles) == 0 {
		return nil, fmt.Errorf("no files found with type %s", config.WatchedFileTypes)
	}
	return matchingFiles, nil
}
