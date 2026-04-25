package dto

import "github.com/luanaands/server-validation-cep/internal/entity"

func FromViaCep(resp *entity.CepViaCepResponse) *CepResponse {
	return &CepResponse{
		Localidade: resp.Localidade,
	}
}

func FromWeather(resp *entity.WeatherResponse) *WeatherResponse {
	return &WeatherResponse{
		TempC: resp.Current.TempC,
		TempF: resp.Current.TempF,
	}
}
