//go:generate easyjson -all model.go

package model

type TransactionResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	NextAction    string `json:"nextAction"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount"`
}
