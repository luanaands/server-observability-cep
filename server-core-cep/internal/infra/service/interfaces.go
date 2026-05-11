package service

import (
	"context"

	"github.com/luanaands/server-core-cep/internal/dto"
)

type CepInterface interface {
	GetViaCep(ctx context.Context, cep string, url string) (*dto.CepResponse, error)
}

type WeatherInterface interface {
	GetWeather(ctx context.Context, city string, apiKey string, baseURL string) (*dto.WeatherResponse, error)
}
