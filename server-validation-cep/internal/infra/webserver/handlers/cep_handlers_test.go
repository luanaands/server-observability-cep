package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/stretchr/testify/assert"
)

// Mock implementations
type mockCepDetailsService struct {
	response *dto.Response
	err      error
}

func (m *mockCepDetailsService) GetCepDetails(cep, url string) (*dto.Response, error) {
	return m.response, m.err
}

func TestGetCep_MissingCep(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := httptest.NewRequest("GET", "/cep", nil)
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "CEP é obrigatório", resp["error"])
}

func TestGetCep_InvalidCepLength(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := httptest.NewRequest("GET", "/cep?cep=123", nil)
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_GetViaCepError(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{err: assert.AnError},
	}
	req := httptest.NewRequest("GET", "/cep?cep=01001000", nil)
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "can not find zipcode", resp["error"])
}

func TestGetCep_GetWeatherError(t *testing.T) {
	viaCepResp := &dto.Response{Localidade: "Sao paulo", TempC: 25.0, TempF: 77.0, TempK: 298.0}
	handler := &CepHandler{
		Service: &mockCepDetailsService{response: viaCepResp},
	}
	req := httptest.NewRequest("GET", "/cep?cep=01001000", nil)
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "can not find weather", resp["error"])
}

func TestGetCep_Success(t *testing.T) {
	cepDetailsResponse := &dto.Response{Localidade: "Sao paulo", TempC: 25.0, TempF: 77.0, TempK: 298.0}
	handler := &CepHandler{
		Service: &mockCepDetailsService{response: cepDetailsResponse},
	}
	req := httptest.NewRequest("GET", "/cep?cep=01001000", nil)
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 25.0, resp.TempC)
	assert.Equal(t, 77.0, resp.TempF)
	assert.Equal(t, 298.0, resp.TempK)
}
