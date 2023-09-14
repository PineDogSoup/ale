package types

type SendResult struct {
	TransactionId string
}

type CallResult struct {
	Data    interface{}
	Message string
}
