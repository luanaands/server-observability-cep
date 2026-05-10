package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/luanaands/server-core-cep/internal/dto"
	"github.com/luanaands/server-core-cep/internal/infra/service"
	zipcode "github.com/luanaands/server-core-cep/internal/zipecode"
	"go.opentelemetry.io/otel/attribute"
)

type CepHandler struct {
	Service        service.CepInterface
	WeatherService service.WeatherInterface
	Config         *dto.TemplateData
}

func NewCepHandler(service service.CepInterface, weatherService service.WeatherInterface, config *dto.TemplateData) *CepHandler {
	return &CepHandler{
		Service:        service,
		WeatherService: weatherService,
		Config:         config,
	}
}

// @Summary Service B - Buscar localidade e clima atual
// @Description Retorna dados do tempo consultando ViaCEP e WeatherAPI
// @Tags CEP
// @Accept json
// @Produce json
// @Param cep body dto.CepRequest true "CEP sem formatação (ex: 01001000)"
// @Router /cep [Post]
func (h *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	viaCepUrl := r.Context().Value("ViaCepHost").(string)
	apiWeatherHost := r.Context().Value("ApiWeatherHost").(string)
	apiWeatherKey := r.Context().Value("ApiWeatherKey").(string)

	ctx := r.Context()
	spanName := strings.TrimSpace(h.Config.Title)
	if spanName == "" {
		spanName = "Service B"
	}
	ctx, span := h.Config.OTELTracer.Start(ctx, spanName+" - "+h.Config.RequestNameOTEL)
	span.SetAttributes(
		attribute.String("service.title", h.Config.Title),
		attribute.String("service.background_color", h.Config.BackgroundColor),
		attribute.String("service.external_call_url", h.Config.ExternalCallURL),
		attribute.String("service.external_call_method", h.Config.ExternalCallMethod),
	)
	defer span.End()

	var req dto.CepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	cep := req.Cep
	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid zipcode"})
		return
	}

	if _, err := strconv.Atoi(cep); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid zipcode"})
		return
	}

	viaCepResponse, err := h.Service.GetViaCep(cep, viaCepUrl)
	if err != nil {
		if errors.Is(err, zipcode.ErrZipcodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "can not find zipcode"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	realtimeWeather, err := h.WeatherService.GetWeather(viaCepResponse.Localidade, apiWeatherKey, apiWeatherHost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	var result = &dto.Response{
		City:  viaCepResponse.Localidade,
		TempC: realtimeWeather.TempC,
		TempF: realtimeWeather.TempF,
		TempK: realtimeWeather.TempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
