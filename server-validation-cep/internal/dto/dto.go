package dto

import "go.opentelemetry.io/otel/trace"

type CepRequest struct {
	Cep string `json:"cep"`
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type TemplateData struct {
	Title           string
	ExternalCallURL string
	OTELTracer      trace.Tracer
}
