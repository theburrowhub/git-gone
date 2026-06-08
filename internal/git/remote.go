package git

import (
	"fmt"
	"os"
	"os/exec"
)

// UpdateRemoteRefs fetches all remotes and prunes deleted references.
func UpdateRemoteRefs() error {
	cmd := exec.Command("git", "fetch", "--all", "--prune")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fetch failed: %s", string(output))
	}
	return nil
}

// UpdateRemoteRefsSync fetches all remotes and prunes, then runs an additional
// "git remote update --prune" pass to fully reconcile tracking refs.
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
