package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/infra/service"
	"github.com/stretchr/testify/assert"
)

type mockCepDetailsService struct {
	response *dto.Response
	err      error
	cep      string
	url      string
}

func (m *mockCepDetailsService) GetCepDetails(_ context.Context, cep, url string) (*dto.Response, error) {
	m.cep = cep
	m.url = url
	return m.response, m.err
}

func newCepRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/cep", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), "MyCoreHost", "http://mock")
	return req.WithContext(ctx)
}

func TestGetCep_MissingCep(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := newCepRequest(`{}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_InvalidRequestBody(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := newCepRequest(`{`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid request body", resp["error"])
}

func TestGetCep_InvalidCepLength(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := newCepRequest(`{"cep":"123"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_InvalidCepCharacters(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{},
	}
	req := newCepRequest(`{"cep":"abcdefgh"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_ZipcodeNotFound(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{err: service.ErrZipcodeNotFound},
	}
	req := newCepRequest(`{"cep":"01001000"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "can not find zipcode", resp["error"])
}

func TestGetCep_ServiceError(t *testing.T) {
	handler := &CepHandler{
		Service: &mockCepDetailsService{err: assert.AnError},
	}
	req := newCepRequest(`{"cep":"01001000"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "internal server error", resp["error"])
}

func TestGetCep_Success(t *testing.T) {
	cepDetailsResponse := &dto.Response{City: "Sao paulo", TempC: 25.0, TempF: 77.0, TempK: 298.0}
	mockService := &mockCepDetailsService{response: cepDetailsResponse}
	handler := &CepHandler{
		Service: mockService,
	}
	req := newCepRequest(`{"cep":"01001000"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "01001000", mockService.cep)
	assert.Equal(t, "http://mock", mockService.url)
	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 25.0, resp.TempC)
	assert.Equal(t, 77.0, resp.TempF)
	assert.Equal(t, 298.0, resp.TempK)
}
