package weather

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
)

var cepPattern = regexp.MustCompile(`^\d{8}$`)

type CityFinder interface {
	City(ctx context.Context, cep string) (string, error)
}

type TemperatureFinder interface {
	CelsiusIn(ctx context.Context, city string) (float64, error)
}

// Handler serves GET /weather/{cep}.
type Handler struct {
	Cities       CityFinder
	Temperatures TemperatureFinder
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	if !cepPattern.MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.Cities.City(r.Context(), cep)
	if errors.Is(err, ErrZipcodeNotFound) {
		http.Error(w, ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("city lookup failed", "cep", cep, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	c, err := h.Temperatures.CelsiusIn(r.Context(), city)
	if err != nil {
		slog.Error("temperature lookup failed", "city", city, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewTemperature(c))
}
