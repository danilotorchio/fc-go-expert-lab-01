package weather

import "math"

type Temperature struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewTemperature(c float64) Temperature {
	return Temperature{
		Celsius:    round(c),
		Fahrenheit: round(c*1.8 + 32),
		Kelvin:     round(c + 273),
	}
}

func round(v float64) float64 {
	return math.Round(v*100) / 100
}
