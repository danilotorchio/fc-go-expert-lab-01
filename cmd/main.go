package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danilotorchio/fc-go-expert-lab-01/internal/config"
	"github.com/danilotorchio/fc-go-expert-lab-01/internal/transport"
	"github.com/danilotorchio/fc-go-expert-lab-01/internal/weather"
)

const (
	viaCepURL     = "https://viacep.com.br"
	weatherApiURL = "https://api.weatherapi.com"
)

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	handler := weather.NewHandler(&weather.NewHandlerOpts{
		Cities:       weather.ViaCEP{BaseURL: viaCepURL, Client: client},
		Temperatures: weather.WeatherAPI{BaseURL: weatherApiURL, Key: cfg.WeatherApiKey, Client: client},
	})

	router := transport.NewRouter()

	server := transport.NewServer(&transport.NewServerOpts{
		Port:    cfg.Port,
		Handler: router.Handler(handler),
	})

	if err := server.Start(ctx, 15*time.Second); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
