package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCepDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"cep_details":{"localidade":"Sao Paulo","temp_c":25.5,"temp_f":77.9,"temp_k":298.6},"erro":{"message":"","status":true}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	s := &CepDetailsService{}
	resp, err := s.GetCepDetails("01001000", server.URL)
	require.NoError(t, err)
	assert.Equal(t, "Sao Paulo", resp.Localidade)
	assert.Equal(t, 25.5, resp.TempC)
	assert.Equal(t, 77.9, resp.TempF)
	assert.Equal(t, 298.6, resp.TempK)
}
