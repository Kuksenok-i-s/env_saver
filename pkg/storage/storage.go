package storage

import (
	"fmt"
	"os/exec"
)

// Here and later I'm going to be using the term "storage" to refer to the git repository.
// Understand and forgive me.

func GitCommit(repoDir string, commitMsg string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}
	return nil
}

func GitPush(repoDir string, commitMsg string) error {
	cmd := exec.Command("git", "push", "origin", "main")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes: %v", err)
	}

	return nil
}

func RestoreFiles(repoDir string) error {
	cmd := exec.Command("git", "pull", "origin", "main")
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore files: %v", err)
	}
	return nil
}
