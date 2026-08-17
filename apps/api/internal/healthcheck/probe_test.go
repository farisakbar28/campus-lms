package healthcheck

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{name: "ok", status: http.StatusOK},
		{name: "service unavailable", status: http.StatusServiceUnavailable, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if request.Method != http.MethodGet {
					t.Errorf("method = %s, want %s", request.Method, http.MethodGet)
				}
				if request.URL.Path != "/healthz" {
					t.Errorf("path = %s, want /healthz", request.URL.Path)
				}
				response.WriteHeader(test.status)
			}))
			defer server.Close()

			err := Probe(context.Background(), server.URL+"/healthz")
			if (err != nil) != test.wantErr {
				t.Fatalf("Probe() error = %v, wantErr %t", err, test.wantErr)
			}
		})
	}
}
