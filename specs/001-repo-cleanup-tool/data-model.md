# Data Model: Git Repository Cleanup Tool

**Date**: 2025-12-16  
**Feature**: 001-repo-cleanup-tool

## Core Entities

### Branch

Represents a local git branch with its associated metadata.

```
Branch
├── Name: string              # Branch name (e.g., "feature/login")
├── IsCurrent: bool           # Is this the currently checked out branch?
├── IsDefault: bool           # Is this the default branch (main/master)?
├── IsMerged: bool            # Has this been merged into default branch?
├── RemoteStatus: enum        # TrackingActive | TrackingGone | NoTracking
└── RemoteName: string?       # Remote tracking branch name (if any)
```

**States**:
- `Safe`: IsMerged=true OR RemoteStatus=TrackingGone
- `Dangerous`: IsMerged=false AND RemoteStatus≠TrackingGone
- `Protected`: IsCurrent=true OR IsDefault=true

**Validation Rules**:
- Name cannot be empty
- Protected branches cannot be deletion candidates

### Tag

Represents a local git tag.

```
Tag
├── Name: string              # Tag name (e.g., "v1.0.0")
├── ExistsOnRemote: bool      # Does this tag exist on origin?
└── IsAnnotated: bool         # Is this an annotated tag?
```

**States**:
- `Stale`: ExistsOnRemote=false
- `Current`: ExistsOnRemote=true

### DeletionCandidate

A unified representation of items that can be deleted.

```
DeletionCandidate
├── Type: enum                # Branch | Tag
├── Name: string              # Display name
├── Reason: enum              # Merged | GoneRemote | StaleTag | Unmerged
├── RiskLevel: enum           # Safe | Dangerous
└── DisplayLabel: string      # Formatted for TUI (e.g., "(!) feature/old")
```

**Display Prefixes by Reason**:
- `Merged`: (none) - safe, default
- `GoneRemote`: (none) - safe, remote was deleted
- `StaleTag`: (none) - safe, tag doesn't exist on remote
- `Unmerged`: `(!) ` - dangerous, requires extra confirmation

### Repository

Represents the git repository context.

```
Repository
├── Path: string              # Absolute path to .git
├── DefaultBranch: string     # main, master, or detected default
├── CurrentBranch: string     # Currently checked out branch
├── HasRemote: bool           # Does origin remote exist?
└── RemoteURL: string?        # Origin URL (for display)
```

## State Transitions

### Branch Lifecycle

```
                    ┌─────────────────┐
                    │   Not Tracked   │
                    └────────┬────────┘
                             │ git push -u
                             ▼
                    ┌─────────────────┐
                    │ Tracking Active │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │ merge                       │ remote deleted
              ▼                             ▼ (squash/rebase)
     ┌─────────────────┐           ┌─────────────────┐
     │     Merged      │           │  Tracking Gone  │
     └────────┬────────┘           └────────┬────────┘
              │                             │
              └──────────────┬──────────────┘
                             │ user deletes
                             ▼
                    ┌─────────────────┐
                    │     Deleted     │
                    └─────────────────┘
```

### Deletion Flow

```
Candidate Found → Selected by User → Confirmation → Deleted
                                          │
                                          │ (if Dangerous)
                                          ▼
                                  Extra Confirmation
                                  (type "DELETE")
```

## Relationships

```
Repository 1 ──────< * Branch
Repository 1 ──────< * Tag
Branch * ──────────< 1 DeletionCandidate (when eligible)
Tag * ─────────────< 1 DeletionCandidate (when stale)
```

## Enumerations

### RemoteStatus
```
TrackingActive  # Branch has active remote tracking
TrackingGone    # Remote tracking branch was deleted
NoTracking      # Branch has no remote tracking
```

### CandidateType
```
Branch
Tag
```

### DeletionReason
```
Merged          # Branch merged into default
GoneRemote      # Remote tracking branch deleted
StaleTag        # Tag doesn't exist on remote
Unmerged        # Branch not merged (dangerous)
```

### RiskLevel
```
Safe            # Can be deleted with simple y/n
Dangerous       # Requires typing "DELETE" to confirm
```

## Data Flow

```
Git Repository
      │
      ▼ (git commands)
┌─────────────────┐
│  internal/git   │
│  - ListBranches │
│  - ListTags     │
│  - GetDefault   │
└────────┬────────┘
         │
         ▼ ([]Branch, []Tag)
┌─────────────────┐
│   cmd/branches  │
│   cmd/tags      │
│   - Filter      │
│   - Categorize  │
└────────┬────────┘
         │
         ▼ ([]DeletionCandidate)
┌─────────────────┐
│   internal/tui  │
│   - Selector    │
│   - Confirm     │
└────────┬────────┘
         │
         ▼ (user selection)
┌─────────────────┐
│   internal/git  │
│   - Delete      │
└─────────────────┘
```

