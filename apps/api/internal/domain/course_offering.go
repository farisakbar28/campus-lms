// Package domain contains transport- and database-independent LMS entities.
package domain

import "errors"

// ErrNotFound avoids revealing whether an offering exists but is inaccessible.
var ErrNotFound = errors.New("course offering not found")

// CourseOffering identifies the SIAKAD-authoritative offering being read.
type CourseOffering struct {
	ID          string
	ExternalID  string
	SectionCode string
	DisplayName string
	LMSStatus   string
	Course      Course
	Term        AcademicTerm
}

type Course struct {
	ID   string
	Code string
	Name string
}

type AcademicTerm struct {
	ID   string
	Code string
	Name string
}

type Participant struct {
	StudentUserID    string
	DisplayName      string
	EnrollmentStatus string
}

type Roster struct {
	Offering     CourseOffering
	Participants []Participant
}
