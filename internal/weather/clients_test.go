package weather

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ws/01001000/json/":
			w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`))
		case "/ws/99999999/json/":
			w.Write([]byte(`{"erro":"true"}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()
	v := ViaCEP{BaseURL: srv.URL, Client: srv.Client()}

	city, err := v.City(t.Context(), "01001000")
	if err != nil || city != "São Paulo" {
		t.Errorf("City(01001000) = %q, %v; want São Paulo, nil", city, err)
	}
	if _, err := v.City(t.Context(), "99999999"); !errors.Is(err, ErrZipcodeNotFound) {
		t.Errorf("City(99999999) err = %v, want ErrZipcodeNotFound", err)
	}
	if _, err := v.City(t.Context(), "bad"); err == nil || errors.Is(err, ErrZipcodeNotFound) {
		t.Errorf("City(bad) err = %v, want upstream error", err)
	}
}

func TestWeatherAPI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != "/v1/current.json" || q.Get("key") != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if q.Get("q") != "São Paulo, Brazil" {
			t.Errorf("q = %q, want São Paulo, Brazil", q.Get("q"))
		}
		w.Write([]byte(`{"location":{"name":"Sao Paulo"},"current":{"temp_c":28.5}}`))
	}))
	defer srv.Close()

	w := WeatherAPI{BaseURL: srv.URL, Key: "secret", Client: srv.Client()}
	c, err := w.CelsiusIn(t.Context(), "São Paulo")
	if err != nil || c != 28.5 {
		t.Errorf("CelsiusIn = %v, %v; want 28.5, nil", c, err)
	}

	w.Key = "wrong"
	if _, err := w.CelsiusIn(t.Context(), "São Paulo"); err == nil {
		t.Error("CelsiusIn with bad key: want error")
	}
}
