package storage

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
)

// Here and later I'm going to be using the term "storage" to refer to the git repository.
// Understand and forgive me.

type GitStorageInterface interface {
	GitInit() error
	GitCommit(commitMsg string) error
	GitPush(commitMsg string) error
	RestoreFiles() error
}

type gitStorage struct {
	config  *config.Config
	repodir string
}

func NewGitStorage(config *config.Config) GitStorageInterface {
	repodir := strings.Split(config.RemoteRepo, "/")[len(strings.Split(config.RemoteRepo, "/"))-1]
	// "https://github.com/Kuksenok-i-s/env_saver" -> "env_saver"
	os.Mkdir(repodir, 0755)
	return &gitStorage{
		config:  config,
		repodir: repodir,
	}
}

func (gs *gitStorage) GitInit() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = gs.repodir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}
	if gs.config.RemoteRepo != "" {
		cmd = exec.Command("git", "remote", "add", "origin", gs.config.RemoteRepo)
		cmd.Dir = gs.repodir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add remote repository: %v", err)
		}
	}
	return nil
}

func (gs *gitStorage) GitCommit(commitMsg string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = gs.repodir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = gs.repodir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}
	return nil
}

func (gs *gitStorage) GitPush(commitMsg string) error {
	cmd := exec.Command("git", "push", "origin", "main")
	cmd.Dir = gs.repodir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes: %v", err)
	}

	return nil
}

func (gs *gitStorage) RestoreFiles() error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = gs.repodir
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to pull changes: %v", err)
		return gs.cloneRemoteRepo()
	}
	return nil
}

func (gs *gitStorage) cloneRemoteRepo() error {
	cmd := exec.Command("git", "clone", "-o", gs.repodir, gs.config.RemoteRepo)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}
