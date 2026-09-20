package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danilotorchio/fc-go-expert-lab-01/internal/weather"
)

func run() error {
	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		return errors.New("WEATHER_API_KEY is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	mux := http.NewServeMux()
	mux.Handle("GET /weather/{cep}", weather.Handler{
		Cities:       weather.ViaCEP{BaseURL: "https://viacep.com.br", Client: client},
		Temperatures: weather.WeatherAPI{BaseURL: "https://api.weatherapi.com", Key: key, Client: client},
	})

	slog.Info("listening", "port", port)
	return http.ListenAndServe(":"+port, mux)
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
