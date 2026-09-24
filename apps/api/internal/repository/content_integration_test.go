package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/jackc/pgx/v5"
)

func TestContentAuthorizationPublicationAndFileVisibility(t *testing.T) {
	ctx := context.Background()
	if _, err := repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET lms_status = 'published', published_at = now(), archived_at = NULL WHERE id = $1::uuid`, offeringAID); err != nil {
		t.Fatalf("publish content offering fixture: %v", err)
	}
	defer func() {
		_, _ = repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET lms_status = 'published', published_at = NULL, archived_at = NULL WHERE id = $1::uuid`, offeringAID)
	}()

	var studentID string
	if err := repositorySuite.owner.QueryRow(ctx, `
SELECT e.student_user_id
FROM enrollments AS e
JOIN memberships AS m ON m.tenant_id = e.tenant_id AND m.user_id = e.student_user_id
JOIN membership_roles AS mr ON mr.tenant_id = m.tenant_id AND mr.membership_id = m.id
WHERE e.tenant_id = $1::uuid AND e.course_offering_id = $2::uuid AND e.status = 'active'
  AND m.status = 'active' AND mr.role = 'student' AND mr.revoked_at IS NULL
LIMIT 1`, tenantAID, offeringAID).Scan(&studentID); err != nil {
		t.Fatalf("select content student fixture: %v", err)
	}

	service := NewContentService(repositorySuite.appPool)
	availableFrom := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	availableUntil := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	module, err := service.CreateModule(ctx, tenantAID, repositorySuite.instructorA, offeringAID, "request-module", domain.CreateModuleInput{Title: "Week 1", Position: 1, AvailableFrom: &availableFrom, AvailableUntil: &availableUntil})
	if err != nil {
		t.Fatalf("create module: %v", err)
	}
	updatedModule, err := service.UpdateModule(ctx, tenantAID, repositorySuite.instructorA, offeringAID, module.ID, "request-module-window", domain.UpdateModuleInput{AvailableFrom: &availableFrom, AvailableUntil: &availableUntil})
	if err != nil {
		t.Fatalf("update module availability: %v", err)
	}
	if updatedModule.AvailableFrom == nil || !updatedModule.AvailableFrom.Equal(availableFrom) || updatedModule.AvailableUntil == nil || !updatedModule.AvailableUntil.Equal(availableUntil) {
		t.Fatalf("updated module availability = %#v/%#v, want %s/%s", updatedModule.AvailableFrom, updatedModule.AvailableUntil, availableFrom, availableUntil)
	}
	lesson, err := service.CreateLesson(ctx, tenantAID, repositorySuite.instructorA, offeringAID, module.ID, "request-lesson", domain.CreateLessonInput{Title: "Lecture", LearningMode: "asynchronous", Position: 1, EstimatedMinutes: 30, AvailableFrom: &availableFrom, AvailableUntil: &availableUntil})
	if err != nil {
		t.Fatalf("create lesson: %v", err)
	}
	updatedLesson, err := service.UpdateLesson(ctx, tenantAID, repositorySuite.instructorA, offeringAID, lesson.ID, "request-lesson-window", domain.UpdateLessonInput{AvailableFrom: &availableFrom, AvailableUntil: &availableUntil})
	if err != nil {
		t.Fatalf("update lesson availability: %v", err)
	}
	if updatedLesson.AvailableFrom == nil || !updatedLesson.AvailableFrom.Equal(availableFrom) || updatedLesson.AvailableUntil == nil || !updatedLesson.AvailableUntil.Equal(availableUntil) {
		t.Fatalf("updated lesson availability = %#v/%#v, want %s/%s", updatedLesson.AvailableFrom, updatedLesson.AvailableUntil, availableFrom, availableUntil)
	}
	textMaterial, err := service.CreateMaterial(ctx, tenantAID, repositorySuite.instructorA, offeringAID, lesson.ID, "request-material", domain.CreateMaterialInput{Title: "Notes", Type: "text", Content: "Week 1 notes", Position: 1})
	if err != nil {
		t.Fatalf("create text material: %v", err)
	}
	file, err := service.CreateFile(ctx, tenantAID, repositorySuite.instructorA, offeringAID, "request-file", domain.CreateFileInput{OriginalFilename: "slides.pdf", MIMEType: "application/pdf", SizeBytes: 42, Checksum: "sha256:test"})
	if err != nil {
		t.Fatalf("create file metadata: %v", err)
	}
	fileMaterial, err := service.CreateMaterial(ctx, tenantAID, repositorySuite.instructorA, offeringAID, lesson.ID, "request-file-material", domain.CreateMaterialInput{Title: "Slides", Type: "file", FileID: file.ID, Position: 2})
	if err != nil {
		t.Fatalf("create file material: %v", err)
	}

	studentContent, err := service.ListContent(ctx, tenantAID, studentID, offeringAID, time.Now().UTC())
	if err != nil {
		t.Fatalf("list unpublished student content: %v", err)
	}
	if len(studentContent.Modules) != 0 {
		t.Fatalf("student saw %d unpublished modules, want 0", len(studentContent.Modules))
	}

	if _, err := service.UpdateModule(ctx, tenantAID, repositorySuite.instructorA, offeringAID, module.ID, "request-denied", domain.UpdateModuleInput{Status: stringPointer("published")}); !errors.Is(err, domain.ErrContentConflict) {
		t.Fatalf("instructor publication error = %v, want content conflict", err)
	}
	if _, err := service.UpdateModule(ctx, tenantAID, fixtureLead, offeringAID, module.ID, "request-module-publish", domain.UpdateModuleInput{Status: stringPointer("published")}); err != nil {
		t.Fatalf("lead publishes module: %v", err)
	}
	if _, err := service.UpdateLesson(ctx, tenantAID, fixtureLead, offeringAID, lesson.ID, "request-lesson-publish", domain.UpdateLessonInput{Status: stringPointer("published")}); err != nil {
		t.Fatalf("lead publishes lesson: %v", err)
	}
	if _, err := service.UpdateMaterial(ctx, tenantAID, fixtureLead, offeringAID, textMaterial.ID, "request-text-publish", domain.UpdateMaterialInput{Published: boolPointer(true)}); err != nil {
		t.Fatalf("lead publishes text material: %v", err)
	}
	if _, err := service.UpdateMaterial(ctx, tenantAID, fixtureLead, offeringAID, fileMaterial.ID, "request-file-publish", domain.UpdateMaterialInput{Published: boolPointer(true)}); err != nil {
		t.Fatalf("lead publishes file material: %v", err)
	}

	studentContent, err = service.ListContent(ctx, tenantAID, studentID, offeringAID, time.Now().UTC())
	if err != nil {
		t.Fatalf("list published student content: %v", err)
	}
	if len(studentContent.Modules) != 1 || len(studentContent.Modules[0].Lessons) != 1 || len(studentContent.Modules[0].Lessons[0].Materials) != 1 {
		t.Fatalf("student content tree = %#v, want one module, lesson, and clean material", studentContent)
	}
	if studentContent.Modules[0].Lessons[0].Materials[0].ID != textMaterial.ID {
		t.Fatalf("student material = %s, want text material %s", studentContent.Modules[0].Lessons[0].Materials[0].ID, textMaterial.ID)
	}

	if _, err := repositorySuite.owner.Exec(ctx, `UPDATE files SET malware_scan_status = 'clean' WHERE id = $1::uuid AND tenant_id = $2::uuid`, file.ID, tenantAID); err != nil {
		t.Fatalf("mark file clean: %v", err)
	}
	studentContent, err = service.ListContent(ctx, tenantAID, studentID, offeringAID, time.Now().UTC())
	if err != nil {
		t.Fatalf("list clean file content: %v", err)
	}
	if got := len(studentContent.Modules[0].Lessons[0].Materials); got != 2 {
		t.Fatalf("student materials after clean scan = %d, want 2", got)
	}
	for _, material := range studentContent.Modules[0].Lessons[0].Materials {
		if material.File != nil && material.File.StorageKey != "" {
			t.Fatalf("student file exposed storage key %q", material.File.StorageKey)
		}
	}

	staffContent, err := service.ListContent(ctx, tenantAID, repositorySuite.instructorA, offeringAID, time.Now().UTC())
	if err != nil {
		t.Fatalf("list staff content: %v", err)
	}
	if got := len(staffContent.Modules[0].Lessons[0].Materials); got != 2 {
		t.Fatalf("staff materials = %d, want 2", got)
	}

	if _, err := service.ListContent(ctx, tenantAID, repositorySuite.instructorA, offeringBID, time.Now().UTC()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("cross-tenant offering error = %v, want not found", err)
	}
	if _, err := service.CreateModule(ctx, tenantAID, repositorySuite.instructorA, offeringAID, "request-archived", domain.CreateModuleInput{Title: "Archived", Position: 3}); err != nil {
		t.Fatalf("create pre-archive module: %v", err)
	}
	if _, err := repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET lms_status = 'archived', archived_at = now() WHERE id = $1::uuid`, offeringAID); err != nil {
		t.Fatalf("archive content offering fixture: %v", err)
	}
	if _, err := service.CreateModule(ctx, tenantAID, repositorySuite.instructorA, offeringAID, "request-archived-denied", domain.CreateModuleInput{Title: "Denied", Position: 4}); !errors.Is(err, domain.ErrContentConflict) {
		t.Fatalf("archived create error = %v, want content conflict", err)
	}
	if _, err := repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET lms_status = 'unexpected', published_at = now(), archived_at = NULL WHERE id = $1::uuid`, offeringAID); err != nil {
		t.Fatalf("set invalid offering lifecycle fixture: %v", err)
	}
	if _, err := service.ListContent(ctx, tenantAID, studentID, offeringAID, time.Now().UTC()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("student invalid lifecycle error = %v, want not found", err)
	}
	if _, err := service.ListContent(ctx, tenantAID, repositorySuite.instructorA, offeringAID, time.Now().UTC()); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("staff invalid lifecycle error = %v, want not found", err)
	}
}

func TestContentListUsesFixedQueryShape(t *testing.T) {
	ctx := context.Background()
	if _, err := repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET lms_status = 'published', published_at = now() WHERE id = $1::uuid`, offeringAID); err != nil {
		t.Fatalf("publish query-count fixture: %v", err)
	}
	defer func() {
		_, _ = repositorySuite.owner.Exec(ctx, `UPDATE course_offerings SET published_at = NULL WHERE id = $1::uuid`, offeringAID)
	}()

	count := 0
	err := repositorySuite.appPool.WithTenantTx(ctx, tenantAID, func(tx pgx.Tx) error {
		_, err := (ContentRepository{}).ListContent(ctx, countingQuerier{Querier: tx, count: &count}, tenantAID, repositorySuite.instructorA, offeringAID, time.Now().UTC())
		return err
	})
	if err != nil {
		t.Fatalf("list content: %v", err)
	}
	if count != 4 {
		t.Fatalf("content list query count = %d, want fixed 4-query shape", count)
	}
}

func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }
