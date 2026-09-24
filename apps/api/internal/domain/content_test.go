package domain

import (
	"errors"
	"testing"
	"time"
)

func TestContentValidation(t *testing.T) {
	from := time.Date(2026, 9, 24, 10, 0, 0, 0, time.UTC)
	until := from.Add(time.Hour)

	tests := []struct {
		name  string
		check func() error
		valid bool
	}{
		{"module title and range", func() error {
			return ValidateCreateModule(CreateModuleInput{Title: "Week 1", AvailableFrom: &from, AvailableUntil: &until})
		}, true},
		{"module invalid range", func() error {
			return ValidateCreateModule(CreateModuleInput{Title: "Week 1", AvailableFrom: &until, AvailableUntil: &from})
		}, false},
		{"lesson learning mode", func() error {
			return ValidateCreateLesson(CreateLessonInput{Title: "Lecture", LearningMode: "asynchronous"})
		}, true},
		{"lesson invalid learning mode", func() error {
			return ValidateCreateLesson(CreateLessonInput{Title: "Lecture", LearningMode: "unknown"})
		}, false},
		{"text material", func() error {
			return ValidateCreateMaterial(CreateMaterialInput{Title: "Notes", Type: "text", Content: "content"})
		}, true},
		{"external material types", func() error {
			return ValidateCreateMaterial(CreateMaterialInput{Title: "Interactive", Type: "embed", ExternalURL: "https://example.test/embed"})
		}, true},
		{"file material missing file", func() error {
			return ValidateCreateMaterial(CreateMaterialInput{Title: "Slides", Type: "file"})
		}, false},
		{"file metadata", func() error {
			return ValidateCreateFile(CreateFileInput{OriginalFilename: "slides.pdf", MIMEType: "application/pdf", SizeBytes: 10, Checksum: "sha256:abc"})
		}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.check()
			if test.valid && err != nil {
				t.Fatalf("validation error = %v", err)
			}
			if !test.valid && !errors.Is(err, ErrInvalidContent) {
				t.Fatalf("validation error = %v, want ErrInvalidContent", err)
			}
		})
	}
}
