package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/infra/service"
)

type CepHandler struct {
	Service service.CepDetailsInterface
}

func NewCepHandler(service service.CepDetailsInterface) *CepHandler {
	return &CepHandler{
		Service: service,
	}
}

// @Summary Service A - Buscar informações do CEP
// @Description Retorna informações do CEP consultando a API do ViaCEP e da WeatherAPI.
// @Tags CEP
// @Accept json
// @Produce json
// @Param request body dto.CepResponse true "CEP sem formatacao (ex: 01001000)"
// @Router /cep [post]
func (h *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	myHost := r.Context().Value("MyCoreHost").(string)

	var request dto.CepResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	cep := request.Cep
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "CEP obrigatorio"})
		return
	}

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

	result, err := h.Service.GetCepDetails(cep, myHost)
	if err != nil {
		if err == service.ErrZipcodeNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "can not find zipcode"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
