package dto

type CepResponse struct {
	Cep string `json:"cep"`
}

type Response struct {
	Localidade string  `json:"localidade"`
	TempC      float64 `json:"temp_c"`
	TempF      float64 `json:"temp_f"`
	TempK      float64 `json:"temp_k"`
}
