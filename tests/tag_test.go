package tests

import (
	"testing"

	"git-gone/internal/git"
)

func TestGetLocalTags_ReturnsAllTags(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create tags
	runGitCmd(t, "tag", "v1.0.0")
	runGitCmd(t, "tag", "v1.1.0")
	runGitCmd(t, "tag", "-a", "v2.0.0", "-m", "Release 2.0.0")

	tags, err := git.GetLocalTags()
	if err != nil {
		t.Fatalf("GetLocalTags failed: %v", err)
	}

	// Should contain all three tags
	expectedTags := []string{"v1.0.0", "v1.1.0", "v2.0.0"}
	for _, exp := range expectedTags {
		found := false
		for _, tag := range tags {
			if tag == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected tag '%s' in result, got: %v", exp, tags)
		}
	}
}

func TestGetLocalTags_ReturnsEmptyForNoTags(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	tags, err := git.GetLocalTags()
	if err != nil {
		t.Fatalf("GetLocalTags failed: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("Expected no tags, got: %v", tags)
	}
}

func TestDeleteTag_DeletesLocalTag(t *testing.T) {
	h := NewTestHelper(t)
	defer h.Cleanup()

	// Create a tag
	runGitCmd(t, "tag", "v-to-delete")

	// Verify tag exists
	tagsBefore, _ := git.GetLocalTags()
	found := false
	for _, tag := range tagsBefore {
		if tag == "v-to-delete" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("Tag should exist before deletion")
	}

	// Delete the tag
	err := git.DeleteTag("v-to-delete")
	if err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}

	// Verify tag is gone
	tagsAfter, _ := git.GetLocalTags()
	for _, tag := range tagsAfter {
		if tag == "v-to-delete" {
			t.Error("Tag should not exist after deletion")
		}
	}
}
