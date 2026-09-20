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

type Handler struct {
	cities       CityFinder
	temperatures TemperatureFinder
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	if !cepPattern.MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.cities.City(r.Context(), cep)
	if errors.Is(err, ErrZipcodeNotFound) {
		http.Error(w, ErrZipcodeNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("city lookup failed", "cep", cep, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	c, err := h.temperatures.CelsiusIn(r.Context(), city)
	if err != nil {
		slog.Error("temperature lookup failed", "city", city, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewTemperature(c))
}

type NewHandlerOpts struct {
	Cities       CityFinder
	Temperatures TemperatureFinder
}

func NewHandler(opts *NewHandlerOpts) *Handler {
	return &Handler{
		cities:       opts.Cities,
		temperatures: opts.Temperatures,
	}
}
