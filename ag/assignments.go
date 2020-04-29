package ag

import (
	"time"
)

const (
	layout = "2006-01-02T15:04:05"
	days   = time.Duration(24 * time.Hour)
	zero   = time.Duration(0)
)

// SinceDeadline returns the duration since the deadline.
// A positive duration means the deadline has passed, whereas
// a negative duration means the deadline has not yet passed.
func (m Assignment) SinceDeadline(now time.Time) (time.Duration, error) {
	deadline, err := time.ParseInLocation(layout, m.GetDeadline(), now.Location())
	if err != nil {
		// this should not happen if deadlines are parsed and recorded correctly
		return zero, err
	}
	return now.Sub(deadline), nil
}

// IsApproved returns true if this assignment is already approved for the
// latest submission, or if the score of the latest submission is sufficient
// to autoapprove the assignment.
func (m Assignment) IsApproved(latest *Submission, score uint32) bool {
	// keep approved status if already approved
	approved := latest.GetApproved()
	if m.GetAutoApprove() && score >= m.GetScoreLimit() {
		approved = true
	}
	return approved
}
