package handlers

import (
	"encoding/json"
	"net/http"

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

// @Summary Buscar informações do CEP
// @Description Retorna informações do CEP consultando a API do ViaCEP e da WeatherAPI.
// @Tags CEP
// @Accept json
// @Produce json
// @Param cep query string true "CEP sem formatação (ex: 01001000)"
// @Router /cep [post]
func (h *CepHandler) GetCep(w http.ResponseWriter, r *http.Request) {
	myHost := r.Context().Value("MyCoreHost").(string)
	cep := r.URL.Query().Get("cep")

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "CEP é obrigatório"})
		return
	}

	if len(cep) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid zipcode"})
		return
	}

	result, err := h.Service.GetCepDetails(cep, myHost)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
