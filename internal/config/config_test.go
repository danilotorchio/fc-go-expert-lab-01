package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Run("defaults port to 8080", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "secret")
		t.Setenv("PORT", "")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() err = %v", err)
		}
		if cfg.Port != "8080" || cfg.WeatherApiKey != "secret" {
			t.Errorf("LoadConfig() = %+v, want port 8080 and key secret", cfg)
		}
	})

	t.Run("reads port from env", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "secret")
		t.Setenv("PORT", "9090")

		cfg, err := LoadConfig()
		if err != nil || cfg.Port != "9090" {
			t.Errorf("LoadConfig() = %+v, %v; want port 9090", cfg, err)
		}
	})

	t.Run("requires weather api key", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "")

		if _, err := LoadConfig(); err == nil {
			t.Error("LoadConfig() without WEATHER_API_KEY: want error")
		}
	})
}
