package types

type MerklePath struct {
	MerklePathNodes []MerklePathNode
}

type MerklePathNode struct {
	Hash            string
	IsLeftChildNode bool
}
