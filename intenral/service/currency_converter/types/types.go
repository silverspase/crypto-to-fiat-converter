package types

type ConvertRequest struct {
	FromToken  string
	ToCurrency string
	Amount     float32
}

type ConvertResponse struct {
	FromToken   string
	ToCurrency  string
	TokenAmount float32
	TotalPrice  float32
}
