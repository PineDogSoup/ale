package client

import (
	"ale/core/types"
	"ale/utils"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"

	"github.com/btcsuite/btcutil/base58"
)

// const const.
const (
	CHAINSTATUS             = "/api/blockChain/chainStatus"
	BLOCKHEIGHT             = "/api/blockChain/blockHeight"
	BLOCKBYHASH             = "/api/blockChain/block"
	BLOCKBYHEIGHT           = "/api/blockChain/blockByHeight"
	TRANSACTIONPOOLSTATUS   = "/api/blockChain/transactionPoolStatus"
	RAWTRANSACTION          = "/api/blockChain/rawTransaction"
	SENDTRANSACTION         = "/api/blockChain/sendTransaction"
	SENDRAWTRANSACTION      = "/api/blockChain/sendRawTransaction"
	TASKQUEUESTATUS         = "/api/blockChain/taskQueueStatus"
	TRANSACTIONRESULT       = "/api/blockChain/transactionResult"
	TRANSACTIONRESULTS      = "/api/blockChain/transactionResults"
	MBYTRANSACTIONID        = "/api/blockChain/merklePathByTransactionId"
	ADDPEER                 = "/api/net/peer"
	REMOVEPEER              = "/api/net/peer"
	PEERS                   = "/api/net/peers"
	NETWORKINFO             = "/api/net/networkInfo"
	SENDTRANSACTIONS        = "/api/blockChain/sendTransactions"
	EXECUTETRANSACTION      = "/api/blockChain/executeTransaction"
	EXECUTERAWTRANSACTION   = "/api/blockChain/executeRawTransaction"
	FILEDESCRIPTOR          = "/api/blockChain/contractFileDescriptorSet"
	CALCULATETRANSACTIONFEE = "/api/blockChain/calculateTransactionFee"

	privateKeyForView = "680afd630d82ae5c97942c4141d60b8a9fedfa5b2864fca84072c17ee1f72d9d"
	addressForView    = "SD6BXDrKT2syNd1WehtPyRo3dPBiXqfGUj8UJym7YP9W9RynM"
)

type AElfAPI interface {
	GetChainStatus() (*types.ChainStatus, error)
	GetContractFileDescriptorSet(address string) ([]byte, error)
	GetChainID() (int, error)
	GetTaskQueueStatus() ([]*types.TaskQueueInfo, error)
	GetNetworkInfo() (*types.NetworkInfo, error)
	GetPeers(withMetrics bool) ([]*types.Peer, error)
	GetBlockHeight() (int64, error)
	GetBlockByHash(blockHash string, includeTransactions bool) (*types.Block, error)
	GetBlockByHeight(blockHeight int64, includeTransactions bool) (*types.Block, error)

	CalculateTransactionFee(rawTransaction string) (*types.CalculateTransactionFee, error)

	CreateRawTransaction(input *types.CreateRawTransactionInput) (*types.CreateRawTransactionOutput, error)

	SendTransaction(transaction string) (*types.SendTransactionOutput, error)
	SendTransactions(rawTransactions string) ([]interface{}, error)
	SendRawTransaction(transaction, signature string, returnTransaction bool) (*types.SendRawTransaction, error)

	ExecuteTransaction(rawTransaction string) (string, error)
	ExecuteRawTransaction(input *types.ExecuteRawTransaction) (string, error)
}

// GetChainStatus Get the current status of the block chain.
func (c *AElfClient) GetChainStatus() (*types.ChainStatus, error) {
	url := c.Host + CHAINSTATUS
	chainBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get ChainStatus error:" + err.Error())
	}
	var chain = new(types.ChainStatus)
	json.Unmarshal(chainBytes, &chain)
	return chain, nil
}

// GetContractFileDescriptorSet Get the definitions of proto-buff related to a contract.
func (c *AElfClient) GetContractFileDescriptorSet(address string) ([]byte, error) {
	url := c.Host + FILEDESCRIPTOR
	params := map[string]interface{}{"address": address}
	data, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get ContractFile Descriptor Set error:" + err.Error())
	}
	return data, err
}

// GetChainID Get id of the chain.
func (c *AElfClient) GetChainID() (int, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return 0, errors.New("Get Chain Status error:" + err.Error())
	}
	chainIDBytes := base58.Decode(chainStatus.ChainId)
	if len(chainIDBytes) < 4 {
		var bs [4]byte
		for i := 0; i < 4; i++ {
			bs[i] = 0
			if len(chainIDBytes) > i {
				bs[i] = chainIDBytes[i]
			}
		}
		chainIDBytes = bs[:]
	}
	return utils.BytesToInt(chainIDBytes), nil
}

// GetTaskQueueStatus Get the status information of the task queue.
func (c *AElfClient) GetTaskQueueStatus() ([]*types.TaskQueueInfo, error) {
	url := c.Host + TASKQUEUESTATUS
	queues, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get Task Queue Status error:" + err.Error())
	}
	var datas interface{}
	json.Unmarshal(queues, &datas)
	var queueInfos []*types.TaskQueueInfo
	for _, data := range datas.([]interface{}) {
		var queue = new(types.TaskQueueInfo)
		queueBytes, err := json.Marshal(data)
		if err != nil {
			return nil, errors.New("json Marshal data error:" + err.Error())
		}
		json.Unmarshal(queueBytes, &queue)
		queueInfos = append(queueInfos, queue)
	}
	return queueInfos, nil
}

// GetNetworkInfo Get the node's network information.
func (c *AElfClient) GetNetworkInfo() (*types.NetworkInfo, error) {
	url := c.Host + NETWORKINFO
	networkBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get Network Info error:" + err.Error())
	}
	var network = new(types.NetworkInfo)
	json.Unmarshal(networkBytes, &network)
	return network, nil
}

// GetPeers Gets information about the peer nodes of the current node.Optional whether to include metrics.
func (c *AElfClient) GetPeers(withMetrics bool) ([]*types.Peer, error) {
	url := c.Host + PEERS
	params := map[string]interface{}{"withMetrics": withMetrics}
	peerBytes, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get Peers error:" + err.Error())
	}
	var datas interface{}
	var peers []*types.Peer
	json.Unmarshal(peerBytes, &datas)
	for _, data := range datas.([]interface{}) {
		var peer = new(types.Peer)
		peerBytes, _ := json.Marshal(data)
		json.Unmarshal(peerBytes, &peer)
		peers = append(peers, peer)
	}
	return peers, nil
}

// GetBlockHeight Get height of the current chain.
func (c *AElfClient) GetBlockHeight() (int64, error) {
	url := c.Host + BLOCKHEIGHT
	heightBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return 0, errors.New("Get BlockHeight error:" + err.Error())
	}
	var data interface{}
	json.Unmarshal(heightBytes, &data)
	return int64(data.(float64)), nil
}

// GetBlockByHash Get information of a block by given block hash. Optional whether to include transaction information.
func (c *AElfClient) GetBlockByHash(blockHash string, includeTransactions bool) (*types.Block, error) {
	_, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, errors.New("transactionID hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{
		"blockHash":           blockHash,
		"includeTransactions": includeTransactions,
	}
	url := c.Host + BLOCKBYHASH
	blockBytes, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get Block ByHash error:" + err.Error())
	}
	var block = new(types.Block)
	json.Unmarshal(blockBytes, &block)
	return block, nil
}

// GetBlockByHeight Get information of a block by specified height. Optional whether to include transaction information.
func (c *AElfClient) GetBlockByHeight(blockHeight int64, includeTransactions bool) (*types.Block, error) {
	params := map[string]interface{}{
		"blockHeight":         blockHeight,
		"includeTransactions": includeTransactions,
	}
	url := c.Host + BLOCKBYHEIGHT
	blockBytes, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get Block ByHeight error:" + err.Error())
	}
	var block = new(types.Block)
	json.Unmarshal(blockBytes, &block)
	return block, nil
}

// ExecuteTransaction  Call a read-only method of a contract.
func (c *AElfClient) ExecuteTransaction(rawTransaction string) (string, error) {
	url := c.Host + EXECUTETRANSACTION
	params := map[string]interface{}{"RawTransaction": rawTransaction}
	transactionBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return "", errors.New("Execute Transaction error:" + err.Error())
	}
	return utils.BytesToString(transactionBytes), nil
}

// ExecuteRawTransaction Call a method of a contract by given serialized strings.
func (c *AElfClient) ExecuteRawTransaction(input *types.ExecuteRawTransaction) (string, error) {
	url := c.Host + EXECUTERAWTRANSACTION
	params := map[string]interface{}{
		"RawTransaction": input.RawTransaction,
		"Signature":      input.Signature,
	}
	transactionBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return "", errors.New("Execute RawTransaction error:" + err.Error())
	}
	//var data interface{}
	//json.Unmarshal(transactionBytes, &data)
	return utils.BytesToString(transactionBytes), nil
}

// SendTransaction Broadcast a transaction.
func (c *AElfClient) SendTransaction(transaction string) (*types.SendTransactionOutput, error) {
	url := c.Host + SENDTRANSACTION
	params := map[string]interface{}{"RawTransaction": transaction}
	transactionBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return nil, errors.New("Send Transaction error:" + err.Error())
	}
	var output = new(types.SendTransactionOutput)
	json.Unmarshal(transactionBytes, &output)
	return output, nil
}

// CreateRawTransaction Creates an unsigned serialized transaction.
func (c *AElfClient) CreateRawTransaction(input *types.CreateRawTransactionInput) (*types.CreateRawTransactionOutput, error) {
	url := c.Host + RAWTRANSACTION
	params := map[string]interface{}{
		"From":           input.From,
		"MethodName":     input.MethodName,
		"Params":         input.Params,
		"RefBlockHash":   input.RefBlockHash,
		"RefBlockNumber": input.RefBlockNumber,
		"To":             input.To,
	}
	transactionBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return nil, errors.New("Create RawTransaction error:" + err.Error())
	}
	var output = new(types.CreateRawTransactionOutput)
	json.Unmarshal(transactionBytes, &output)
	return output, nil
}

// SendRawTransaction Broadcast a serialized transaction.
func (c *AElfClient) SendRawTransaction(transaction, signature string, returnTransaction bool) (*types.SendRawTransaction, error) {
	url := c.Host + SENDRAWTRANSACTION
	params := map[string]interface{}{
		"Transaction":       transaction,
		"Signature":         signature,
		"ReturnTransaction": returnTransaction,
	}
	rawTransactionBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return nil, errors.New("Send RawTransaction error:" + err.Error())
	}
	var rawTransaction = new(types.SendRawTransaction)
	json.Unmarshal(rawTransactionBytes, &rawTransaction)
	return rawTransaction, nil
}

// SendTransactions Broadcast volume transactions.
func (c *AElfClient) SendTransactions(rawTransactions string) ([]interface{}, error) {
	url := c.Host + SENDTRANSACTIONS
	params := map[string]interface{}{
		"RawTransactions": rawTransactions,
	}
	transactionsBytes, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return nil, errors.New("Send Transaction error:" + err.Error())
	}
	var data interface{}
	json.Unmarshal(transactionsBytes, &data)
	var transactions []interface{}
	for _, d := range data.([]interface{}) {
		transactions = append(transactions, d)
	}
	return transactions, nil
}

func (c *AElfClient) CalculateTransactionFee(rawTransaction string) (*types.CalculateTransactionFee, error) {
	url := c.Host + CALCULATETRANSACTIONFEE
	params := map[string]interface{}{
		"RawTransaction": rawTransaction,
	}
	transactionFeeResult, err := utils.PostRequest(url, c.Version, params)
	if err != nil {
		return nil, errors.New("CalculateTransactionFee error:" + err.Error())
	}
	var feeResult = new(types.CalculateTransactionFee)
	json.Unmarshal(transactionFeeResult, &feeResult)
	spew.Dump("CalculateTransactionFee : ", feeResult.Success)
	return feeResult, nil

}
