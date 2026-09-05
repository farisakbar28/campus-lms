package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/farisakbar28/campus-lms/apps/api/internal/database"
	"github.com/farisakbar28/campus-lms/apps/api/internal/domain"
	"github.com/farisakbar28/campus-lms/apps/api/internal/middleware"
)

const validOfferingID = "2888c021-06ae-73da-2f57-884c1dd5d059"
const validTenantID = "19cd4773-2aeb-d614-028f-e21bf9b73d0c"
const validUserID = "4ec42919-bc1a-17bc-10b0-d75b8343dff8"

type countingRosterStub struct {
	calls  int
	roster domain.Roster
	err    error
}

func (stub *countingRosterStub) AuthorizedRoster(context.Context, string, string, string) (domain.Roster, error) {
	stub.calls++
	return stub.roster, stub.err
}

func TestCourseOfferingParticipants(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		principal  *middleware.Principal
		stub       countingRosterStub
		wantStatus int
		wantCalls  int
		wantCode   string
	}{
		{
			name:       "rejects malformed offering ID",
			path:       "/course-offerings/not-a-uuid/participants",
			wantStatus: nethttp.StatusBadRequest,
			wantCalls:  0,
			wantCode:   "invalid_course_offering_id",
		},
		{
			name:       "fails closed without principal",
			path:       "/course-offerings/" + validOfferingID + "/participants",
			wantStatus: nethttp.StatusUnauthorized,
			wantCalls:  0,
			wantCode:   "unauthenticated",
		},
		{
			name:      "maps not found without leaking access state",
			path:      "/course-offerings/" + validOfferingID + "/participants",
			principal: &middleware.Principal{TenantID: validTenantID, UserID: validUserID},
			stub: countingRosterStub{
				err: domain.ErrNotFound,
			},
			wantStatus: nethttp.StatusNotFound,
			wantCalls:  1,
			wantCode:   "course_offering_not_found",
		},
		{
			name:      "maps database unavailable",
			path:      "/course-offerings/" + validOfferingID + "/participants",
			principal: &middleware.Principal{TenantID: validTenantID, UserID: validUserID},
			stub: countingRosterStub{
				err: errors.Join(database.ErrUnavailable, errors.New("connection reset")),
			},
			wantStatus: nethttp.StatusServiceUnavailable,
			wantCalls:  1,
			wantCode:   "service_unavailable",
		},
		{
			name:      "maps unexpected errors safely",
			path:      "/course-offerings/" + validOfferingID + "/participants",
			principal: &middleware.Principal{TenantID: validTenantID, UserID: validUserID},
			stub: countingRosterStub{
				err: errors.New("SQLSTATE secret database detail"),
			},
			wantStatus: nethttp.StatusInternalServerError,
			wantCalls:  1,
			wantCode:   "internal_error",
		},
		{
			name:      "returns roster response",
			path:      "/course-offerings/" + validOfferingID + "/participants",
			principal: &middleware.Principal{TenantID: validTenantID, UserID: validUserID},
			stub: countingRosterStub{
				roster: domain.Roster{
					Offering: domain.CourseOffering{ID: validOfferingID, DisplayName: "Algorithms"},
					Participants: []domain.Participant{{
						StudentUserID:    "72d328d4-7352-bc28-f2cf-2ae3306dbfcf",
						DisplayName:      "Seed Student 1",
						EnrollmentStatus: "active",
					}},
				},
			},
			wantStatus: nethttp.StatusOK,
			wantCalls:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := courseOfferingParticipants(&test.stub, slog.New(slog.NewTextHandler(io.Discard, nil)))
			request := httptest.NewRequest(nethttp.MethodGet, test.path, nil)
			request.SetPathValue("id", strings.TrimSuffix(strings.TrimPrefix(test.path, "/course-offerings/"), "/participants"))
			if test.principal != nil {
				request = request.WithContext(middleware.WithPrincipal(request.Context(), *test.principal))
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if test.stub.calls != test.wantCalls {
				t.Fatalf("repository calls = %d, want %d", test.stub.calls, test.wantCalls)
			}

			if test.wantCode != "" {
				var body errorResponse
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if body.Code != test.wantCode {
					t.Errorf("error code = %q, want %q", body.Code, test.wantCode)
				}
				if body.Message == "SQLSTATE secret database detail" {
					t.Error("response leaked internal database detail")
				}
			}
		})
	}
}
