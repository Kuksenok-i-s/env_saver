package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"

	"github.com/fsnotify/fsnotify"
)

type FilesHandler interface {
	RestoreFiles(repoDir string) error
	WatchFileChanges(events chan<- string, errors chan<- error)
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
	files, err := fh.getFilesByType()
	if err != nil {
		return err
	}
	for _, file := range files {
		err := copyFile(file, repoDir)
		if err != nil {
			return fmt.Errorf("error copying file %s: %v", file, err)
		}
	}
	return nil
}

func copyFile(src, destDir string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	destPath := filepath.Join(destDir, filepath.Base(src))
	err = os.WriteFile(destPath, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (fh *filesHandler) WatchFileChanges(events chan<- string, errors chan<- error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errors <- fmt.Errorf("error creating watcher: %v", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(fh.config.WatchDir)
	if err != nil {
		errors <- fmt.Errorf("error watching directory: %v", err)
		return
	}

	fmt.Printf("Watching changes in %s\n", fh.config.WatchDir)
	// TODO: refactor
	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				for _, fileType := range fh.getFileTypes() {
					if strings.HasSuffix(event.Name, fileType) {
						log.Printf("Detected change: %s\n", event.Name)
						eventMsg := fmt.Sprintf("File changed: %s at %s", event.Name, time.Now().Format(time.RFC1123))
						events <- eventMsg
						break
					}
				}
			}
		case err := <-watcher.Errors:
			errors <- fmt.Errorf("watcher error: %v", err)
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
func (fh *filesHandler) getFilesByType() ([]string, error) {
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
