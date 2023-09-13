package types

type TransferInput struct {
	To     string `json:"to"`
	Symbol string `json:"symbol"`
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}
