package ag

import (
	"fmt"
	"time"
)

// TODO(meling) must also handle groups; if group assignment is late; withdraw one slip day per group member.
// TODO(meling) also need to propogate slip day information to frontend to show to user

// UpdateSlipDays updates the number of slipdays for the given assignment/submission.
func (m *Enrollment) UpdateSlipDays(start time.Time, assignment *Assignment, submission *Submission) error {
	if m.GetCourseID() != assignment.GetCourseID() {
		return fmt.Errorf("invariant violation (enrollment.CourseID != assignment.CourseID) (%d != %d)", m.CourseID, assignment.CourseID)
	}
	if assignment.GetID() != submission.GetAssignmentID() {
		return fmt.Errorf("invariant violation (assignment.ID != submission.AssignmentID) (%d != %d)", assignment.ID, submission.AssignmentID)
	}
	sinceDeadline, err := assignment.SinceDeadline(start)
	if err != nil {
		return err
	}
	if !submission.GetApproved() && sinceDeadline > 0 {
		// deadline exceeded; calculate used slipdays for this assignment
		m.updateSlipDays(assignment.GetID(), uint32(sinceDeadline/days))
	}
	return nil
}

// UpdateSlipDays updates the number of slipdays for the given assignment.
func (m *Enrollment) updateSlipDays(assignmentID uint64, slipDays uint32) {
	for _, val := range m.GetUsedSlipDays() {
		if val.AssignmentID == assignmentID {
			val.UsedSlipDays = slipDays
			return
		}
	}
	// not found; add new entry to the slice
	m.UsedSlipDays = append(m.UsedSlipDays, &SlipDays{AssignmentID: assignmentID, UsedSlipDays: slipDays})
}

// TotalSlipDays returns the total number of slipdays used for this enrollment.
func (m Enrollment) TotalSlipDays() uint32 {
	var total uint32
	for _, val := range m.GetUsedSlipDays() {
		total += val.GetUsedSlipDays()
	}
	return total
}

// RemainingSlipDays returns the remaining number of slip days for this
// user/course enrollment. Note that if the returned amount is negative,
// the user has used up all slip days.
func (m Enrollment) RemainingSlipDays() int32 {
	return int32(m.Course.GetSlipDays() - m.TotalSlipDays())
}
