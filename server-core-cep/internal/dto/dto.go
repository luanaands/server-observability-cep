package dto

import "errors"

type CepRequest struct {
	Cep string `json:"cep"`
}

type CepResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

var ErrZipcodeNotFound = errors.New("zipcode not found")
