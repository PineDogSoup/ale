package types

type Transaction struct {
	From           string
	To             string
	RefBlockNumber int64
	RefBlockPrefix string
	MethodName     string
	Params         string
	Signature      string
}

type TransactionResult struct {
	TransactionId string
	Status        string
	Logs          []LogEvent
	Bloom         string
	BlockNumber   int64
	BlockHash     string
	Transaction   Transaction
	ReturnValue   string
	Error         string
}
