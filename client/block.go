package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// GetBlockHeight Get height of the current chain.
func (ac *AElfClient) GetBlockHeight() (int64, error) {
	url := ac.Host + BLOCKHEIGHT
	heightBytes, err := utils.GetRequest("GET", url, ac.Version, nil)
	if err != nil {
		return 0, errors.New("Get BlockHeight error:" + err.Error())
	}
	var data interface{}
	json.Unmarshal(heightBytes, &data)
	return int64(data.(float64)), nil
}

// GetBlockByHash Get information of a block by given block hash. Optional whether to include transaction information.
func (ac *AElfClient) GetBlockByHash(blockHash string, includeTransactions bool) (*types.Block, error) {
	_, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, errors.New("transactionID hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{
		"blockHash":           blockHash,
		"includeTransactions": includeTransactions,
	}
	url := ac.Host + BLOCKBYHASH
	blockBytes, err := utils.GetRequest("GET", url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Get Block ByHash error:" + err.Error())
	}
	var block = new(types.Block)
	json.Unmarshal(blockBytes, &block)
	return block, nil
}

// GetBlockByHeight Get information of a block by specified height. Optional whether to include transaction information.
func (ac *AElfClient) GetBlockByHeight(blockHeight int64, includeTransactions bool) (*types.Block, error) {
	params := map[string]interface{}{
		"blockHeight":         blockHeight,
		"includeTransactions": includeTransactions,
	}
	url := ac.Host + BLOCKBYHEIGHT
	blockBytes, err := utils.GetRequest("GET", url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Get Block ByHeight error:" + err.Error())
	}
	var block = new(types.Block)
	json.Unmarshal(blockBytes, &block)
	return block, nil
}
