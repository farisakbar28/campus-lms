package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContentService struct {
	transactions TenantTransactioner
	repository   ContentRepository
}

func NewContentService(transactions TenantTransactioner) ContentService {
	return ContentService{transactions: transactions, repository: ContentRepository{}}
}

type ContentRepository struct{}

type contentAccess struct {
	role          string
	offeringState string
}

const authorizeContentSQL = `
SELECT
    o.lms_status,
    o.published_at,
    COALESCE(
        (
            SELECT cs.role
            FROM course_staff AS cs
            JOIN memberships AS m
              ON m.tenant_id = cs.tenant_id
             AND m.user_id = cs.user_id
            JOIN membership_roles AS mr
              ON mr.tenant_id = m.tenant_id
             AND mr.membership_id = m.id
            WHERE cs.tenant_id = o.tenant_id
              AND cs.course_offering_id = o.id
              AND cs.user_id = $3::uuid
              AND cs.active = true
              AND cs.role IN ('instructor', 'lead_instructor')
              AND m.status = 'active'
              AND mr.role = 'lecturer'
              AND mr.revoked_at IS NULL
            LIMIT 1
        ),
        CASE WHEN EXISTS (
            SELECT 1
            FROM enrollments AS e
            JOIN memberships AS m
              ON m.tenant_id = e.tenant_id
             AND m.user_id = e.student_user_id
            JOIN membership_roles AS mr
              ON mr.tenant_id = m.tenant_id
             AND mr.membership_id = m.id
            WHERE e.tenant_id = o.tenant_id
              AND e.course_offering_id = o.id
              AND e.student_user_id = $3::uuid
              AND e.status = 'active'
              AND m.status = 'active'
              AND mr.role = 'student'
              AND mr.revoked_at IS NULL
        ) THEN 'student' ELSE '' END
    ) AS access_role
FROM course_offerings AS o
WHERE o.tenant_id = $1::uuid
  AND o.id = $2::uuid`

const modulesForContentSQL = `
SELECT
    m.id,
    m.course_offering_id,
    m.title,
    m.description,
    m.position,
    m.status,
    m.available_from,
    m.available_until,
    m.created_by
FROM modules AS m
WHERE m.tenant_id = $1::uuid
  AND m.course_offering_id = $2::uuid
  AND (
      $3::text <> 'student'
      OR (
          m.status = 'published'
          AND (m.available_from IS NULL OR m.available_from <= $4::timestamptz)
          AND (m.available_until IS NULL OR m.available_until > $4::timestamptz)
      )
  )
ORDER BY m.position, m.id`

const lessonsForContentSQL = `
SELECT
    l.id,
    l.module_id,
    l.title,
    l.description,
    l.position,
    l.learning_mode,
    l.estimated_minutes,
    l.available_from,
    l.available_until,
    l.status,
    l.created_by
FROM lessons AS l
JOIN modules AS m
  ON m.tenant_id = l.tenant_id
 AND m.id = l.module_id
WHERE l.tenant_id = $1::uuid
  AND m.course_offering_id = $2::uuid
  AND (
      $3::text <> 'student'
      OR (
          m.status = 'published'
          AND l.status = 'published'
          AND (m.available_from IS NULL OR m.available_from <= $4::timestamptz)
          AND (m.available_until IS NULL OR m.available_until > $4::timestamptz)
          AND (l.available_from IS NULL OR l.available_from <= $4::timestamptz)
          AND (l.available_until IS NULL OR l.available_until > $4::timestamptz)
      )
  )
ORDER BY l.module_id, l.position, l.id`

const materialsForContentSQL = `
SELECT
    material.id,
    material.lesson_id,
    material.title,
    material.description,
    material.type,
    material.file_id,
    material.external_url,
    material.content,
    material.position,
    material.published,
    material.created_by,
    file.storage_key,
    file.original_filename,
    file.mime_type,
    file.size_bytes,
    file.checksum,
    file.malware_scan_status,
    file.uploaded_by
FROM materials AS material
JOIN lessons AS lesson
  ON lesson.tenant_id = material.tenant_id
 AND lesson.id = material.lesson_id
JOIN modules AS module
  ON module.tenant_id = lesson.tenant_id
 AND module.id = lesson.module_id
LEFT JOIN files AS file
  ON file.tenant_id = material.tenant_id
 AND file.id = material.file_id
WHERE material.tenant_id = $1::uuid
  AND module.course_offering_id = $2::uuid
  AND (
      $3::text <> 'student'
      OR (
          module.status = 'published'
          AND lesson.status = 'published'
          AND material.published = true
          AND (module.available_from IS NULL OR module.available_from <= $4::timestamptz)
          AND (module.available_until IS NULL OR module.available_until > $4::timestamptz)
          AND (lesson.available_from IS NULL OR lesson.available_from <= $4::timestamptz)
          AND (lesson.available_until IS NULL OR lesson.available_until > $4::timestamptz)
          AND (material.type <> 'file' OR file.malware_scan_status = 'clean')
      )
  )
ORDER BY material.lesson_id, material.position, material.id`

func (s ContentService) ListContent(ctx context.Context, tenantID, userID, offeringID string, now time.Time) (content domain.Content, err error) {
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		content, err = s.repository.ListContent(ctx, tx, tenantID, userID, offeringID, now)
		return err
	})
	return content, err
}

func (s ContentService) CreateModule(ctx context.Context, tenantID, userID, offeringID, requestID string, input domain.CreateModuleInput) (module domain.Module, err error) {
	if err := domain.ValidateCreateModule(input); err != nil {
		return domain.Module{}, err
	}
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		module, err = s.repository.CreateModule(ctx, tx, tenantID, userID, offeringID, requestID, input)
		return err
	})
	return module, err
}

func (s ContentService) UpdateModule(ctx context.Context, tenantID, userID, offeringID, moduleID, requestID string, input domain.UpdateModuleInput) (module domain.Module, err error) {
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		module, err = s.repository.UpdateModule(ctx, tx, tenantID, userID, offeringID, moduleID, requestID, input)
		return err
	})
	return module, err
}

func (s ContentService) CreateLesson(ctx context.Context, tenantID, userID, offeringID, moduleID, requestID string, input domain.CreateLessonInput) (lesson domain.Lesson, err error) {
	if err := domain.ValidateCreateLesson(input); err != nil {
		return domain.Lesson{}, err
	}
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		lesson, err = s.repository.CreateLesson(ctx, tx, tenantID, userID, offeringID, moduleID, requestID, input)
		return err
	})
	return lesson, err
}

func (s ContentService) UpdateLesson(ctx context.Context, tenantID, userID, offeringID, lessonID, requestID string, input domain.UpdateLessonInput) (lesson domain.Lesson, err error) {
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		lesson, err = s.repository.UpdateLesson(ctx, tx, tenantID, userID, offeringID, lessonID, requestID, input)
		return err
	})
	return lesson, err
}

func (s ContentService) CreateMaterial(ctx context.Context, tenantID, userID, offeringID, lessonID, requestID string, input domain.CreateMaterialInput) (material domain.Material, err error) {
	if err := domain.ValidateCreateMaterial(input); err != nil {
		return domain.Material{}, err
	}
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		material, err = s.repository.CreateMaterial(ctx, tx, tenantID, userID, offeringID, lessonID, requestID, input)
		return err
	})
	return material, err
}

func (s ContentService) UpdateMaterial(ctx context.Context, tenantID, userID, offeringID, materialID, requestID string, input domain.UpdateMaterialInput) (material domain.Material, err error) {
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		material, err = s.repository.UpdateMaterial(ctx, tx, tenantID, userID, offeringID, materialID, requestID, input)
		return err
	})
	return material, err
}

func (s ContentService) CreateFile(ctx context.Context, tenantID, userID, offeringID, requestID string, input domain.CreateFileInput) (file domain.File, err error) {
	if err := domain.ValidateCreateFile(input); err != nil {
		return domain.File{}, err
	}
	err = s.transactions.WithTenantTx(ctx, tenantID, func(tx pgx.Tx) error {
		file, err = s.repository.CreateFile(ctx, tx, tenantID, userID, offeringID, requestID, input)
		return err
	})
	return file, err
}

func (ContentRepository) authorizeContent(ctx context.Context, queries Querier, tenantID, userID, offeringID string) (contentAccess, error) {
	var access contentAccess
	var publishedAt pgtype.Timestamptz
	if err := queries.QueryRow(ctx, authorizeContentSQL, tenantID, offeringID, userID).Scan(&access.offeringState, &publishedAt, &access.role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return contentAccess{}, domain.ErrNotFound
		}
		return contentAccess{}, classifyDatabaseError("authorize content access", err)
	}
	if access.role == "" || (access.role == "student" && (access.offeringState == "draft" || !publishedAt.Valid)) {
		return contentAccess{}, domain.ErrNotFound
	}
	return access, nil
}

func (repository ContentRepository) authorizeStaff(ctx context.Context, queries Querier, tenantID, userID, offeringID string) (string, error) {
	access, err := repository.authorizeContent(ctx, queries, tenantID, userID, offeringID)
	if err != nil {
		return "", err
	}
	if access.role != "instructor" && access.role != "lead_instructor" {
		return "", domain.ErrNotFound
	}
	if access.offeringState == "archived" {
		return "", domain.ErrContentConflict
	}
	return access.role, nil
}

func (ContentRepository) ListContent(ctx context.Context, queries Querier, tenantID, userID, offeringID string, now time.Time) (domain.Content, error) {
	access, err := (ContentRepository{}).authorizeContent(ctx, queries, tenantID, userID, offeringID)
	if err != nil {
		return domain.Content{}, err
	}
	content := domain.Content{OfferingID: offeringID, Modules: make([]domain.Module, 0)}
	moduleIndexes := make(map[string]int)
	lessonIndexes := make(map[string]struct{ moduleIndex, lessonIndex int })

	rows, err := queries.Query(ctx, modulesForContentSQL, tenantID, offeringID, access.role, now.UTC())
	if err != nil {
		return domain.Content{}, classifyDatabaseError("list content modules", err)
	}
	for rows.Next() {
		var module domain.Module
		var availableFrom, availableUntil pgtype.Timestamptz
		if err := rows.Scan(&module.ID, &module.OfferingID, &module.Title, &module.Description, &module.Position, &module.Status, &availableFrom, &availableUntil, &module.CreatedBy); err != nil {
			rows.Close()
			return domain.Content{}, classifyDatabaseError("scan content module", err)
		}
		module.AvailableFrom = timestamptzPointer(availableFrom)
		module.AvailableUntil = timestamptzPointer(availableUntil)
		module.Lessons = make([]domain.Lesson, 0)
		moduleIndexes[module.ID] = len(content.Modules)
		content.Modules = append(content.Modules, module)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.Content{}, classifyDatabaseError("iterate content modules", err)
	}
	rows.Close()

	rows, err = queries.Query(ctx, lessonsForContentSQL, tenantID, offeringID, access.role, now.UTC())
	if err != nil {
		return domain.Content{}, classifyDatabaseError("list content lessons", err)
	}
	for rows.Next() {
		var lesson domain.Lesson
		var availableFrom, availableUntil pgtype.Timestamptz
		if err := rows.Scan(&lesson.ID, &lesson.ModuleID, &lesson.Title, &lesson.Description, &lesson.Position, &lesson.LearningMode, &lesson.EstimatedMinutes, &availableFrom, &availableUntil, &lesson.Status, &lesson.CreatedBy); err != nil {
			rows.Close()
			return domain.Content{}, classifyDatabaseError("scan content lesson", err)
		}
		moduleIndex, ok := moduleIndexes[lesson.ModuleID]
		if !ok {
			continue
		}
		lesson.AvailableFrom = timestamptzPointer(availableFrom)
		lesson.AvailableUntil = timestamptzPointer(availableUntil)
		lesson.Materials = make([]domain.Material, 0)
		lessonIndexes[lesson.ID] = struct{ moduleIndex, lessonIndex int }{moduleIndex: moduleIndex, lessonIndex: len(content.Modules[moduleIndex].Lessons)}
		content.Modules[moduleIndex].Lessons = append(content.Modules[moduleIndex].Lessons, lesson)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.Content{}, classifyDatabaseError("iterate content lessons", err)
	}
	rows.Close()

	rows, err = queries.Query(ctx, materialsForContentSQL, tenantID, offeringID, access.role, now.UTC())
	if err != nil {
		return domain.Content{}, classifyDatabaseError("list content materials", err)
	}
	for rows.Next() {
		var material domain.Material
		var materialExternalURL, materialContent pgtype.Text
		var fileID, storageKey, originalFilename, mimeType, checksum, scanStatus, uploadedBy pgtype.Text
		var fileUUID pgtype.UUID
		var fileSize pgtype.Int8
		if err := rows.Scan(&material.ID, &material.LessonID, &material.Title, &material.Description, &material.Type, &fileUUID, &materialExternalURL, &materialContent, &material.Position, &material.Published, &material.CreatedBy, &storageKey, &originalFilename, &mimeType, &fileSize, &checksum, &scanStatus, &uploadedBy); err != nil {
			rows.Close()
			return domain.Content{}, classifyDatabaseError("scan content material", err)
		}
		lessonIndex, ok := lessonIndexes[material.LessonID]
		if !ok {
			continue
		}
		if materialExternalURL.Valid {
			material.ExternalURL = materialExternalURL.String
		}
		if materialContent.Valid {
			material.Content = materialContent.String
		}
		if fileUUID.Valid {
			fileID = pgtype.Text{String: uuid.UUID(fileUUID.Bytes).String(), Valid: true}
			file := &domain.File{ID: fileID.String, OriginalFilename: originalFilename.String, MIMEType: mimeType.String, SizeBytes: fileSize.Int64, Checksum: checksum.String, MalwareScanStatus: scanStatus.String, UploadedBy: uploadedBy.String}
			if access.role != "student" {
				file.StorageKey = storageKey.String
			}
			material.File = file
		}
		content.Modules[lessonIndex.moduleIndex].Lessons[lessonIndex.lessonIndex].Materials = append(content.Modules[lessonIndex.moduleIndex].Lessons[lessonIndex.lessonIndex].Materials, material)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return domain.Content{}, classifyDatabaseError("iterate content materials", err)
	}
	rows.Close()
	return content, nil
}

func (repository ContentRepository) CreateModule(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, _ string, input domain.CreateModuleInput) (domain.Module, error) {
	if _, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID); err != nil {
		return domain.Module{}, err
	}
	id := uuid.New()
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `
INSERT INTO modules (id, tenant_id, course_offering_id, title, description, position, status, available_from, available_until, created_by, created_at, updated_at)
VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, 'draft', $7, $8, $9::uuid, $10, $10)`, id, tenantID, offeringID, input.Title, input.Description, input.Position, input.AvailableFrom, input.AvailableUntil, userID, now); err != nil {
		return domain.Module{}, classifyDatabaseError("create content module", err)
	}
	return domain.Module{ID: id.String(), OfferingID: offeringID, Title: input.Title, Description: input.Description, Position: input.Position, Status: "draft", AvailableFrom: input.AvailableFrom, AvailableUntil: input.AvailableUntil, CreatedBy: userID, Lessons: []domain.Lesson{}}, nil
}

func (repository ContentRepository) UpdateModule(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, moduleID, requestID string, input domain.UpdateModuleInput) (domain.Module, error) {
	role, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID)
	if err != nil {
		return domain.Module{}, err
	}
	id, err := parseContentUUID(moduleID)
	if err != nil {
		return domain.Module{}, err
	}
	var current domain.Module
	var availableFrom, availableUntil pgtype.Timestamptz
	if err := tx.QueryRow(ctx, `SELECT title, description, position, status, available_from, available_until, created_by FROM modules WHERE tenant_id = $1::uuid AND course_offering_id = $2::uuid AND id = $3 FOR UPDATE`, tenantID, offeringID, id).Scan(&current.Title, &current.Description, &current.Position, &current.Status, &availableFrom, &availableUntil, &current.CreatedBy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Module{}, domain.ErrNotFound
		}
		return domain.Module{}, classifyDatabaseError("read content module", err)
	}
	current.ID, current.OfferingID = id.String(), offeringID
	current.AvailableFrom = timestamptzPointer(availableFrom)
	current.AvailableUntil = timestamptzPointer(availableUntil)
	beforeStatus := current.Status
	if input.Title != nil {
		current.Title = *input.Title
	}
	if input.Description != nil {
		current.Description = *input.Description
	}
	if input.Position != nil {
		current.Position = *input.Position
	}
	if input.AvailableFrom != nil {
		current.AvailableFrom = input.AvailableFrom
	}
	if input.AvailableUntil != nil {
		current.AvailableUntil = input.AvailableUntil
	}
	if input.Status != nil {
		if !validContentStatus(*input.Status) {
			return domain.Module{}, domain.ErrInvalidContent
		}
		if role != "lead_instructor" && *input.Status != current.Status {
			return domain.Module{}, domain.ErrContentConflict
		}
		current.Status = *input.Status
	}
	if strings.TrimSpace(current.Title) == "" || current.Position < 0 || domain.ValidateCreateModule(domain.CreateModuleInput{Title: current.Title, Position: current.Position, AvailableFrom: current.AvailableFrom, AvailableUntil: current.AvailableUntil}) != nil {
		return domain.Module{}, domain.ErrInvalidContent
	}
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `UPDATE modules SET title = $1, description = $2, position = $3, status = $4, available_from = $5, available_until = $6, updated_at = $7 WHERE tenant_id = $8::uuid AND course_offering_id = $9::uuid AND id = $10`, current.Title, current.Description, current.Position, current.Status, current.AvailableFrom, current.AvailableUntil, now, tenantID, offeringID, id); err != nil {
		return domain.Module{}, classifyDatabaseError("update content module", err)
	}
	if beforeStatus != current.Status {
		if err := insertContentAudit(ctx, tx, tenantID, userID, role, "content.publication_changed", "module", id, offeringID, requestID, map[string]string{"status": beforeStatus}, map[string]string{"status": current.Status}); err != nil {
			return domain.Module{}, err
		}
	}
	current.Lessons = []domain.Lesson{}
	return current, nil
}

func (repository ContentRepository) CreateLesson(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, moduleID, _ string, input domain.CreateLessonInput) (domain.Lesson, error) {
	if _, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID); err != nil {
		return domain.Lesson{}, err
	}
	moduleUUID, err := parseContentUUID(moduleID)
	if err != nil {
		return domain.Lesson{}, err
	}
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM modules WHERE tenant_id = $1::uuid AND course_offering_id = $2::uuid AND id = $3)`, tenantID, offeringID, moduleUUID).Scan(&exists); err != nil {
		return domain.Lesson{}, classifyDatabaseError("authorize content module", err)
	}
	if !exists {
		return domain.Lesson{}, domain.ErrNotFound
	}
	id := uuid.New()
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `
INSERT INTO lessons (id, tenant_id, module_id, title, description, position, learning_mode, estimated_minutes, available_from, available_until, status, created_by, created_at, updated_at)
VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10, 'draft', $11::uuid, $12, $12)`, id, tenantID, moduleUUID, input.Title, input.Description, input.Position, input.LearningMode, input.EstimatedMinutes, input.AvailableFrom, input.AvailableUntil, userID, now); err != nil {
		return domain.Lesson{}, classifyDatabaseError("create content lesson", err)
	}
	return domain.Lesson{ID: id.String(), ModuleID: moduleID, Title: input.Title, Description: input.Description, Position: input.Position, LearningMode: input.LearningMode, EstimatedMinutes: input.EstimatedMinutes, AvailableFrom: input.AvailableFrom, AvailableUntil: input.AvailableUntil, Status: "draft", CreatedBy: userID, Materials: []domain.Material{}}, nil
}

func (repository ContentRepository) UpdateLesson(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, lessonID, requestID string, input domain.UpdateLessonInput) (domain.Lesson, error) {
	role, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID)
	if err != nil {
		return domain.Lesson{}, err
	}
	id, err := parseContentUUID(lessonID)
	if err != nil {
		return domain.Lesson{}, err
	}
	var current domain.Lesson
	var availableFrom, availableUntil pgtype.Timestamptz
	if err := tx.QueryRow(ctx, `
	SELECT lesson.title, lesson.description, lesson.position, lesson.learning_mode, lesson.estimated_minutes, lesson.available_from, lesson.available_until, lesson.status, lesson.created_by, lesson.module_id
FROM lessons AS lesson
JOIN modules AS module ON module.tenant_id = lesson.tenant_id AND module.id = lesson.module_id
	WHERE lesson.tenant_id = $1::uuid AND module.course_offering_id = $2::uuid AND lesson.id = $3 FOR UPDATE`, tenantID, offeringID, id).Scan(&current.Title, &current.Description, &current.Position, &current.LearningMode, &current.EstimatedMinutes, &availableFrom, &availableUntil, &current.Status, &current.CreatedBy, &current.ModuleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, domain.ErrNotFound
		}
		return domain.Lesson{}, classifyDatabaseError("read content lesson", err)
	}
	current.ID = id.String()
	current.AvailableFrom = timestamptzPointer(availableFrom)
	current.AvailableUntil = timestamptzPointer(availableUntil)
	beforeStatus := current.Status
	if input.Title != nil {
		current.Title = *input.Title
	}
	if input.Description != nil {
		current.Description = *input.Description
	}
	if input.Position != nil {
		current.Position = *input.Position
	}
	if input.LearningMode != nil {
		current.LearningMode = *input.LearningMode
	}
	if input.EstimatedMinutes != nil {
		current.EstimatedMinutes = *input.EstimatedMinutes
	}
	if input.AvailableFrom != nil {
		current.AvailableFrom = input.AvailableFrom
	}
	if input.AvailableUntil != nil {
		current.AvailableUntil = input.AvailableUntil
	}
	if input.Status != nil {
		if !validContentStatus(*input.Status) {
			return domain.Lesson{}, domain.ErrInvalidContent
		}
		if role != "lead_instructor" && *input.Status != current.Status {
			return domain.Lesson{}, domain.ErrContentConflict
		}
		current.Status = *input.Status
	}
	if strings.TrimSpace(current.Title) == "" || current.Position < 0 || current.EstimatedMinutes < 0 || domain.ValidateCreateLesson(domain.CreateLessonInput{Title: current.Title, Position: current.Position, LearningMode: current.LearningMode, EstimatedMinutes: current.EstimatedMinutes, AvailableFrom: current.AvailableFrom, AvailableUntil: current.AvailableUntil}) != nil {
		return domain.Lesson{}, domain.ErrInvalidContent
	}
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `UPDATE lessons SET title = $1, description = $2, position = $3, learning_mode = $4, estimated_minutes = $5, status = $6, available_from = $7, available_until = $8, updated_at = $9 WHERE tenant_id = $10::uuid AND id = $11`, current.Title, current.Description, current.Position, current.LearningMode, current.EstimatedMinutes, current.Status, current.AvailableFrom, current.AvailableUntil, now, tenantID, id); err != nil {
		return domain.Lesson{}, classifyDatabaseError("update content lesson", err)
	}
	if beforeStatus != current.Status {
		if err := insertContentAudit(ctx, tx, tenantID, userID, role, "content.publication_changed", "lesson", id, offeringID, requestID, map[string]string{"status": beforeStatus}, map[string]string{"status": current.Status}); err != nil {
			return domain.Lesson{}, err
		}
	}
	current.Materials = []domain.Material{}
	return current, nil
}

func (repository ContentRepository) CreateMaterial(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, lessonID, _ string, input domain.CreateMaterialInput) (domain.Material, error) {
	if _, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID); err != nil {
		return domain.Material{}, err
	}
	lessonUUID, err := parseContentUUID(lessonID)
	if err != nil {
		return domain.Material{}, err
	}
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM lessons AS lesson JOIN modules AS module ON module.tenant_id = lesson.tenant_id AND module.id = lesson.module_id WHERE lesson.tenant_id = $1::uuid AND module.course_offering_id = $2::uuid AND lesson.id = $3)`, tenantID, offeringID, lessonUUID).Scan(&exists); err != nil {
		return domain.Material{}, classifyDatabaseError("authorize content lesson", err)
	}
	if !exists {
		return domain.Material{}, domain.ErrNotFound
	}
	var fileID any
	if input.FileID != "" {
		fileID, err = parseContentUUID(input.FileID)
		if err != nil {
			return domain.Material{}, err
		}
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM files WHERE tenant_id = $1::uuid AND id = $2)`, tenantID, fileID).Scan(&exists); err != nil {
			return domain.Material{}, classifyDatabaseError("authorize content file", err)
		}
		if !exists {
			return domain.Material{}, domain.ErrNotFound
		}
	}
	id := uuid.New()
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `
INSERT INTO materials (id, tenant_id, lesson_id, title, description, type, file_id, external_url, content, position, published, created_by, created_at, updated_at)
VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, false, $11::uuid, $12, $12)`, id, tenantID, lessonUUID, input.Title, input.Description, input.Type, fileID, input.ExternalURL, input.Content, input.Position, userID, now); err != nil {
		return domain.Material{}, classifyDatabaseError("create content material", err)
	}
	return domain.Material{ID: id.String(), LessonID: lessonID, Title: input.Title, Description: input.Description, Type: input.Type, ExternalURL: input.ExternalURL, Content: input.Content, Position: input.Position, Published: false, CreatedBy: userID}, nil
}

func (repository ContentRepository) UpdateMaterial(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, materialID, requestID string, input domain.UpdateMaterialInput) (domain.Material, error) {
	role, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID)
	if err != nil {
		return domain.Material{}, err
	}
	id, err := parseContentUUID(materialID)
	if err != nil {
		return domain.Material{}, err
	}
	var current domain.Material
	var fileUUID pgtype.UUID
	var externalURL, content pgtype.Text
	if err := tx.QueryRow(ctx, `
SELECT material.lesson_id, material.title, material.description, material.type, material.file_id, material.external_url, material.content, material.position, material.published, material.created_by
FROM materials AS material
JOIN lessons AS lesson ON lesson.tenant_id = material.tenant_id AND lesson.id = material.lesson_id
JOIN modules AS module ON module.tenant_id = lesson.tenant_id AND module.id = lesson.module_id
WHERE material.tenant_id = $1::uuid AND module.course_offering_id = $2::uuid AND material.id = $3 FOR UPDATE`, tenantID, offeringID, id).Scan(&current.LessonID, &current.Title, &current.Description, &current.Type, &fileUUID, &externalURL, &content, &current.Position, &current.Published, &current.CreatedBy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Material{}, domain.ErrNotFound
		}
		return domain.Material{}, classifyDatabaseError("read content material", err)
	}
	current.ID = id.String()
	if externalURL.Valid {
		current.ExternalURL = externalURL.String
	}
	if content.Valid {
		current.Content = content.String
	}
	if input.Title != nil {
		current.Title = *input.Title
	}
	if input.Description != nil {
		current.Description = *input.Description
	}
	if input.ExternalURL != nil {
		current.ExternalURL = *input.ExternalURL
	}
	if input.Content != nil {
		current.Content = *input.Content
	}
	if input.Position != nil {
		current.Position = *input.Position
	}
	beforePublished := current.Published
	if input.Published != nil {
		if role != "lead_instructor" && *input.Published != current.Published {
			return domain.Material{}, domain.ErrContentConflict
		}
		current.Published = *input.Published
	}
	if strings.TrimSpace(current.Title) == "" || current.Position < 0 {
		return domain.Material{}, domain.ErrInvalidContent
	}
	if current.Type == "text" && (current.Content == "" || current.ExternalURL != "") {
		return domain.Material{}, domain.ErrInvalidContent
	}
	if current.Type != "text" && current.Type != "file" && current.ExternalURL == "" {
		return domain.Material{}, domain.ErrInvalidContent
	}
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `UPDATE materials SET title = $1, description = $2, external_url = NULLIF($3, ''), content = NULLIF($4, ''), position = $5, published = $6, updated_at = $7 WHERE tenant_id = $8::uuid AND id = $9`, current.Title, current.Description, current.ExternalURL, current.Content, current.Position, current.Published, now, tenantID, id); err != nil {
		return domain.Material{}, classifyDatabaseError("update content material", err)
	}
	if beforePublished != current.Published {
		if err := insertContentAudit(ctx, tx, tenantID, userID, role, "content.publication_changed", "material", id, offeringID, requestID, map[string]bool{"published": beforePublished}, map[string]bool{"published": current.Published}); err != nil {
			return domain.Material{}, err
		}
	}
	return current, nil
}

func (repository ContentRepository) CreateFile(ctx context.Context, tx pgx.Tx, tenantID, userID, offeringID, _ string, input domain.CreateFileInput) (domain.File, error) {
	if _, err := repository.authorizeStaff(ctx, tx, tenantID, userID, offeringID); err != nil {
		return domain.File{}, err
	}
	id := uuid.New()
	storageKey := "pending/" + id.String()
	now := time.Now().UTC()
	if _, err := tx.Exec(ctx, `
INSERT INTO files (id, tenant_id, storage_key, original_filename, mime_type, size_bytes, checksum, malware_scan_status, uploaded_by, created_at)
VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, 'pending', $8::uuid, $9)`, id, tenantID, storageKey, input.OriginalFilename, input.MIMEType, input.SizeBytes, input.Checksum, userID, now); err != nil {
		return domain.File{}, classifyDatabaseError("create content file metadata", err)
	}
	return domain.File{ID: id.String(), StorageKey: storageKey, OriginalFilename: input.OriginalFilename, MIMEType: input.MIMEType, SizeBytes: input.SizeBytes, Checksum: input.Checksum, MalwareScanStatus: "pending", UploadedBy: userID}, nil
}

func insertContentAudit(ctx context.Context, tx pgx.Tx, tenantID, userID, role, action, entityType string, entityID uuid.UUID, offeringID, requestID string, before, after any) error {
	beforeJSON, err := json.Marshal(before)
	if err != nil {
		return fmt.Errorf("marshal content audit before state: %w", err)
	}
	afterJSON, err := json.Marshal(after)
	if err != nil {
		return fmt.Errorf("marshal content audit after state: %w", err)
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO audit_logs (id, tenant_id, actor_user_id, actor_role, action, entity_type, entity_id, course_offering_id, before_data, after_data, request_id, occurred_at)
VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8::uuid, $9, $10, $11, $12)`, uuid.New(), tenantID, userID, role, action, entityType, entityID, offeringID, beforeJSON, afterJSON, requestID, time.Now().UTC()); err != nil {
		return classifyDatabaseError("write content audit", err)
	}
	return nil
}

func parseContentUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, domain.ErrInvalidContent
	}
	return id, nil
}

func validContentStatus(value string) bool {
	switch value {
	case "draft", "published", "hidden":
		return true
	default:
		return false
	}
}

func timestamptzPointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}
