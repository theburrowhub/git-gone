# Feature Specification: Git Repository Cleanup Tool

**Feature Branch**: `001-repo-cleanup-tool`  
**Created**: 2025-12-16  
**Status**: Draft  
**Input**: User description: "Build a CLI tool as the ideal git plugin to help users keep their repository clean of branches, tags and other elements that are no longer needed, in a simple, friendly, verbose manner without risk of accidental destructive operations"

## Assumptions

The following reasonable defaults have been assumed:

- Tool operates on a single local git repository at a time
- Users have basic git knowledge (understand branches, tags, remotes)
- Confirmation prompts use simple y/N responses for standard operations
- "Verbose" means clear status messages with visual indicators (emojis, colors)
- Tool integrates with git as a plugin (`git gone`) and standalone (`git-gone`)
- Default behavior is always the safest option (require confirmation for any deletion)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Clean Merged Branches (Priority: P1)

As a developer, I want to remove local branches that have already been merged into the main branch so I can keep my branch list manageable and focused on active work.

**Why this priority**: This is the most common cleanup task developers perform daily. Merged branches are the safest to delete since their changes are already integrated.

**Independent Test**: Can be fully tested by creating a repository with merged branches and verifying only merged branches are offered for deletion, with proper confirmation flow.

**Acceptance Scenarios**:

1. **Given** a repository with 5 local branches (3 merged into main, 2 unmerged), **When** the user runs the cleanup command, **Then** only the 3 merged branches are shown as candidates for deletion
2. **Given** merged branches are displayed, **When** the user selects branches and confirms, **Then** selected branches are deleted and a summary is shown
3. **Given** the user is on a branch that would be deleted, **When** cleanup runs, **Then** the current branch is excluded from deletion candidates
4. **Given** the default branch (main/master), **When** cleanup runs, **Then** the default branch is never shown as a deletion candidate

---

### User Story 2 - Clean Branches with Deleted Remotes (Priority: P1)

As a developer working with pull requests, I want to remove local branches whose remote tracking branch has been deleted (squash/rebase merged) so I don't accumulate stale branches.

**Why this priority**: Equally important as merged branches since many teams use squash merges, leaving local branches orphaned even though changes are integrated.

**Independent Test**: Can be tested by creating branches with remote tracking, deleting the remote, and verifying the tool detects and offers these for cleanup.

**Acceptance Scenarios**:

1. **Given** a branch with remote tracking branch marked as "gone", **When** cleanup runs, **Then** this branch appears as a deletion candidate
2. **Given** branches with active remote tracking, **When** cleanup runs, **Then** these branches are NOT shown as deletion candidates
3. **Given** a mix of merged and gone-remote branches, **When** cleanup runs, **Then** both types are shown with clear indication of why each is deletable

---

### User Story 3 - Interactive Branch Selection (Priority: P1)

As a developer, I want to interactively select which branches to delete using a fuzzy finder so I can review and choose exactly what to remove.

**Why this priority**: User control is essential for safety. No deletions should happen without explicit user choice.

**Independent Test**: Can be tested by verifying the interactive selection UI allows navigation, filtering, multi-select, and cancellation.

**Acceptance Scenarios**:

1. **Given** multiple deletion candidates, **When** selection UI appears, **Then** user can navigate with arrow keys and toggle selection with Tab
2. **Given** the selection UI, **When** user types text, **Then** the list filters to show matching branch names
3. **Given** branches are selected, **When** user presses Enter, **Then** selection is confirmed and process continues
4. **Given** the selection UI, **When** user presses Escape, **Then** operation is cancelled with no changes made

---

### User Story 4 - Safe Deletion with Confirmation (Priority: P2)

As a developer, I want to see a clear summary of what will be deleted and confirm before any deletion occurs so I never accidentally lose work.

**Why this priority**: Safety is critical but builds on the selection mechanism from P1 stories.

**Independent Test**: Can be tested by verifying confirmation prompt appears, shows correct branches, and respects user response.

**Acceptance Scenarios**:

1. **Given** branches selected for deletion, **When** confirmation prompt appears, **Then** all selected branch names are listed clearly
2. **Given** confirmation prompt, **When** user responds "n" or presses Enter (default No), **Then** no branches are deleted
3. **Given** confirmation prompt, **When** user responds "y", **Then** branches are deleted and results shown
4. **Given** the --force flag is used, **When** branches are selected, **Then** deletion proceeds without confirmation prompt

---

### User Story 5 - Clean Stale Tags (Priority: P2)

As a developer, I want to remove local tags that no longer exist on the remote so my tag list reflects the current state of the project.

**Why this priority**: Tags are less frequently used than branches but still accumulate over time. Important for release management workflows.

**Independent Test**: Can be tested by creating local tags, removing them from remote, and verifying the tool detects and offers these for cleanup.

**Acceptance Scenarios**:

1. **Given** local tags that don't exist on the remote, **When** tag cleanup runs, **Then** these tags are shown as deletion candidates
2. **Given** tags that exist both locally and on remote, **When** tag cleanup runs, **Then** these tags are NOT shown as deletion candidates
3. **Given** stale tags selected for deletion, **When** user confirms, **Then** tags are deleted locally with a clear summary

---

### User Story 6 - Delete Unmerged Branches (Priority: P3)

As a developer, I want the option to include unmerged branches in the cleanup list so I can remove abandoned work, but with extra safety confirmations.

**Why this priority**: Useful for cleaning abandoned branches but inherently risky since changes may be lost permanently.

**Independent Test**: Can be tested by verifying unmerged branches only appear with explicit flag, are clearly marked as dangerous, and require explicit typed confirmation.

**Acceptance Scenarios**:

1. **Given** unmerged branches exist, **When** cleanup runs without special flag, **Then** unmerged branches are NOT shown
2. **Given** the --unmerged flag is used, **When** cleanup runs, **Then** unmerged branches appear marked with a warning indicator
3. **Given** unmerged branches are selected, **When** confirmation appears, **Then** user must type "DELETE" to proceed (not just y/n)
4. **Given** unmerged branches selected, **When** deletion proceeds, **Then** both local and remote branches are deleted

---

### User Story 7 - Verbose Progress Feedback (Priority: P3)

As a developer, I want clear, friendly feedback throughout the cleanup process so I always know what the tool is doing and what happened.

**Why this priority**: Enhances user experience but tool is functional without verbose output.

**Independent Test**: Can be tested by verifying each operation phase produces appropriate status messages.

**Acceptance Scenarios**:

1. **Given** cleanup starts, **When** remote references are updated, **Then** a progress message is shown (e.g., "🔄 Updating remote references...")
2. **Given** cleanup analysis completes, **When** results are ready, **Then** summary shows counts by category (merged, gone remotes, etc.)
3. **Given** branches are deleted, **When** each deletion completes, **Then** success/failure is shown per branch with clear indicator
4. **Given** cleanup completes, **When** all operations finish, **Then** final summary shows total deleted count

---

### Edge Cases

- What happens when user runs cleanup outside a git repository? → Clear error message: "Not in a git repository"
- What happens when no branches are eligible for deletion? → Friendly message: "No branches to delete"
- What happens when deletion of a branch fails? → Show error for that branch, continue with others, report in summary
- What happens when remote is unreachable? → Warning about stale data, continue with available information
- What happens when user's locale is non-English? → Tool uses LC_ALL=C for consistent git output parsing

## Requirements *(mandatory)*

### Functional Requirements

**Branch Cleanup**
- **FR-001**: System MUST identify branches merged into the default branch
- **FR-002**: System MUST identify branches whose remote tracking branch is deleted
- **FR-003**: System MUST exclude the current branch from deletion candidates
- **FR-004**: System MUST exclude the default branch (main/master) from deletion candidates
- **FR-005**: System MUST support an option to include unmerged branches with enhanced safety

**Tag Cleanup**
- **FR-006**: System MUST identify local tags not present on remote
- **FR-007**: System MUST provide a dedicated command/flag for tag cleanup

**Interactive Selection**
- **FR-008**: System MUST provide fuzzy-search filtering of candidates
- **FR-009**: System MUST support multi-selection of items
- **FR-010**: System MUST allow cancellation at any point before deletion

**Safety**
- **FR-011**: System MUST show confirmation prompt before any deletion by default
- **FR-012**: System MUST support a --force flag to skip confirmation for safe deletions
- **FR-013**: System MUST require explicit typed confirmation ("DELETE") for unmerged branch deletion
- **FR-014**: System MUST never delete unmerged branches without the explicit --unmerged flag

**User Experience**
- **FR-015**: System MUST display clear progress indicators during operations
- **FR-016**: System MUST show success/failure status for each deleted item
- **FR-017**: System MUST provide a summary at completion
- **FR-018**: System MUST work as both standalone command and git plugin

**Compatibility**
- **FR-019**: System MUST work regardless of user's system locale
- **FR-020**: System MUST work on repositories with any default branch name

### Key Entities

- **Branch**: A git branch with attributes: name, merge status, remote tracking status, is-current, is-default
- **Tag**: A git tag with attributes: name, exists-on-remote
- **Deletion Candidate**: An item eligible for deletion with: type (branch/tag), name, reason (merged/gone/stale), risk-level (safe/dangerous), display-label (formatted for TUI)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete a typical branch cleanup workflow in under 30 seconds
- **SC-002**: Zero accidental deletions occur when users follow the confirmation prompts
- **SC-003**: Tool correctly identifies 100% of merged and gone-remote branches
- **SC-004**: 95% of users successfully complete their first cleanup without consulting documentation
- **SC-005**: Tool starts and displays candidates in under 2 seconds for repositories with up to 100 branches
- **SC-006**: All destructive operations are reversible through git reflog for 30 days (standard git behavior)
- **SC-007**: Users report the tool feels "safe" to use (qualitative feedback)
