package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/infra/service"
	"go.opentelemetry.io/otel/attribute"
)

type CepHandler struct {
	Service service.CepDetailsInterface
	Config  *dto.TemplateData
}

func NewCepHandler(service service.CepDetailsInterface, config *dto.TemplateData) *CepHandler {
	return &CepHandler{
		Service: service,
		Config:  config,
	}
}

// @Summary Service A - Buscar informações do CEP
// @Description Retorna informações do CEP
// @Tags CEP
// @Accept json
// @Produce json
// @Param request body dto.CepRequest true "CEP sem formatacao (ex: 01001000)"
// @Router /cep [post]
func (h *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	myHost := r.Context().Value("MyCoreHost").(string)
	ctx := r.Context()
	ctx, span := h.Config.OTELTracer.Start(ctx, "Service-a.request "+"/cep")
	span.SetAttributes(
		attribute.String("service.title", h.Config.Title),
		attribute.String("service.external_call_url", h.Config.ExternalCallURL),
	)
	defer span.End()

	var request dto.CepRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	cep := request.Cep
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

	result, err := h.Service.GetCepDetails(ctx, cep, myHost)
	if err != nil {
		if err == service.ErrZipcodeNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "can not find zipcode"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
