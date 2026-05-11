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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type mockCepDetailsService struct {
	response *dto.Response
	err      error
	cep      string
	url      string
}

type mockSpan struct {
	trace.Span
}

func (ms *mockSpan) End(options ...trace.SpanEndOption)                  {}
func (ms *mockSpan) AddEvent(name string, options ...trace.EventOption)  {}
func (ms *mockSpan) IsRecording() bool                                   { return false }
func (ms *mockSpan) RecordError(err error, options ...trace.EventOption) {}
func (ms *mockSpan) SetAttributes(kv ...attribute.KeyValue)              {}
func (ms *mockSpan) SetName(name string)                                 {}
func (ms *mockSpan) SetStatus(code codes.Code, description string)       {}
func (ms *mockSpan) TracerProvider() trace.TracerProvider                { return nil }
func (ms *mockSpan) AddLink(link trace.Link)                             {}
func (ms *mockSpan) SpanContext() trace.SpanContext                      { return trace.SpanContext{} }

type mockTracer struct {
	trace.Tracer
}

func (mt *mockTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, &mockSpan{}
}

func (m *mockCepDetailsService) GetCepDetails(_ context.Context, cep, url string) (*dto.Response, error) {
	m.cep = cep
	m.url = url
	return m.response, m.err
}

func newTestConfig() *dto.TemplateData {
	return &dto.TemplateData{
		Title:           "Service A",
		ExternalCallURL: "http://external",
		OTELTracer:      &mockTracer{},
	}
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
		Config:  newTestConfig(),
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
		Config:  newTestConfig(),
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
		Config:  newTestConfig(),
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
		Config:  newTestConfig(),
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
		Config:  newTestConfig(),
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
		Config:  newTestConfig(),
	}
	req := newCepRequest(`{"cep":"01001000"}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, assert.AnError.Error(), resp["error"])
}

func TestGetCep_Success(t *testing.T) {
	cepDetailsResponse := &dto.Response{City: "Sao paulo", TempC: 25.0, TempF: 77.0, TempK: 298.0}
	mockService := &mockCepDetailsService{response: cepDetailsResponse}
	handler := &CepHandler{
		Service: mockService,
		Config:  newTestConfig(),
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
