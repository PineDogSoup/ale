package test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestGetChainStatus(t *testing.T) {
	chainStatus, err := mainClient.AElf.GetChainStatus()

	assert.NoError(t, err)
	assert.Equal(t, "AELF", chainStatus.ChainId)
	assert.NotEmpty(t, chainStatus.Branches)
	assert.True(t, chainStatus.LongestChainHeight > 0)
	assert.NotEmpty(t, chainStatus.LongestChainHash)
	assert.NotEmpty(t, chainStatus.GenesisContractAddress)
	assert.NotEmpty(t, chainStatus.GenesisBlockHash)
	assert.True(t, chainStatus.LastIrreversibleBlockHeight > 0)
	assert.NotEmpty(t, chainStatus.LastIrreversibleBlockHash)
	assert.True(t, chainStatus.BestChainHeight > 0)
	assert.NotEmpty(t, chainStatus.BestChainHash)

	longestChainBlock, err := mainClient.AElf.GetBlockByHash(chainStatus.LongestChainHash, false)
	assert.Equal(t, longestChainBlock.Header.Height, chainStatus.LongestChainHeight)

	genesisBlock, err := mainClient.AElf.GetBlockByHash(chainStatus.GenesisBlockHash, false)
	assert.Equal(t, int64(1), genesisBlock.Header.Height)

	lastIrreversibleBlock, err := mainClient.AElf.GetBlockByHash(chainStatus.LastIrreversibleBlockHash, false)
	assert.Equal(t, lastIrreversibleBlock.Header.Height, chainStatus.LastIrreversibleBlockHeight)

	bestChainBlock, err := mainClient.AElf.GetBlockByHash(chainStatus.BestChainHash, false)
	assert.Equal(t, bestChainBlock.Header.Height, chainStatus.BestChainHeight)

	genesisContractAddress, err := mainClient.AElf.GetGenesisContractAddress()
	assert.Equal(t, genesisContractAddress, chainStatus.GenesisContractAddress)
}

func TestGetChainID(t *testing.T) {
	chainID, err := mainClient.AElf.GetChainID()
	assert.NoError(t, err)
	assert.Equal(t, 9992731, chainID)
}

func TestGetTaskQueueStatus(t *testing.T) {
	taskQueueStatus, err := mainClient.AElf.GetTaskQueueStatus()
	assert.NoError(t, err)
	spew.Dump("Get Task Queue Status Result", taskQueueStatus)
}

func TestGetContractFileDescriptorSet(t *testing.T) {
	contractAddr, err := mainClient.AElf.GetGenesisContractAddress()
	assert.NoError(t, err)
	contractFile, err := mainClient.AElf.GetContractFileDescriptorSet(contractAddr)
	assert.NoError(t, err)
	assert.NotEmpty(t, contractFile)
	spew.Dump("Get contract File Descriptor Set Result", contractFile)
}
