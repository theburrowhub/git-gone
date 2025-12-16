package git

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// UpdateRemoteRefs fetches all remotes and prunes deleted references.
// This operation runs asynchronously using goroutines.
func UpdateRemoteRefs() error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Fetch all remotes with prune
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd := exec.Command("git", "fetch", "--all", "--prune")
		cmd.Env = append(os.Environ(), "LC_ALL=C")
		output, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- fmt.Errorf("fetch failed: %s", string(output))
		}
	}()

	// Wait for fetch to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateRemoteRefsSync fetches all remotes and prunes synchronously.
func UpdateRemoteRefsSync() error {
	cmd := exec.Command("git", "fetch", "--all", "--prune")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fetch failed: %s", string(output))
	}

	cmd = exec.Command("git", "remote", "update", "--prune")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remote update failed: %s", string(output))
	}

	return nil
}

// HasRemote checks if the repository has an origin remote.
func HasRemote() bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	return cmd.Run() == nil
}
