package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/luanaands/server-validation-cep/internal/dto"
	"github.com/luanaands/server-validation-cep/internal/entity"
)

type CepDetailsService struct {
	client *http.Client
}

func NewCepDetailsService() *CepDetailsService {
	return &CepDetailsService{
		client: &http.Client{},
	}
}

func (s *CepDetailsService) GetCepDetails(cep string, url string) (*dto.Response, error) {
	body := map[string]string{"cep": cep}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response *entity.Response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	if response.Erro.Success != true {
		return nil, errors.New(response.Erro.Message)
	}

	dtoResponse := &dto.Response{}
	dtoResponse.Localidade = response.CepDetails.Localidade
	dtoResponse.TempC = response.CepDetails.TempC
	dtoResponse.TempF = response.CepDetails.TempF
	dtoResponse.TempK = response.CepDetails.TempK

	return dtoResponse, nil
}
