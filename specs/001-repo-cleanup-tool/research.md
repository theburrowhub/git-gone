# Research: Git Repository Cleanup Tool

**Date**: 2025-12-16  
**Feature**: 001-repo-cleanup-tool

## Technology Decisions

### 1. CLI Framework: Cobra

**Decision**: Use `github.com/spf13/cobra` for all CLI functionality

**Rationale**:
- Industry standard for Go CLI applications
- Built-in support for subcommands, flags, help generation
- Automatic shell completion (bash, zsh, fish, powershell)
- Excellent documentation and community support
- Already used in current codebase

**Alternatives Considered**:
- `urfave/cli`: Good but less feature-rich for complex subcommand structures
- `kong`: Struct-tag based, different paradigm from existing code
- Standard `flag`: Too low-level for command/subcommand pattern

### 2. TUI Framework: Bubbletea

**Decision**: Use `github.com/charmbracelet/bubbletea` for interactive TUI elements

**Rationale**:
- Elm-architecture makes state management predictable
- Rich ecosystem (bubbles for components, lipgloss for styling)
- Excellent performance and responsiveness
- Active development and community
- Works well with goroutines for async operations

**Alternatives Considered**:
- `go-fzf` only: Limited to fuzzy selection, can't do rich TUI
- `termui`: Less active, older architecture
- `tview`: More complex, overkill for this use case

**Integration Strategy**:
- Keep `go-fzf` for simple fuzzy selection (already working)
- Use Bubbletea for enhanced TUI when more control needed (spinners, progress)
- Gradual migration path: enhance TUI over time

### 3. Self-Update Mechanism

**Decision**: Implement GitHub release-based self-update

**Rationale**:
- GitHub Releases API is well-documented and reliable
- Users expect tools to self-update
- Single binary makes update simple (replace binary)

**Implementation Pattern**:
```
1. Check GitHub API for latest release
2. Compare with current version (semantic versioning)
3. If newer available, download appropriate binary
4. Verify checksum
5. Replace current binary
6. Report success/failure
```

**Libraries**:
- `github.com/blang/semver` - Semantic version comparison
- `net/http` - GitHub API calls (no external dependency needed)

### 4. Concurrency Pattern: Goroutines

**Decision**: Use goroutines with channels for async git operations

**Rationale**:
- Git operations (fetch, branch listing) can be slow
- User shouldn't wait for operations that can run in background
- Go's concurrency model is well-suited for this

**Pattern**:
```go
// Background fetch while showing TUI
fetchDone := make(chan error)
go func() {
    fetchDone <- git.FetchAll()
}()

// Show UI immediately, update when fetch completes
```

**Key Async Operations**:
- `git fetch --all --prune` (network I/O)
- Remote branch status checking
- Self-update download

### 5. Git Operations

**Decision**: Use `os/exec` to call git commands directly

**Rationale**:
- Most reliable way to get consistent git behavior
- Avoids libgit2 CGO dependencies
- Single binary remains simple
- `LC_ALL=C` ensures consistent output parsing

**Alternatives Considered**:
- `go-git`: Pure Go but doesn't support all git features
- `libgit2/git2go`: CGO dependency, complicates build

## Best Practices Applied

### Cobra Best Practices

1. **Root command**: Define persistent flags here (--verbose, --force)
2. **Subcommands**: Each in separate file (`branches.go`, `tags.go`)
3. **Init pattern**: Use `init()` to register subcommands
4. **Help customization**: Use Cobra's template system
5. **Completion**: Enable auto-generated shell completions

### Bubbletea Best Practices

1. **Model-View-Update**: Keep state in model, pure update functions
2. **Commands**: Return Cmd from Update for async operations
3. **Styling**: Use lipgloss for consistent, cross-platform styling
4. **Components**: Use bubbles library for standard components

### Performance Best Practices

1. **Lazy loading**: Don't fetch data until needed
2. **Streaming**: Process branches as found, don't load all into memory
3. **Cancellation**: Respect Ctrl+C, clean up goroutines
4. **Caching**: Cache git operations within single run

## Integration Points

### Git Plugin Integration

Binary named `git-gone` automatically becomes `git gone` when in PATH.

```bash
# Both work identically
git-gone branches
git gone branches
```

### GitHub Releases

Release naming convention: `v{MAJOR}.{MINOR}.{PATCH}`

Asset naming: `git-gone-{os}-{arch}` (e.g., `git-gone-darwin-arm64`)

## Open Questions Resolved

| Question | Resolution |
|----------|------------|
| Keep go-fzf or replace with Bubbletea? | Keep both - go-fzf for simple selection, Bubbletea for enhanced TUI |
| How to handle concurrent git operations? | Goroutines with channels, background fetch while showing UI |
| Self-update security? | Checksum verification, HTTPS only |

