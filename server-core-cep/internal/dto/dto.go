package dto

import (
	"errors"

	"go.opentelemetry.io/otel/trace"
)

type CepRequest struct {
	Cep string `json:"cep"`
}

type CepResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type TemplateData struct {
	Title              string
	BackgroundColor    string
	ExternalCallURL    string
	ExternalCallMethod string
	RequestNameOTEL    string
	OTELTracer         trace.Tracer
}

var ErrZipcodeNotFound = errors.New("zipcode not found")
