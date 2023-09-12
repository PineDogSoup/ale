package client

import (
	"ale/core/Contract"
	"ale/core/types"
	client "ale/protobuf/generated"
	"ale/utils"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/protobuf/proto"
)

// GetTransactionPoolStatus Get information about the current transaction pool.
func (ac *AElfClient) GetTransactionPoolStatus() (*types.TransactionPoolStatusOutput, error) {
	url := ac.Host + TRANSACTIONPOOLSTATUS
	//resp, err := ac.c.Do("GET", url, nil)
	transactionPoolBytes, err := utils.GetRequest("GET", url, ac.Version, nil)
	if err != nil {
		return nil, errors.New("Get Transaction Pool Status error:" + err.Error())
	}
	var transactionPool = new(types.TransactionPoolStatusOutput)
	json.Unmarshal(transactionPoolBytes, &transactionPool)
	return transactionPool, nil
}

// GetTransactionResult Gets the result of transaction execution by the given transactionId.
func (ac *AElfClient) GetTransactionResult(transactionID string) (*types.TransactionResult, error) {
	url := ac.Host + TRANSACTIONRESULT
	_, err := hex.DecodeString(transactionID)
	if err != nil {
		return nil, errors.New("transactionID hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{"transactionId": transactionID}
	transactionBytes, err := utils.GetRequest("GET", url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Get Transaction Result error:" + err.Error())
	}
	var transaction = new(types.TransactionResult)
	json.Unmarshal(transactionBytes, &transaction)
	return transaction, nil
}

// GetTransactionResults Get results of multiple transactions by specified blockHash.
func (ac *AElfClient) GetTransactionResults(blockHash string, offset, limit int) ([]*types.TransactionResult, error) {
	url := ac.Host + TRANSACTIONRESULTS
	_, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, errors.New("blockHash hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{
		"blockHash": blockHash,
		"offset":    offset,
		"limit":     limit,
	}
	transactionsBytes, err := utils.GetRequest("GET", url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Get Transaction Results error:" + err.Error())
	}
	var datas interface{}
	json.Unmarshal(transactionsBytes, &datas)
	var transactions []*types.TransactionResult
	for _, d := range datas.([]interface{}) {
		var transaction = new(types.TransactionResult)
		Bytes, _ := json.Marshal(d)
		json.Unmarshal(Bytes, &transaction)
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// GetMerklePathByTransactionID Get merkle path of a transaction.
func (ac *AElfClient) GetMerklePathByTransactionID(transactionID string) (*types.MerklePath, error) {
	url := ac.Host + MBYTRANSACTIONID
	_, err := hex.DecodeString(transactionID)
	if err != nil {
		return nil, errors.New("transactionID hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{"transactionId": transactionID}
	merkleBytes, err := utils.GetRequest("GET", url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Get MerklePath By TransactionID error:" + err.Error())
	}
	var merkle = new(types.MerklePath)
	json.Unmarshal(merkleBytes, &merkle)
	return merkle, nil
}

// ExecuteTransaction  Call a read-only method of a contract.
func (ac *AElfClient) ExecuteTransaction(rawTransaction string) (string, error) {
	url := ac.Host + EXECUTETRANSACTION
	params := map[string]interface{}{"RawTransaction": rawTransaction}
	transactionBytes, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return "", errors.New("Execute Transaction error:" + err.Error())
	}
	return utils.BytesToString(transactionBytes), nil
}

// ExecuteRawTransaction Call a method of a contract by given serialized strings.
func (ac *AElfClient) ExecuteRawTransaction(input *types.ExecuteRawTransaction) (string, error) {
	url := ac.Host + EXECUTERAWTRANSACTION
	params := map[string]interface{}{
		"RawTransaction": input.RawTransaction,
		"Signature":      input.Signature,
	}
	transactionBytes, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return "", errors.New("Execute RawTransaction error:" + err.Error())
	}
	//var data interface{}
	//json.Unmarshal(transactionBytes, &data)
	return utils.BytesToString(transactionBytes), nil
}

// SendTransaction Broadcast a transaction.
func (ac *AElfClient) SendTransaction(transaction string) (*types.SendTransactionOutput, error) {
	url := ac.Host + SENDTRANSACTION
	params := map[string]interface{}{"RawTransaction": transaction}
	transactionBytes, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Send Transaction error:" + err.Error())
	}
	var output = new(types.SendTransactionOutput)
	json.Unmarshal(transactionBytes, &output)
	return output, nil
}

// CreateRawTransaction Creates an unsigned serialized transaction.
func (ac *AElfClient) CreateRawTransaction(input *types.CreateRawTransactionInput) (*types.CreateRawTransactionOutput, error) {
	url := ac.Host + RAWTRANSACTION
	params := map[string]interface{}{
		"From":           input.From,
		"MethodName":     input.MethodName,
		"Params":         input.Params,
		"RefBlockHash":   input.RefBlockHash,
		"RefBlockNumber": input.RefBlockNumber,
		"To":             input.To,
	}
	transactionBytes, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Create RawTransaction error:" + err.Error())
	}
	var output = new(types.CreateRawTransactionOutput)
	json.Unmarshal(transactionBytes, &output)
	return output, nil
}

// SendRawTransaction Broadcast a serialized transaction.
func (ac *AElfClient) SendRawTransaction(transaction, signature string, returnTransaction bool) (*types.SendRawTransaction, error) {
	url := ac.Host + SENDRAWTRANSACTION
	params := map[string]interface{}{
		"Transaction":       transaction,
		"Signature":         signature,
		"ReturnTransaction": returnTransaction,
	}
	rawTransactionBytes, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return nil, errors.New("Send RawTransaction error:" + err.Error())
	}
	var rawTransaction = new(types.SendRawTransaction)
	json.Unmarshal(rawTransactionBytes, &rawTransaction)
	return rawTransaction, nil
}

// SendTransactions Broadcast volume transactions.
func (ac *AElfClient) SendTransactions(rawTransactions string) ([]interface{}, error) {
	url := ac.Host + SENDTRANSACTIONS
	params := map[string]interface{}{
		"RawTransactions": rawTransactions,
	}
	transactionsBytes, err := utils.PostRequest(url, ac.Version, params)
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

func (ac *AElfClient) CalculateTransactionFee(rawTransaction string) (*types.CalculateTransactionFee, error) {
	url := ac.Host + CALCULATETRANSACTIONFEE
	params := map[string]interface{}{
		"RawTransaction": rawTransaction,
	}
	transactionFeeResult, err := utils.PostRequest(url, ac.Version, params)
	if err != nil {
		return nil, errors.New("CalculateTransactionFee error:" + err.Error())
	}
	var feeResult = new(types.CalculateTransactionFee)
	json.Unmarshal(transactionFeeResult, &feeResult)
	spew.Dump("CalculateTransactionFee : ", feeResult.Success)
	return feeResult, nil

}

func (ac *AElfClient) GetTransactionFees(transactionResultDto types.TransactionResult) (map[string][]map[string]interface{}, error) {
	var feeDicts = map[string][]map[string]interface{}{}
	eventLogs := transactionResultDto.Logs
	if len(eventLogs) == 0 {
		return nil, errors.New("transaction Result Dto not found  Logs error")
	}
	for _, log := range eventLogs {
		nonIndexedBytes, _ := utils.Base64DecodeBytes(log.NonIndexed)
		if log.Name == "TransactionFeeCharged" {
			var feeCharged = new(client.TransactionFeeCharged)
			proto.Unmarshal(nonIndexedBytes, feeCharged)
			var feeMap = map[string]interface{}{feeCharged.Symbol: feeCharged.Amount}
			feeDicts["TransactionFeeCharged"] = append(feeDicts["TransactionFeeCharged"], feeMap)
		}

		if log.Name == "ResourceTokenCharged" {
			var tokenCharged = new(client.ResourceTokenCharged)
			proto.Unmarshal(nonIndexedBytes, tokenCharged)
			var feeMap = map[string]interface{}{tokenCharged.Symbol: tokenCharged.Amount}
			feeDicts["ResourceTokenCharged"] = append(feeDicts["ResourceTokenCharged"], feeMap)
		}
	}
	return feeDicts, nil
}

func (ac *AElfClient) GetTransferred(txId string) []*client.Transferred {
	transffereds := make([]*client.Transferred, 0)
	result, err := ac.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return transffereds
	}

	contractAddr, _ := ac.GetContractAddressByName(Contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == Contract.TransferredLogEventName && log.Address == contractAddr {
			transferred := new(client.Transferred)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, transferred)
			}
			if fromBytes, err := utils.Base64DecodeBytes(log.Indexed[0]); err == nil {
				temp := new(client.Transferred)
				proto.Unmarshal(fromBytes, temp)
				transferred.From = temp.From
			}
			if toBytes, err := utils.Base64DecodeBytes(log.Indexed[1]); err == nil {
				temp := new(client.Transferred)
				proto.Unmarshal(toBytes, temp)
				transferred.To = temp.To
			}
			if symbolBytes, err := utils.Base64DecodeBytes(log.Indexed[2]); err == nil {
				temp := new(client.Transferred)
				proto.Unmarshal(symbolBytes, temp)
				transferred.Symbol = temp.Symbol
			}
			transffereds = append(transffereds, transferred)
		}
	}

	return transffereds
}

func (ac *AElfClient) GetCrossChainTransferred(txId string) []*client.CrossChainTransferred {
	crossChainTransferreds := make([]*client.CrossChainTransferred, 0)
	result, err := ac.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return crossChainTransferreds
	}

	contractAddr, _ := ac.GetContractAddressByName(Contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == Contract.CrossChainTransferredLogEventName && log.Address == contractAddr {
			crossChainTransferred := new(client.CrossChainTransferred)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, crossChainTransferred)
			}
			crossChainTransferreds = append(crossChainTransferreds, crossChainTransferred)
		}
	}

	return crossChainTransferreds
}

func (ac *AElfClient) GetCrossChainReceived(txId string) []*client.CrossChainReceived {
	crossChainReceiveds := make([]*client.CrossChainReceived, 0)
	result, err := ac.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return crossChainReceiveds
	}

	contractAddr, _ := ac.GetContractAddressByName(Contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == Contract.CrossChainReceivedLogEventName && log.Address == contractAddr {
			crossChainReceived := new(client.CrossChainReceived)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, crossChainReceived)
			}
			crossChainReceiveds = append(crossChainReceiveds, crossChainReceived)
		}
	}

	return crossChainReceiveds
}
