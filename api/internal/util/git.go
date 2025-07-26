package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneGitRepo(url string, branch string) (string, error) {

	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("repo-%d", time.Now().UnixNano()))

	fmt.Println("⏳ Cloning repo to:", tmpDir)

	_, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})

	if err != nil {
		return "", fmt.Errorf("❌ clone failed: %w", err)
	}

	return tmpDir, nil
}
