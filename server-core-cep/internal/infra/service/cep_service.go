package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/luanaands/server-core-cep/internal/dto"
	"github.com/luanaands/server-core-cep/internal/entity"
	zipcode "github.com/luanaands/server-core-cep/internal/zipecode"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type CepService struct {
	client *http.Client
}

func NewCepService() *CepService {
	return &CepService{
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (s *CepService) GetViaCep(cep string, url string) (*dto.CepResponse, error) {
	req, err := http.NewRequest("GET", url+"/"+cep+"/json", nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("viacep service error: status %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, zipcode.ErrZipcodeNotFound
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *entity.CepViaCepResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Erro == "true" {
		return nil, zipcode.ErrZipcodeNotFound
	}

	var dtoResponse *dto.CepResponse
	dtoResponse = dto.FromViaCep(response)
	return dtoResponse, nil
}
