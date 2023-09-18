package types

type GetHolderInfoInput struct {
	CaHash string `json:"caHash"`
}

type GetBalanceInput struct {
	Owner  string
	Symbol string
}
