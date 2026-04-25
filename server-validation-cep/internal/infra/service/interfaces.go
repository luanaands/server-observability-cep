package service

import (
	"github.com/luanaands/server-validation-cep/internal/dto"
)

type CepDetailsInterface interface {
	GetCepDetails(cep string, url string) (*dto.Response, error)
}
