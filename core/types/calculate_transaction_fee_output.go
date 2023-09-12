package types

type CalculateTransactionFee struct {
	Success        bool
	TransactionFee map[string]interface{}
	ResourceFee    map[string]interface{}
}
