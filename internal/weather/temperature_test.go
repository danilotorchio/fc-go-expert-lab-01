package weather

import "testing"

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		c    float64
		want Temperature
	}{
		{28.5, Temperature{28.5, 83.3, 301.5}},
		{0, Temperature{0, 32, 273}},
		{100, Temperature{100, 212, 373}},
		{-40, Temperature{-40, -40, 233}},
		{21.37, Temperature{21.37, 70.47, 294.37}},
	}
	for _, tt := range tests {
		if got := NewTemperature(tt.c); got != tt.want {
			t.Errorf("NewTemperature(%v) = %+v, want %+v", tt.c, got, tt.want)
		}
	}
}
