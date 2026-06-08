<!--
  SYNC IMPACT REPORT
  ==================
  Version change: N/A → 1.0.0 (initial creation)
  
  Added Principles:
  - I. Code Quality
  - II. Testing Standards
  - III. User Experience Consistency
  - IV. Performance Requirements
  
  Added Sections:
  - Technical Constraints
  - Development Workflow
  - Governance
  
  Removed Sections: None (initial creation)
  
  Templates Validation:
  - .specify/templates/plan-template.md: ✅ Compatible (Constitution Check section exists)
  - .specify/templates/spec-template.md: ✅ Compatible (success criteria align with performance principle)
  - .specify/templates/tasks-template.md: ✅ Compatible (testing phases align with testing standards)
  
  Follow-up TODOs: None
-->

# git-gone Constitution

## Core Principles

### I. Code Quality

All code MUST adhere to these non-negotiable standards:

- **Functional paradigm preferred**: Use pure functions, avoid side effects where possible, prefer immutability
- **Single responsibility**: Each function/module MUST have one clear purpose
- **Explicit error handling**: All errors MUST be handled explicitly; no silent failures
- **No dead code**: Remove unused code, commented-out blocks, and orphaned functions
- **Consistent formatting**: All code MUST pass `go fmt` and follow Go idioms
- **Self-documenting code**: Code MUST be readable without extensive comments; comments explain "why", not "what"

**Rationale**: git-gone is a developer tool; developers expect clean, maintainable code they can trust and contribute to.

### II. Testing Standards

Testing is MANDATORY for all user-facing functionality:

- **TDD encouraged**: Write tests before implementation when adding new features
- **Test coverage**: All CLI commands MUST have integration tests verifying expected behavior
- **Edge cases**: Error paths, empty inputs, and boundary conditions MUST be tested
- **No flaky tests**: Tests MUST be deterministic and reproducible
- **Test naming**: Test names MUST describe the scenario being tested (e.g., `TestDeleteBranch_WithUnmergedBranch_RequiresConfirmation`)

**Rationale**: Users trust git-gone to manage their branches safely; comprehensive testing ensures reliability.

### III. User Experience Consistency

The CLI experience MUST be predictable and intuitive:

- **Consistent output format**: All messages MUST use emoji prefixes for status (✅ success, ❌ error, ⚠️ warning, 🔄 progress, 📍 info)
- **Clear feedback**: Every user action MUST produce visible feedback
- **Safe defaults**: Destructive operations MUST require explicit confirmation unless `--force` is used
- **Graceful degradation**: Tool MUST handle missing git repos, network issues, and invalid states with clear error messages
- **No surprises**: Flag behavior MUST be documented and consistent across all commands
- **Internationalization-ready**: Use `LC_ALL=C` for git commands to ensure consistent parsing regardless of user locale

**Rationale**: Developers use git-gone in their daily workflow; a consistent, predictable experience builds trust.

### IV. Performance Requirements

The tool MUST remain fast and responsive:

- **Startup time**: CLI MUST start in under 100ms on standard hardware
- **Branch operations**: Listing and filtering branches MUST complete in under 1 second for repositories with up to 1000 branches
- **Memory efficiency**: Tool MUST not load unnecessary data; stream where possible
- **No blocking UI**: Interactive selection MUST remain responsive during background operations
- **Minimal dependencies**: Only add dependencies that provide significant value; prefer standard library

**Rationale**: Developers run git-gone frequently; slow tools disrupt workflow and reduce adoption.

## Technical Constraints

The following technical decisions are binding:

- **Language**: Go (minimum version specified in go.mod)
- **CLI Framework**: Cobra for command structure
- **Interactive Selection**: go-fzf for fuzzy finding
- **Build**: Single binary with no runtime dependencies
- **Platforms**: MUST support Linux, macOS, and Windows
- **Git Integration**: MUST work as both standalone (`git-gone`) and git plugin (`git gone`)

## Development Workflow

All contributions MUST follow this workflow:

1. **Branch naming**: `feat-*`, `fix-*`, `docs-*`, `refactor-*` prefixes required
2. **Atomic commits**: Each commit MUST represent a single logical change
3. **Build verification**: Code MUST compile without warnings before commit
4. **Lint check**: Code MUST pass `go vet` and `go fmt` checks
5. **Test verification**: All tests MUST pass before merging
6. **Version bumping**: Follow semantic versioning (MAJOR.MINOR.PATCH)

## Governance

This constitution supersedes all other development practices for git-gone:

- **Compliance**: All PRs and code reviews MUST verify adherence to these principles
- **Amendments**: Changes to this constitution require:
  - Documentation of the proposed change
  - Justification for why current principles are insufficient
  - Version increment following semantic versioning
- **Exceptions**: Complexity beyond these principles MUST be explicitly justified in PR description
- **Runtime guidance**: See `docs/ARCHITECTURE.md` for implementation patterns and architectural decisions

**Version**: 1.0.0 | **Ratified**: 2025-12-16 | **Last Amended**: 2025-12-16
