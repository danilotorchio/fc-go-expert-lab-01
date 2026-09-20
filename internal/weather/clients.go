package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ErrZipcodeNotFound is returned when ViaCEP has no record for the CEP.
var ErrZipcodeNotFound = errors.New("can not find zipcode")

// ViaCEP looks up the city of a CEP.
type ViaCEP struct {
	BaseURL string // e.g. https://viacep.com.br
	Client  *http.Client
}

func (v ViaCEP) City(ctx context.Context, cep string) (string, error) {
	var body struct {
		Localidade string `json:"localidade"`
	}
	if err := getJSON(ctx, v.Client, v.BaseURL+"/ws/"+cep+"/json/", &body); err != nil {
		return "", fmt.Errorf("viacep: %w", err)
	}
	// Unknown CEPs come back as 200 {"erro": "true"}, so an empty city means not found.
	if body.Localidade == "" {
		return "", ErrZipcodeNotFound
	}
	return body.Localidade, nil
}

// WeatherAPI fetches the current temperature of a city from weatherapi.com.
type WeatherAPI struct {
	BaseURL string // e.g. https://api.weatherapi.com
	Key     string
	Client  *http.Client
}

func (w WeatherAPI) CelsiusIn(ctx context.Context, city string) (float64, error) {
	q := url.Values{"key": {w.Key}, "q": {city + ", Brazil"}, "aqi": {"no"}}
	var body struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := getJSON(ctx, w.Client, w.BaseURL+"/v1/current.json?"+q.Encode(), &body); err != nil {
		return 0, fmt.Errorf("weatherapi: %w", err)
	}
	return body.Current.TempC, nil
}

func getJSON(ctx context.Context, c *http.Client, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}
