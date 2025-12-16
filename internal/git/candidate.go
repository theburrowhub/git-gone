package git

// CandidateType represents the type of deletion candidate.
type CandidateType int

const (
	// CandidateBranch represents a branch candidate.
	CandidateBranch CandidateType = iota
	// CandidateTag represents a tag candidate.
	CandidateTag
)

// DeletionReason represents why an item is a deletion candidate.
type DeletionReason int

const (
	// ReasonMerged means the branch was merged into default.
	ReasonMerged DeletionReason = iota
	// ReasonGoneRemote means the remote tracking branch was deleted.
	ReasonGoneRemote
	// ReasonStaleTag means the tag doesn't exist on remote.
	ReasonStaleTag
	// ReasonUnmerged means the branch is not merged (dangerous).
	ReasonUnmerged
)

// RiskLevel represents how dangerous a deletion is.
type RiskLevel int

const (
	// RiskSafe means deletion can be confirmed with y/n.
	RiskSafe RiskLevel = iota
	// RiskDangerous means deletion requires typing "DELETE".
	RiskDangerous
)

// DeletionCandidate represents an item that can be deleted.
type DeletionCandidate struct {
	Type         CandidateType
	Name         string
	Reason       DeletionReason
	RiskLevel    RiskLevel
	DisplayLabel string
}

// UnmergedPrefix is the prefix used to mark unmerged branches.
const UnmergedPrefix = "(!) "

// NewBranchCandidate creates a DeletionCandidate for a branch.
func NewBranchCandidate(name string, reason DeletionReason) DeletionCandidate {
	risk := RiskSafe
	displayLabel := name
	if reason == ReasonUnmerged {
		risk = RiskDangerous
		displayLabel = UnmergedPrefix + name
	}
	return DeletionCandidate{
		Type:         CandidateBranch,
		Name:         name,
		Reason:       reason,
		RiskLevel:    risk,
		DisplayLabel: displayLabel,
	}
}

// NewTagCandidate creates a DeletionCandidate for a tag.
func NewTagCandidate(name string) DeletionCandidate {
	return DeletionCandidate{
		Type:         CandidateTag,
		Name:         name,
		Reason:       ReasonStaleTag,
		RiskLevel:    RiskSafe,
		DisplayLabel: name,
	}
}
