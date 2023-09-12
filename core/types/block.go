package types

import "time"

type Block struct {
	BlockHash string
	Header    BlockHeader
	Body      BlockBody
}

type BlockBody struct {
	TransactionsCount int
	Transactions      []string
}

type BlockHeader struct {
	PreviousBlockHash                string
	MerkleTreeRootOfTransactions     string
	MerkleTreeRootOfWorldState       string
	MerkleTreeRootOfTransactionState string
	Extra                            string
	Height                           int64
	Time                             time.Time
	ChainId                          string
	Bloom                            string
	SignerPubkey                     string
}
