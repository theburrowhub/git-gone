package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Tag represents a local git tag.
type Tag struct {
	Name           string
	ExistsOnRemote bool
	IsAnnotated    bool
}

// IsStale returns true if the tag doesn't exist on the remote.
func (t *Tag) IsStale() bool {
	return !t.ExistsOnRemote
}

// GetLocalTags returns all local tag names.
func GetLocalTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "-l")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var tags []string
	for _, line := range lines {
		tag := strings.TrimSpace(line)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

// GetRemoteTags returns all tag names from the remote.
func GetRemoteTags() ([]string, error) {
	cmd := exec.Command("git", "ls-remote", "--tags", "origin")
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var tags []string
	seen := make(map[string]bool)
	for _, line := range lines {
		// Format: SHA refs/tags/tagname or SHA refs/tags/tagname^{}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			ref := parts[1]
			// Extract tag name from refs/tags/tagname
			if strings.HasPrefix(ref, "refs/tags/") {
				tagName := strings.TrimPrefix(ref, "refs/tags/")
				tagName = strings.TrimSuffix(tagName, "^{}") // Remove annotated tag suffix
				if !seen[tagName] {
					tags = append(tags, tagName)
					seen[tagName] = true
				}
			}
		}
	}
	return tags, nil
}

// GetStaleTags returns local tags that don't exist on the remote.
func GetStaleTags() ([]string, error) {
	localTags, err := GetLocalTags()
	if err != nil {
		return nil, fmt.Errorf("failed to get local tags: %w", err)
	}

	remoteTags, err := GetRemoteTags()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote tags: %w", err)
	}

	// Create a set of remote tags for quick lookup
	remoteTagSet := make(map[string]bool)
	for _, tag := range remoteTags {
		remoteTagSet[tag] = true
	}

	// Find local tags not on remote
	var staleTags []string
	for _, tag := range localTags {
		if !remoteTagSet[tag] {
			staleTags = append(staleTags, tag)
		}
	}

	return staleTags, nil
}

// DeleteTag deletes a local tag.
func DeleteTag(name string) error {
	cmd := exec.Command("git", "tag", "-d", name)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}
