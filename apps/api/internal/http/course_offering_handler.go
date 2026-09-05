package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	nethttp "net/http"
	"strings"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/farisakbar28/campus-lms/apps/api/internal/middleware"
)

type rosterReader interface {
	AuthorizedRoster(context.Context, string, string, string) (domain.Roster, error)
}

func courseOfferingParticipants(reader rosterReader, logger *slog.Logger) nethttp.HandlerFunc {
	return func(response nethttp.ResponseWriter, request *nethttp.Request) {
		offeringID := request.PathValue("id")
		if !isUUID(offeringID) {
			writeError(response, nethttp.StatusBadRequest, "invalid_course_offering_id", "course offering ID must be a UUID")
			return
		}

		principal, ok := middleware.PrincipalFromContext(request.Context())
		if !ok || !isUUID(principal.TenantID) || !isUUID(principal.UserID) {
			writeError(response, nethttp.StatusUnauthorized, "unauthenticated", "authentication is required")
			return
		}

		roster, err := reader.AuthorizedRoster(request.Context(), principal.TenantID, principal.UserID, offeringID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrNotFound):
				writeError(response, nethttp.StatusNotFound, "course_offering_not_found", "course offering was not found")
			case errors.Is(err, database.ErrUnavailable):
				logger.Error("read course offering participants", "error", err)
				writeError(response, nethttp.StatusServiceUnavailable, "service_unavailable", "service is temporarily unavailable")
			default:
				logger.Error("read course offering participants", "error", err)
				writeError(response, nethttp.StatusInternalServerError, "internal_error", "internal server error")
			}
			return
		}

		writeJSON(response, nethttp.StatusOK, rosterResponseFromDomain(roster))
	}
}

type rosterResponse struct {
	CourseOffering courseOfferingResponse `json:"course_offering"`
	Participants   []participantResponse  `json:"participants"`
}

type courseOfferingResponse struct {
	ID          string               `json:"id"`
	ExternalID  string               `json:"external_id"`
	SectionCode string               `json:"section_code"`
	DisplayName string               `json:"display_name"`
	LMSStatus   string               `json:"lms_status"`
	Course      courseResponse       `json:"course"`
	Term        academicTermResponse `json:"academic_term"`
}

type courseResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type academicTermResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type participantResponse struct {
	StudentUserID    string `json:"student_user_id"`
	DisplayName      string `json:"display_name"`
	EnrollmentStatus string `json:"enrollment_status"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func rosterResponseFromDomain(roster domain.Roster) rosterResponse {
	participants := make([]participantResponse, 0, len(roster.Participants))
	for _, participant := range roster.Participants {
		participants = append(participants, participantResponse{
			StudentUserID:    participant.StudentUserID,
			DisplayName:      participant.DisplayName,
			EnrollmentStatus: participant.EnrollmentStatus,
		})
	}

	return rosterResponse{
		CourseOffering: courseOfferingResponse{
			ID:          roster.Offering.ID,
			ExternalID:  roster.Offering.ExternalID,
			SectionCode: roster.Offering.SectionCode,
			DisplayName: roster.Offering.DisplayName,
			LMSStatus:   roster.Offering.LMSStatus,
			Course: courseResponse{
				ID:   roster.Offering.Course.ID,
				Code: roster.Offering.Course.Code,
				Name: roster.Offering.Course.Name,
			},
			Term: academicTermResponse{
				ID:   roster.Offering.Term.ID,
				Code: roster.Offering.Term.Code,
				Name: roster.Offering.Term.Name,
			},
		},
		Participants: participants,
	}
}

func writeError(response nethttp.ResponseWriter, status int, code, message string) {
	writeJSON(response, status, errorResponse{Code: code, Message: message})
}

func writeJSON(response nethttp.ResponseWriter, status int, body any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(body); err != nil {
		return
	}
}

func isUUID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for index, character := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			if character != '-' {
				return false
			}
			continue
		}
		if !strings.ContainsRune("0123456789abcdefABCDEF", character) {
			return false
		}
	}

	return true
}
