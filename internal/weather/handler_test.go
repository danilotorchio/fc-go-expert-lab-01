package weather

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeCities map[string]string

func (f fakeCities) City(_ context.Context, cep string) (string, error) {
	if cep == "00000000" {
		return "", errors.New("upstream down")
	}
	if city, ok := f[cep]; ok {
		return city, nil
	}
	return "", ErrZipcodeNotFound
}

type fakeTemps float64

func (f fakeTemps) CelsiusIn(context.Context, string) (float64, error) { return float64(f), nil }

func TestHandler(t *testing.T) {
	mux := http.NewServeMux()

	mux.Handle("GET /weather/{cep}", NewHandler(&NewHandlerOpts{
		Cities:       fakeCities{"01001000": "São Paulo"},
		Temperatures: fakeTemps(28.5),
	}))

	tests := []struct {
		name, cep string
		status    int
		body      string
	}{
		{"success", "01001000", http.StatusOK, `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`},
		{"too short", "0100100", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"too long", "010010000", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"letters", "0100100a", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"with dash", "01001-000", http.StatusUnprocessableEntity, "invalid zipcode"},
		{"not found", "99999999", http.StatusNotFound, "can not find zipcode"},
		{"upstream error", "00000000", http.StatusInternalServerError, "internal server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/weather/"+tt.cep, nil))
			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tt.body {
				t.Errorf("body = %q, want %q", got, tt.body)
			}
		})
	}
}
