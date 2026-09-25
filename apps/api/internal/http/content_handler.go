package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/farisakbar28/campus-lms/apps/api/internal/middleware"
	"github.com/google/uuid"
)

type contentService interface {
	ListContent(context.Context, string, string, string, time.Time) (domain.Content, error)
	CreateModule(context.Context, string, string, string, string, domain.CreateModuleInput) (domain.Module, error)
	UpdateModule(context.Context, string, string, string, string, string, domain.UpdateModuleInput) (domain.Module, error)
	CreateLesson(context.Context, string, string, string, string, string, domain.CreateLessonInput) (domain.Lesson, error)
	UpdateLesson(context.Context, string, string, string, string, string, domain.UpdateLessonInput) (domain.Lesson, error)
	CreateMaterial(context.Context, string, string, string, string, string, domain.CreateMaterialInput) (domain.Material, error)
	UpdateMaterial(context.Context, string, string, string, string, string, domain.UpdateMaterialInput) (domain.Material, error)
	CreateFile(context.Context, string, string, string, string, domain.CreateFileInput) (domain.File, error)
}

type contentHandler struct {
	service contentService
	logger  *slog.Logger
}

type createModuleRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Position       int        `json:"position"`
	AvailableFrom  *time.Time `json:"available_from"`
	AvailableUntil *time.Time `json:"available_until"`
}

type updateModuleRequest struct {
	Title          *string      `json:"title"`
	Description    *string      `json:"description"`
	Position       *int         `json:"position"`
	Status         *string      `json:"status"`
	AvailableFrom  optionalTime `json:"available_from"`
	AvailableUntil optionalTime `json:"available_until"`
}

type createLessonRequest struct {
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Position         int        `json:"position"`
	LearningMode     string     `json:"learning_mode"`
	EstimatedMinutes int        `json:"estimated_minutes"`
	AvailableFrom    *time.Time `json:"available_from"`
	AvailableUntil   *time.Time `json:"available_until"`
}

type updateLessonRequest struct {
	Title            *string      `json:"title"`
	Description      *string      `json:"description"`
	Position         *int         `json:"position"`
	LearningMode     *string      `json:"learning_mode"`
	EstimatedMinutes *int         `json:"estimated_minutes"`
	Status           *string      `json:"status"`
	AvailableFrom    optionalTime `json:"available_from"`
	AvailableUntil   optionalTime `json:"available_until"`
}

type optionalTime struct {
	set   bool
	value *time.Time
}

func (value *optionalTime) UnmarshalJSON(data []byte) error {
	value.set = true
	if string(data) == "null" {
		value.value = nil
		return nil
	}
	var parsed time.Time
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	value.value = &parsed
	return nil
}

func (value optionalTime) domainValue() domain.OptionalTime {
	return domain.OptionalTime{Set: value.set, Value: value.value}
}

type createMaterialRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	FileID      string `json:"file_id"`
	ExternalURL string `json:"external_url"`
	Content     string `json:"content"`
	Position    int    `json:"position"`
}

type updateMaterialRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	ExternalURL *string `json:"external_url"`
	Content     *string `json:"content"`
	Position    *int    `json:"position"`
	Published   *bool   `json:"published"`
}

type createFileRequest struct {
	OriginalFilename string `json:"original_filename"`
	MIMEType         string `json:"mime_type"`
	SizeBytes        int64  `json:"size_bytes"`
	Checksum         string `json:"checksum"`
}

func (handler contentHandler) list(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, ok := handler.principalAndOffering(response, request)
	if !ok {
		return
	}
	content, err := handler.service.ListContent(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, time.Now().UTC())
	if err != nil {
		handler.writeError(response, "list content", err)
		return
	}
	writeJSON(response, nethttp.StatusOK, contentResponseFromDomain(content))
}

func (handler contentHandler) createModule(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, ok := handler.principalAndOffering(response, request)
	if !ok {
		return
	}
	var body createModuleRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	module, err := handler.service.CreateModule(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, requestID, domain.CreateModuleInput{Title: body.Title, Description: body.Description, Position: body.Position, AvailableFrom: body.AvailableFrom, AvailableUntil: body.AvailableUntil})
	if err != nil {
		handler.writeError(response, "create content module", err)
		return
	}
	writeJSON(response, nethttp.StatusCreated, moduleResponseFromDomain(module))
}

func (handler contentHandler) updateModule(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, moduleID, ok := handler.principalOfferingAndID(response, request, "module_id")
	if !ok {
		return
	}
	var body updateModuleRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	module, err := handler.service.UpdateModule(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, moduleID, requestID, domain.UpdateModuleInput{Title: body.Title, Description: body.Description, Position: body.Position, Status: body.Status, AvailableFrom: body.AvailableFrom.domainValue(), AvailableUntil: body.AvailableUntil.domainValue()})
	if err != nil {
		handler.writeError(response, "update content module", err)
		return
	}
	writeJSON(response, nethttp.StatusOK, moduleResponseFromDomain(module))
}

func (handler contentHandler) createLesson(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, moduleID, ok := handler.principalOfferingAndID(response, request, "module_id")
	if !ok {
		return
	}
	var body createLessonRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	lesson, err := handler.service.CreateLesson(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, moduleID, requestID, domain.CreateLessonInput{Title: body.Title, Description: body.Description, Position: body.Position, LearningMode: body.LearningMode, EstimatedMinutes: body.EstimatedMinutes, AvailableFrom: body.AvailableFrom, AvailableUntil: body.AvailableUntil})
	if err != nil {
		handler.writeError(response, "create content lesson", err)
		return
	}
	writeJSON(response, nethttp.StatusCreated, lessonResponseFromDomain(lesson))
}

func (handler contentHandler) updateLesson(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, lessonID, ok := handler.principalOfferingAndID(response, request, "lesson_id")
	if !ok {
		return
	}
	var body updateLessonRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	lesson, err := handler.service.UpdateLesson(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, lessonID, requestID, domain.UpdateLessonInput{Title: body.Title, Description: body.Description, Position: body.Position, LearningMode: body.LearningMode, EstimatedMinutes: body.EstimatedMinutes, Status: body.Status, AvailableFrom: body.AvailableFrom.domainValue(), AvailableUntil: body.AvailableUntil.domainValue()})
	if err != nil {
		handler.writeError(response, "update content lesson", err)
		return
	}
	writeJSON(response, nethttp.StatusOK, lessonResponseFromDomain(lesson))
}

func (handler contentHandler) createMaterial(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, lessonID, ok := handler.principalOfferingAndID(response, request, "lesson_id")
	if !ok {
		return
	}
	var body createMaterialRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	material, err := handler.service.CreateMaterial(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, lessonID, requestID, domain.CreateMaterialInput{Title: body.Title, Description: body.Description, Type: body.Type, FileID: body.FileID, ExternalURL: body.ExternalURL, Content: body.Content, Position: body.Position})
	if err != nil {
		handler.writeError(response, "create content material", err)
		return
	}
	writeJSON(response, nethttp.StatusCreated, materialResponseFromDomain(material))
}

func (handler contentHandler) updateMaterial(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, materialID, ok := handler.principalOfferingAndID(response, request, "material_id")
	if !ok {
		return
	}
	var body updateMaterialRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	material, err := handler.service.UpdateMaterial(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, materialID, requestID, domain.UpdateMaterialInput{Title: body.Title, Description: body.Description, ExternalURL: body.ExternalURL, Content: body.Content, Position: body.Position, Published: body.Published})
	if err != nil {
		handler.writeError(response, "update content material", err)
		return
	}
	writeJSON(response, nethttp.StatusOK, materialResponseFromDomain(material))
}

func (handler contentHandler) createFile(response nethttp.ResponseWriter, request *nethttp.Request) {
	principal, offeringID, ok := handler.principalAndOffering(response, request)
	if !ok {
		return
	}
	var body createFileRequest
	if !decodeContentJSON(response, request, &body) {
		return
	}
	requestID := setContentRequestID(response)
	file, err := handler.service.CreateFile(request.Context(), principal.TenantID.String(), principal.UserID.String(), offeringID, requestID, domain.CreateFileInput{OriginalFilename: body.OriginalFilename, MIMEType: body.MIMEType, SizeBytes: body.SizeBytes, Checksum: body.Checksum})
	if err != nil {
		handler.writeError(response, "create content file metadata", err)
		return
	}
	writeJSON(response, nethttp.StatusCreated, fileResponseFromDomain(file))
}

func (handler contentHandler) principalAndOffering(response nethttp.ResponseWriter, request *nethttp.Request) (middleware.Principal, string, bool) {
	principal, ok := principalFromContentRequest(response, request)
	if !ok {
		return middleware.Principal{}, "", false
	}
	offeringID := request.PathValue("offering_id")
	if !isUUID(offeringID) {
		writeError(response, nethttp.StatusBadRequest, "invalid_course_offering_id", "course offering ID must be a UUID")
		return middleware.Principal{}, "", false
	}
	return principal, offeringID, true
}

func (handler contentHandler) principalOfferingAndID(response nethttp.ResponseWriter, request *nethttp.Request, key string) (middleware.Principal, string, string, bool) {
	principal, offeringID, ok := handler.principalAndOffering(response, request)
	if !ok {
		return middleware.Principal{}, "", "", false
	}
	entityID := request.PathValue(key)
	if !isUUID(entityID) {
		writeError(response, nethttp.StatusBadRequest, "invalid_content_id", "content ID must be a UUID")
		return middleware.Principal{}, "", "", false
	}
	return principal, offeringID, entityID, true
}

func (handler contentHandler) writeError(response nethttp.ResponseWriter, operation string, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(response, nethttp.StatusNotFound, "content_not_found", "content was not found")
	case errors.Is(err, domain.ErrInvalidContent):
		writeError(response, nethttp.StatusBadRequest, "invalid_content", "content request is invalid")
	case errors.Is(err, domain.ErrContentConflict):
		writeError(response, nethttp.StatusConflict, "content_conflict", "content cannot be changed in its current state")
	case errors.Is(err, database.ErrUnavailable):
		handler.logger.Error(operation, "error", err)
		writeError(response, nethttp.StatusServiceUnavailable, "service_unavailable", "service is temporarily unavailable")
	default:
		handler.logger.Error(operation, "error", err)
		writeError(response, nethttp.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func principalFromContentRequest(response nethttp.ResponseWriter, request *nethttp.Request) (middleware.Principal, bool) {
	principal, ok := middleware.PrincipalFromContext(request.Context())
	if !ok || principal.TenantID == uuid.Nil || principal.UserID == uuid.Nil || principal.SessionID == uuid.Nil || principal.MembershipID == uuid.Nil {
		writeError(response, nethttp.StatusUnauthorized, "unauthenticated", "authentication is required")
		return middleware.Principal{}, false
	}
	return principal, true
}

func decodeContentJSON(response nethttp.ResponseWriter, request *nethttp.Request, target any) bool {
	request.Body = nethttp.MaxBytesReader(response, request.Body, 1<<20)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(response, nethttp.StatusBadRequest, "invalid_json", "request body is invalid")
		return false
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(response, nethttp.StatusBadRequest, "invalid_json", "request body must contain one JSON object")
		return false
	}
	return true
}

func setContentRequestID(response nethttp.ResponseWriter) string {
	requestID := uuid.NewString()
	response.Header().Set("X-Request-ID", requestID)
	return requestID
}

type contentResponse struct {
	OfferingID string           `json:"offering_id"`
	Modules    []moduleResponse `json:"modules"`
}

type moduleResponse struct {
	ID             string           `json:"id"`
	OfferingID     string           `json:"offering_id"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Position       int              `json:"position"`
	Status         string           `json:"status"`
	AvailableFrom  *time.Time       `json:"available_from,omitempty"`
	AvailableUntil *time.Time       `json:"available_until,omitempty"`
	CreatedBy      string           `json:"created_by"`
	Lessons        []lessonResponse `json:"lessons"`
}

type lessonResponse struct {
	ID               string             `json:"id"`
	ModuleID         string             `json:"module_id"`
	Title            string             `json:"title"`
	Description      string             `json:"description"`
	Position         int                `json:"position"`
	LearningMode     string             `json:"learning_mode"`
	EstimatedMinutes int                `json:"estimated_minutes"`
	AvailableFrom    *time.Time         `json:"available_from,omitempty"`
	AvailableUntil   *time.Time         `json:"available_until,omitempty"`
	Status           string             `json:"status"`
	CreatedBy        string             `json:"created_by"`
	Materials        []materialResponse `json:"materials"`
}

type materialResponse struct {
	ID          string        `json:"id"`
	LessonID    string        `json:"lesson_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	File        *fileResponse `json:"file,omitempty"`
	ExternalURL string        `json:"external_url,omitempty"`
	Content     string        `json:"content,omitempty"`
	Position    int           `json:"position"`
	Published   bool          `json:"published"`
	CreatedBy   string        `json:"created_by"`
}

type fileResponse struct {
	ID                string `json:"id"`
	StorageKey        string `json:"storage_key,omitempty"`
	OriginalFilename  string `json:"original_filename"`
	MIMEType          string `json:"mime_type"`
	SizeBytes         int64  `json:"size_bytes"`
	Checksum          string `json:"checksum"`
	MalwareScanStatus string `json:"malware_scan_status"`
	UploadedBy        string `json:"uploaded_by"`
}

func contentResponseFromDomain(content domain.Content) contentResponse {
	result := contentResponse{OfferingID: content.OfferingID, Modules: make([]moduleResponse, 0, len(content.Modules))}
	for _, module := range content.Modules {
		result.Modules = append(result.Modules, moduleResponseFromDomain(module))
	}
	return result
}

func moduleResponseFromDomain(module domain.Module) moduleResponse {
	result := moduleResponse{ID: module.ID, OfferingID: module.OfferingID, Title: module.Title, Description: module.Description, Position: module.Position, Status: module.Status, AvailableFrom: module.AvailableFrom, AvailableUntil: module.AvailableUntil, CreatedBy: module.CreatedBy, Lessons: make([]lessonResponse, 0, len(module.Lessons))}
	for _, lesson := range module.Lessons {
		result.Lessons = append(result.Lessons, lessonResponseFromDomain(lesson))
	}
	return result
}

func lessonResponseFromDomain(lesson domain.Lesson) lessonResponse {
	result := lessonResponse{ID: lesson.ID, ModuleID: lesson.ModuleID, Title: lesson.Title, Description: lesson.Description, Position: lesson.Position, LearningMode: lesson.LearningMode, EstimatedMinutes: lesson.EstimatedMinutes, AvailableFrom: lesson.AvailableFrom, AvailableUntil: lesson.AvailableUntil, Status: lesson.Status, CreatedBy: lesson.CreatedBy, Materials: make([]materialResponse, 0, len(lesson.Materials))}
	for _, material := range lesson.Materials {
		result.Materials = append(result.Materials, materialResponseFromDomain(material))
	}
	return result
}

func materialResponseFromDomain(material domain.Material) materialResponse {
	result := materialResponse{ID: material.ID, LessonID: material.LessonID, Title: material.Title, Description: material.Description, Type: material.Type, ExternalURL: material.ExternalURL, Content: material.Content, Position: material.Position, Published: material.Published, CreatedBy: material.CreatedBy}
	if material.File != nil {
		result.File = &fileResponse{ID: material.File.ID, StorageKey: material.File.StorageKey, OriginalFilename: material.File.OriginalFilename, MIMEType: material.File.MIMEType, SizeBytes: material.File.SizeBytes, Checksum: material.File.Checksum, MalwareScanStatus: material.File.MalwareScanStatus, UploadedBy: material.File.UploadedBy}
	}
	return result
}

func fileResponseFromDomain(file domain.File) fileResponse {
	return fileResponse{ID: file.ID, StorageKey: file.StorageKey, OriginalFilename: file.OriginalFilename, MIMEType: file.MIMEType, SizeBytes: file.SizeBytes, Checksum: file.Checksum, MalwareScanStatus: file.MalwareScanStatus, UploadedBy: file.UploadedBy}
}

func registerContentRoutes(mux *nethttp.ServeMux, protect func(nethttp.Handler) nethttp.Handler, service contentService, logger *slog.Logger) {
	handler := contentHandler{service: service, logger: logger}
	mux.Handle("GET /tenants/{tenant_id}/course-offerings/{offering_id}/content", protect(nethttp.HandlerFunc(handler.list)))
	mux.Handle("POST /tenants/{tenant_id}/course-offerings/{offering_id}/modules", protect(nethttp.HandlerFunc(handler.createModule)))
	mux.Handle("PATCH /tenants/{tenant_id}/course-offerings/{offering_id}/modules/{module_id}", protect(nethttp.HandlerFunc(handler.updateModule)))
	mux.Handle("POST /tenants/{tenant_id}/course-offerings/{offering_id}/modules/{module_id}/lessons", protect(nethttp.HandlerFunc(handler.createLesson)))
	mux.Handle("PATCH /tenants/{tenant_id}/course-offerings/{offering_id}/lessons/{lesson_id}", protect(nethttp.HandlerFunc(handler.updateLesson)))
	mux.Handle("POST /tenants/{tenant_id}/course-offerings/{offering_id}/lessons/{lesson_id}/materials", protect(nethttp.HandlerFunc(handler.createMaterial)))
	mux.Handle("PATCH /tenants/{tenant_id}/course-offerings/{offering_id}/materials/{material_id}", protect(nethttp.HandlerFunc(handler.updateMaterial)))
	mux.Handle("POST /tenants/{tenant_id}/course-offerings/{offering_id}/files", protect(nethttp.HandlerFunc(handler.createFile)))
}
