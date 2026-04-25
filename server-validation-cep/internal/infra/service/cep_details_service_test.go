package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCepDetails(t *testing.T) {
	s := &CepDetailsService{}
	resp, err := s.GetCepDetails("01001000", "https://viacep.com.br/ws/01001000/json/")
	assert.Nil(t, err)
	assert.Equal(t, resp.Localidade, "São Paulo")
}
