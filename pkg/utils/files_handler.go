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
	RestoreFiles(repoDir string) error
	WatchFileChanges()
	GetFilesByType() ([]string, error)
}

type filesHandler struct {
	config *config.Config
}

func NewFilesHandler(config *config.Config) FilesHandler {
	return &filesHandler{
		config: config,
	}
}

func (fh *filesHandler) RestoreFiles(repoDir string) error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore files: %v", err)
	}
	return nil
}

func (fh *filesHandler) WatchFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(fh.config.WatchDir)
	if err != nil {
		fmt.Printf("Error watching directory: %v\n", err)
		return
	}

	fmt.Printf("Watching changes in %s\n", fh.config.WatchDir)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Printf("Detected change: %s\n", event.Name)
				commitMsg := fmt.Sprintf("Config changed on %s", time.Now().Format(time.RFC1123))
				if err := storage.GitCommit(fh.config.RemoteRepo, commitMsg); err != nil {
					fmt.Printf("Error committing changes: %v\n", err)
				}
				if err := storage.GitPush(fh.config.RemoteRepo, commitMsg); err != nil {
					fmt.Printf("Error pushing changes: %v\n", err)
				}
			}
		case err := <-watcher.Errors:
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

func (fh *filesHandler) getFileTypes() []string {
	if fh.config.WatchedFileTypes == "" {
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

	types := strings.Split(fh.config.WatchedFileTypes, ",")

	for i := range types {
		types[i] = strings.TrimSpace(types[i])
	}
	return types
}

// TODO: refactor
func (fh *filesHandler) GetFilesByType() ([]string, error) {
	files, err := os.ReadDir(fh.config.WatchDir)
	if err != nil {
		return nil, err
	}
	var matchingFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		for _, fileType := range fh.getFileTypes() {
			if strings.HasSuffix(file.Name(), fileType) {
				matchingFiles = append(matchingFiles, file.Name())
			}
		}
	}
	if len(matchingFiles) == 0 {
		return nil, fmt.Errorf("no files found with type %s", fh.config.WatchedFileTypes)
	}
	return matchingFiles, nil
}
