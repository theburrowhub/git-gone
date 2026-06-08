# Tasks: Git Repository Cleanup Tool

**Input**: Design documents from `/specs/001-repo-cleanup-tool/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Integration tests REQUIRED per constitution (Principle II: Testing Standards).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1-US7)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Project structure and dependencies

- [X] T001 Create internal/git/ directory structure
- [X] T002 Create internal/tui/ directory structure
- [X] T003 Create internal/updater/ directory structure
- [X] T004 [P] Add Bubbletea dependencies to go.mod (bubbletea, bubbles, lipgloss)
- [X] T005 [P] Add semver dependency to go.mod (github.com/blang/semver)

**Checkpoint**: Project structure ready for implementation ✅

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T006 Create Repository type and git check in internal/git/repository.go
- [X] T007 [P] Create Branch type with all fields in internal/git/branch.go
- [X] T008 [P] Create Tag type with all fields in internal/git/tag.go
- [X] T009 [P] Create DeletionCandidate type in internal/git/candidate.go
- [X] T010 [P] Create TUI styles with lipgloss in internal/tui/styles.go
- [X] T011 Implement GetDefaultBranch() in internal/git/repository.go
- [X] T012 Implement GetCurrentBranch() in internal/git/repository.go
- [X] T013 Implement CheckGitRepository() in internal/git/repository.go

**Checkpoint**: Foundation ready - user story implementation can now begin ✅

---

## Phase 3: User Story 1 - Clean Merged Branches (Priority: P1) 🎯 MVP

**Goal**: Remove local branches that have been merged into the default branch

**Independent Test**: Create repo with merged/unmerged branches, verify only merged shown

### Tests for User Story 1

- [X] T014 [P] [US1] Create test file tests/branch_test.go with test helpers
- [X] T015 [P] [US1] Test GetMergedBranches returns only merged branches in tests/branch_test.go
- [X] T016 [P] [US1] Test FilterProtectedBranches excludes current/default in tests/branch_test.go

### Implementation for User Story 1

- [X] T017 [US1] Implement GetMergedBranches(defaultBranch) in internal/git/branch.go
- [X] T018 [US1] Implement FilterProtectedBranches() to exclude current/default in internal/git/branch.go
- [X] T019 [US1] Implement DeleteBranch(name) in internal/git/branch.go
- [X] T020 [US1] Create branches command skeleton with Cobra in cmd/branches.go
- [X] T021 [US1] Wire GetMergedBranches into branches command in cmd/branches.go
- [X] T022 [US1] Display merged branches summary in cmd/branches.go

**Checkpoint**: User Story 1 complete - can list and delete merged branches ✅

---

## Phase 4: User Story 2 - Clean Branches with Deleted Remotes (Priority: P1)

**Goal**: Remove local branches whose remote tracking branch was deleted

**Independent Test**: Create branch with gone remote, verify detection

### Tests for User Story 2

- [X] T023 [P] [US2] Test GetGoneBranches detects [gone] status in tests/branch_test.go
- [X] T024 [P] [US2] Test branches with active tracking are excluded in tests/branch_test.go

### Implementation for User Story 2

- [X] T025 [US2] Implement GetGoneBranches() in internal/git/branch.go
- [X] T026 [US2] Implement UpdateRemoteRefs() with goroutine in internal/git/remote.go
- [X] T027 [US2] Integrate gone branches detection into branches command in cmd/branches.go
- [X] T028 [US2] Combine merged and gone branches with deduplication in cmd/branches.go
- [X] T029 [US2] Show branch deletion reason in summary output in cmd/branches.go

**Checkpoint**: User Story 2 complete - detects both merged and gone-remote branches ✅

---

## Phase 5: User Story 3 - Interactive Branch Selection (Priority: P1)

**Goal**: Interactive fuzzy-finder selection with multi-select

**Independent Test**: Verify navigation, filtering, multi-select, cancellation

### Tests for User Story 3

- [X] T030 [P] [US3] Test SelectBranches returns selected indices in tests/tui_test.go
- [X] T031 [P] [US3] Test Escape cancellation returns empty selection in tests/tui_test.go

### Implementation for User Story 3

- [X] T032 [US3] Create SelectBranches() with go-fzf in internal/tui/selector.go
- [X] T033 [US3] Configure fzf options (multi-select, prompt, prefixes) in internal/tui/selector.go
- [X] T034 [US3] Handle Escape cancellation in internal/tui/selector.go
- [X] T035 [US3] Integrate selector into branches command in cmd/branches.go
- [X] T036 [US3] Add --all flag to skip interactive selection in cmd/branches.go

**Checkpoint**: User Story 3 complete - full interactive branch selection ✅

---

## Phase 6: User Story 4 - Safe Deletion with Confirmation (Priority: P2)

**Goal**: Confirmation prompts before any deletion

**Independent Test**: Verify confirmation appears, respects y/n, --force skips

### Tests for User Story 4

- [X] T037 [P] [US4] Test ConfirmDeletion prompts and returns response in tests/tui_test.go
- [X] T038 [P] [US4] Test --force flag skips confirmation in tests/cmd_test.go

### Implementation for User Story 4

- [X] T039 [US4] Create ConfirmDeletion(items) prompt in internal/tui/confirm.go
- [X] T040 [US4] Display selected branches before confirmation in cmd/branches.go
- [X] T041 [US4] Implement --force flag to skip confirmation in cmd/branches.go
- [X] T042 [US4] Add flag incompatibility check (-a and -f) in cmd/branches.go
- [X] T043 [US4] Show per-branch deletion success/failure in cmd/branches.go
- [X] T044 [US4] Show final deletion summary with count in cmd/branches.go

**Checkpoint**: User Story 4 complete - safe deletion with confirmation flow ✅

---

## Phase 7: User Story 5 - Clean Stale Tags (Priority: P2)

**Goal**: Remove local tags that don't exist on remote

**Independent Test**: Create local-only tags, verify detection and cleanup

### Tests for User Story 5

- [X] T045 [P] [US5] Test GetStaleTags returns tags not on remote in tests/tag_test.go
- [X] T046 [P] [US5] Test DeleteTag removes local tag in tests/tag_test.go

### Implementation for User Story 5

- [X] T047 [US5] Implement GetLocalTags() in internal/git/tag.go
- [X] T048 [US5] Implement GetRemoteTags() in internal/git/tag.go
- [X] T049 [US5] Implement GetStaleTags() comparing local vs remote in internal/git/tag.go
- [X] T050 [US5] Implement DeleteTag(name) in internal/git/tag.go
- [X] T051 [US5] Create tags command with list/clean subcommands in cmd/tags.go
- [X] T052 [US5] Implement tags list subcommand in cmd/tags.go
- [X] T053 [US5] Implement tags clean subcommand with selection in cmd/tags.go
- [X] T054 [US5] Add --force and --all flags to tags clean in cmd/tags.go

**Checkpoint**: User Story 5 complete - full tag cleanup functionality ✅

---

## Phase 8: User Story 6 - Delete Unmerged Branches (Priority: P3)

**Goal**: Option to include unmerged branches with extra safety

**Independent Test**: Verify --unmerged flag, (!) marker, DELETE confirmation

### Tests for User Story 6

- [X] T055 [P] [US6] Test --unmerged flag includes unmerged branches in tests/cmd_test.go
- [X] T056 [P] [US6] Test TypedConfirmation requires exact match in tests/tui_test.go

### Implementation for User Story 6

- [X] T057 [US6] Implement GetAllLocalBranches() in internal/git/branch.go
- [X] T058 [US6] Add --unmerged flag to branches command in cmd/branches.go
- [X] T059 [US6] Mark unmerged branches with (!) prefix in cmd/branches.go
- [X] T060 [US6] Create TypedConfirmation(prompt, expected) for DELETE in internal/tui/confirm.go
- [X] T061 [US6] Require typed DELETE confirmation for unmerged branches in cmd/branches.go
- [X] T062 [US6] Implement DeleteBranchWithRemote() in internal/git/branch.go
- [X] T063 [US6] Delete both local and remote for unmerged branches in cmd/branches.go

**Checkpoint**: User Story 6 complete - unmerged branch cleanup with safety ✅

---

## Phase 9: User Story 7 - Verbose Progress Feedback (Priority: P3)

**Goal**: Clear, friendly feedback throughout the process

**Independent Test**: Verify emoji messages, progress, summaries

### Tests for User Story 7

- [X] T064 [P] [US7] Test emoji prefixes are consistent in output in tests/output_test.go

### Implementation for User Story 7

- [X] T065 [US7] Create Spinner component with Bubbletea in internal/tui/spinner.go
- [X] T066 [US7] Add spinner during remote fetch operations in cmd/branches.go
- [X] T067 [US7] Standardize emoji prefixes (✅❌⚠️🔄📍🎉) in internal/tui/styles.go
- [X] T068 [US7] Add verbose summary with category counts in cmd/branches.go
- [X] T069 [US7] Add (!) Unmerged legend when -u flag used in cmd/branches.go

**Checkpoint**: User Story 7 complete - polished user experience ✅

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Integration Tests

- [X] T070 [P] Test "not in a git repository" error in tests/integration_test.go
- [X] T071 [P] Test "no branches to delete" message in tests/integration_test.go
- [X] T072 [P] Test partial deletion failure handling in tests/integration_test.go
- [X] T073 [P] Test remote unreachable warning in tests/integration_test.go
- [X] T074 [P] Test non-English locale handling (LC_ALL=C) in tests/integration_test.go

### Verification Tasks

- [X] T075 [P] Update root.go with all global flags in cmd/root.go
- [X] T076 [P] Ensure version command shows commit and build time in cmd/version.go
- [X] T077 [P] Verify self-update command functionality in cmd/selfupdate.go
- [X] T078 [P] Add LC_ALL=C to all git command executions
- [X] T079 Verify git plugin mode works (git gone vs git-gone)
- [X] T080 Test with different default branch names (main, master, develop)
- [X] T081 Run go fmt and go vet on all files
- [X] T082 Update README.md with new commands documentation
- [X] T083 Run quickstart.md validation manually

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup) ──────────────────────────────────────────┐
         │                                                │
         ▼                                                │
Phase 2 (Foundational) ◄──────── BLOCKS ALL ─────────────┤
         │                                                │
         ├──────────────────────────────────────┐         │
         ▼                                      ▼         │
Phase 3 (US1)                            Phase 4 (US2)   │
         │                                      │         │
         └──────────────┬───────────────────────┘         │
                        ▼                                 │
                 Phase 5 (US3) ◄── depends on US1+US2    │
                        │                                 │
                        ▼                                 │
                 Phase 6 (US4) ◄── depends on US3        │
                        │                                 │
         ┌──────────────┴──────────────┐                 │
         ▼                             ▼                 │
  Phase 7 (US5)                 Phase 8 (US6)           │
         │                             │                 │
         └──────────────┬──────────────┘                 │
                        ▼                                 │
                 Phase 9 (US7)                           │
                        │                                 │
                        ▼                                 │
                 Phase 10 (Polish)                       │
```

### User Story Dependencies

| Story | Depends On | Can Parallelize With |
|-------|------------|---------------------|
| US1 | Foundational | US2 |
| US2 | Foundational | US1 |
| US3 | US1, US2 | - |
| US4 | US3 | - |
| US5 | US4 (uses confirmation TUI) | US6 |
| US6 | US4 (uses confirmation TUI) | US5 |
| US7 | US5, US6 | - |

### Within Each User Story

- Models/types before business logic
- Business logic before CLI integration
- Core implementation before flags/options
- Story complete before moving to next priority

### Parallel Opportunities

**Phase 1 (Setup)**:
```
T004 + T005 (dependencies)
```

**Phase 2 (Foundational)**:
```
T007 + T008 + T009 + T010 (types and styles)
```

**Phase 3-4 (US1 + US2)** - Tests and implementation can run in parallel:
```
US1 Tests: T014 + T015 + T016 (parallel)
US1 Impl:  T017 → T018 → T019 → T020 → T021 → T022
US2 Tests: T023 + T024 (parallel)
US2 Impl:  T025 → T026 → T027 → T028 → T029
```

**Phase 10 (Polish)**:
```
T070-T078 (integration tests and verification - parallel)
```

---

## Implementation Strategy

### MVP First (User Stories 1-4)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (merged branches)
4. Complete Phase 4: User Story 2 (gone remotes)
5. Complete Phase 5: User Story 3 (interactive selection)
6. Complete Phase 6: User Story 4 (confirmation)
7. **STOP and VALIDATE**: MVP is functional
8. Deploy/demo if ready

### Incremental Delivery

| Milestone | Stories | Value Delivered |
|-----------|---------|-----------------|
| MVP | US1-US4 | Branch cleanup with safety |
| v1.1 | +US5 | Tag cleanup |
| v1.2 | +US6 | Unmerged branch cleanup |
| v1.3 | +US7 | Polished UX |

### Task Summary

| Phase | Tasks | Tests | Parallel |
|-------|-------|-------|----------|
| Setup | 5 | 0 | 2 |
| Foundational | 8 | 0 | 5 |
| US1 (P1) | 9 | 3 | 3 |
| US2 (P1) | 7 | 2 | 2 |
| US3 (P1) | 7 | 2 | 2 |
| US4 (P2) | 8 | 2 | 2 |
| US5 (P2) | 10 | 2 | 2 |
| US6 (P3) | 9 | 2 | 2 |
| US7 (P3) | 6 | 1 | 1 |
| Polish | 14 | 5 | 9 |
| **Total** | **83** | **19** | **30** |

---

## Notes

- [P] tasks = different files, no dependencies on incomplete tasks
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Existing code in cmd/ provides foundation - extend rather than rewrite

