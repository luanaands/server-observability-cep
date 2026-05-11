package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/luanaands/server-core-cep/internal/dto"
	zipcode "github.com/luanaands/server-core-cep/internal/zipecode"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Mock implementations
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

func newTestConfig() *dto.TemplateData {
	return &dto.TemplateData{
		Title:              "Service B",
		BackgroundColor:    "#FF5733",
		ExternalCallURL:    "http://external",
		ExternalCallMethod: "POST",
		RequestNameOTEL:    "GetCep",
		OTELTracer:         &mockTracer{},
	}
}

func newCepRequest(cep string) *http.Request {
	body := dto.CepRequest{Cep: cep}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/cep", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(context.WithValue(context.WithValue(req.Context(), "ViaCepHost", "http://mock"), "ApiWeatherHost", "http://mock"), "ApiWeatherKey", "key")
	return req.WithContext(ctx)
}

func newCepRequestWithBody(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/cep", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(context.WithValue(context.WithValue(req.Context(), "ViaCepHost", "http://mock"), "ApiWeatherHost", "http://mock"), "ApiWeatherKey", "key")
	return req.WithContext(ctx)
}

type mockCepService struct {
	response *dto.CepResponse
	err      error
}

func (m *mockCepService) GetViaCep(ctx context.Context, cep, url string) (*dto.CepResponse, error) {
	return m.response, m.err
}

type mockWeatherService struct {
	response *dto.WeatherResponse
	err      error
}

func (m *mockWeatherService) GetWeather(ctx context.Context, city, key, host string) (*dto.WeatherResponse, error) {
	return m.response, m.err
}

func TestGetCep_MissingCep(t *testing.T) {
	handler := &CepHandler{
		Service:        &mockCepService{},
		WeatherService: &mockWeatherService{},
		Config:         newTestConfig(),
	}
	req := newCepRequestWithBody(`{}`)
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_InvalidCepLength(t *testing.T) {
	handler := &CepHandler{
		Service:        &mockCepService{},
		WeatherService: &mockWeatherService{},
		Config:         newTestConfig(),
	}
	req := newCepRequest("123")
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_InvalidCepCharacters(t *testing.T) {
	handler := &CepHandler{
		Service:        &mockCepService{},
		WeatherService: &mockWeatherService{},
		Config:         newTestConfig(),
	}
	req := newCepRequest("abcdefgh")
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid zipcode", resp["error"])
}

func TestGetCep_GetViaCepError(t *testing.T) {
	handler := &CepHandler{
		Service:        &mockCepService{err: zipcode.ErrZipcodeNotFound},
		WeatherService: &mockWeatherService{},
		Config:         newTestConfig(),
	}
	req := newCepRequest("01001000")
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "can not find zipcode", resp["error"])
}

func TestGetCep_GetWeatherError(t *testing.T) {
	viaCepResp := &dto.CepResponse{Localidade: "São Paulo"}
	handler := &CepHandler{
		Service:        &mockCepService{response: viaCepResp},
		WeatherService: &mockWeatherService{err: assert.AnError},
		Config:         newTestConfig(),
	}
	req := newCepRequest("01001000")
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "internal server error", resp["error"])
}

func TestGetCep_Success(t *testing.T) {
	viaCepResp := &dto.CepResponse{Localidade: "São Paulo"}
	weatherResp := &dto.WeatherResponse{
		TempC: 25.0,
		TempF: 77.0,
	}
	handler := &CepHandler{
		Service:        &mockCepService{response: viaCepResp},
		WeatherService: &mockWeatherService{response: weatherResp},
		Config:         newTestConfig(),
	}
	req := newCepRequest("01001000")
	w := httptest.NewRecorder()

	handler.GetCep(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "São Paulo", resp.City)
	assert.Equal(t, 25.0, resp.TempC)
	assert.Equal(t, 77.0, resp.TempF)
	assert.Equal(t, 298.0, resp.TempK)
}
