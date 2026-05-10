package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/entity"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var ErrZipcodeNotFound = errors.New("zipcode not found")

type CepDetailsService struct {
	client *http.Client
}

func NewCepDetailsService() *CepDetailsService {
	return &CepDetailsService{
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (s *CepDetailsService) GetCepDetails(ctx context.Context, cep string, url string) (*dto.Response, error) {
	body := map[string]string{"cep": cep}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("Service B error: status %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrZipcodeNotFound
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *entity.CepDetails
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	dtoResponse := &dto.Response{
		City:  response.City,
		TempC: response.TempC,
		TempF: response.TempF,
		TempK: response.TempK,
	}
	return dtoResponse, nil
}
