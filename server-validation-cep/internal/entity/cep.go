package entity

type Response struct {
	CepDetails CepDetails `json:"cep_details"`
	Erro       Erro       `json:"erro"`
}

type Erro struct {
	Message string `json:"message"`
	Success bool   `json:"status"`
}

type CepDetails struct {
	Localidade string  `json:"localidade"`
	TempC      float64 `json:"temp_c"`
	TempF      float64 `json:"temp_f"`
	TempK      float64 `json:"temp_k"`
}
