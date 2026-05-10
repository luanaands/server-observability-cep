package service

import (
	"context"

	"github.com/luanaands/server-validation-cep/internal/dto"
)

type CepDetailsInterface interface {
	GetCepDetails(ctx context.Context, cep string, url string) (*dto.Response, error)
}
