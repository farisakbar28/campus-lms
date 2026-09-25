package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	// ErrInvalidContent marks a request that violates content validation rules.
	ErrInvalidContent = errors.New("invalid content")
	// ErrContentConflict marks a valid request that cannot be applied to the current lifecycle.
	ErrContentConflict = errors.New("content conflict")
)

type Content struct {
	OfferingID string
	Modules    []Module
}

type Module struct {
	ID             string
	OfferingID     string
	Title          string
	Description    string
	Position       int
	Status         string
	AvailableFrom  *time.Time
	AvailableUntil *time.Time
	CreatedBy      string
	Lessons        []Lesson
}

type Lesson struct {
	ID               string
	ModuleID         string
	Title            string
	Description      string
	Position         int
	LearningMode     string
	EstimatedMinutes int
	AvailableFrom    *time.Time
	AvailableUntil   *time.Time
	Status           string
	CreatedBy        string
	Materials        []Material
}

type Material struct {
	ID          string
	LessonID    string
	Title       string
	Description string
	Type        string
	File        *File
	ExternalURL string
	Content     string
	Position    int
	Published   bool
	CreatedBy   string
}

type File struct {
	ID                string
	StorageKey        string
	OriginalFilename  string
	MIMEType          string
	SizeBytes         int64
	Checksum          string
	MalwareScanStatus string
	UploadedBy        string
}

type CreateModuleInput struct {
	Title          string
	Description    string
	Position       int
	AvailableFrom  *time.Time
	AvailableUntil *time.Time
}

type UpdateModuleInput struct {
	Title          *string
	Description    *string
	Position       *int
	Status         *string
	AvailableFrom  OptionalTime
	AvailableUntil OptionalTime
}

// OptionalTime distinguishes an omitted patch field from an explicit null.
type OptionalTime struct {
	Set   bool
	Value *time.Time
}

type CreateLessonInput struct {
	Title            string
	Description      string
	Position         int
	LearningMode     string
	EstimatedMinutes int
	AvailableFrom    *time.Time
	AvailableUntil   *time.Time
}

type UpdateLessonInput struct {
	Title            *string
	Description      *string
	Position         *int
	LearningMode     *string
	EstimatedMinutes *int
	Status           *string
	AvailableFrom    OptionalTime
	AvailableUntil   OptionalTime
}

type CreateMaterialInput struct {
	Title       string
	Description string
	Type        string
	FileID      string
	ExternalURL string
	Content     string
	Position    int
}

type UpdateMaterialInput struct {
	Title       *string
	Description *string
	ExternalURL *string
	Content     *string
	Position    *int
	Published   *bool
}

type CreateFileInput struct {
	OriginalFilename string
	MIMEType         string
	SizeBytes        int64
	Checksum         string
}

func ValidateCreateModule(input CreateModuleInput) error {
	if strings.TrimSpace(input.Title) == "" || input.Position < 0 {
		return ErrInvalidContent
	}
	return validateTimeRange(input.AvailableFrom, input.AvailableUntil)
}

func ValidateCreateLesson(input CreateLessonInput) error {
	if strings.TrimSpace(input.Title) == "" || input.Position < 0 || input.EstimatedMinutes < 0 {
		return ErrInvalidContent
	}
	if !validLearningMode(input.LearningMode) {
		return ErrInvalidContent
	}
	return validateTimeRange(input.AvailableFrom, input.AvailableUntil)
}

func ValidateCreateMaterial(input CreateMaterialInput) error {
	if strings.TrimSpace(input.Title) == "" || input.Position < 0 {
		return ErrInvalidContent
	}
	switch input.Type {
	case "text":
		if input.FileID != "" || input.ExternalURL != "" || input.Content == "" {
			return ErrInvalidContent
		}
	case "file":
		if input.FileID == "" || input.ExternalURL != "" || input.Content != "" {
			return ErrInvalidContent
		}
	case "link", "video", "audio", "embed", "learning_package", "external_tool":
		if input.FileID != "" || input.ExternalURL == "" || input.Content != "" {
			return ErrInvalidContent
		}
	default:
		return ErrInvalidContent
	}
	return nil
}

func ValidateCreateFile(input CreateFileInput) error {
	if strings.TrimSpace(input.OriginalFilename) == "" || strings.TrimSpace(input.MIMEType) == "" || input.SizeBytes < 0 || strings.TrimSpace(input.Checksum) == "" {
		return ErrInvalidContent
	}
	return nil
}

func validateTimeRange(from, until *time.Time) error {
	if from != nil && until != nil && !from.Before(*until) {
		return ErrInvalidContent
	}
	return nil
}

func validLearningMode(value string) bool {
	switch value {
	case "asynchronous", "synchronous", "blended", "onsite":
		return true
	default:
		return false
	}
}
