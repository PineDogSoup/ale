package client

import (
	"ale/core/contract"
	"ale/core/types"
	pb "ale/protobuf/generated"

	"ale/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	secp256 "github.com/haltingstate/secp256k1-go"
	"google.golang.org/protobuf/proto"
)

// SignTransaction Sign a transaction using private key.
func (c *AElfClient) SignTransaction(privateKey string, transaction *pb.Transaction) ([]byte, error) {
	transactionBytes, _ := proto.Marshal(transaction)
	txDataBytes := sha256.Sum256(transactionBytes)
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	signatureBytes := secp256.Sign(txDataBytes[:], privateKeyBytes)
	return signatureBytes, nil
}

// CreateTransaction create a transaction from the input parameters.
func (c *AElfClient) CreateTransaction(from, to, method string, params []byte) (*pb.Transaction, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return nil, errors.New("Get Chain Status error ")
	}
	prefixBytes, _ := hex.DecodeString(chainStatus.BestChainHash)
	fromAddressBytes, _ := utils.Base58StringToAddress(from)
	toAddressBytes, _ := utils.Base58StringToAddress(to)
	var transaction = &pb.Transaction{
		From:           fromAddressBytes,
		To:             toAddressBytes,
		MethodName:     method,
		RefBlockNumber: chainStatus.BestChainHeight,
		RefBlockPrefix: prefixBytes[:4],
		Params:         params,
	}
	return transaction, nil
}

func (c *AElfClient) GetTransactionFees(transactionResultDto types.TransactionResult) (map[string][]map[string]interface{}, error) {
	var feeDicts = map[string][]map[string]interface{}{}
	eventLogs := transactionResultDto.Logs
	if len(eventLogs) == 0 {
		return nil, errors.New("transaction Result Dto not found  Logs error")
	}
	for _, log := range eventLogs {
		nonIndexedBytes, _ := utils.Base64DecodeBytes(log.NonIndexed)
		if log.Name == "TransactionFeeCharged" {
			var feeCharged = new(pb.TransactionFeeCharged)
			proto.Unmarshal(nonIndexedBytes, feeCharged)
			var feeMap = map[string]interface{}{feeCharged.Symbol: feeCharged.Amount}
			feeDicts["TransactionFeeCharged"] = append(feeDicts["TransactionFeeCharged"], feeMap)
		}

		if log.Name == "ResourceTokenCharged" {
			var tokenCharged = new(pb.ResourceTokenCharged)
			proto.Unmarshal(nonIndexedBytes, tokenCharged)
			var feeMap = map[string]interface{}{tokenCharged.Symbol: tokenCharged.Amount}
			feeDicts["ResourceTokenCharged"] = append(feeDicts["ResourceTokenCharged"], feeMap)
		}
	}
	return feeDicts, nil
}

func (c *AElfClient) GetTransferred(txId string) []*pb.Transferred {
	transffereds := make([]*pb.Transferred, 0)
	result, err := c.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return transffereds
	}

	contractAddr, _ := c.GetContractAddressByName(contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == contract.TransferredLogEventName && log.Address == contractAddr {
			transferred := new(pb.Transferred)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, transferred)
			}
			if fromBytes, err := utils.Base64DecodeBytes(log.Indexed[0]); err == nil {
				temp := new(pb.Transferred)
				proto.Unmarshal(fromBytes, temp)
				transferred.From = temp.From
			}
			if toBytes, err := utils.Base64DecodeBytes(log.Indexed[1]); err == nil {
				temp := new(pb.Transferred)
				proto.Unmarshal(toBytes, temp)
				transferred.To = temp.To
			}
			if symbolBytes, err := utils.Base64DecodeBytes(log.Indexed[2]); err == nil {
				temp := new(pb.Transferred)
				proto.Unmarshal(symbolBytes, temp)
				transferred.Symbol = temp.Symbol
			}
			transffereds = append(transffereds, transferred)
		}
	}

	return transffereds
}

func (c *AElfClient) GetCrossChainTransferred(txId string) []*pb.CrossChainTransferred {
	crossChainTransferreds := make([]*pb.CrossChainTransferred, 0)
	result, err := c.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return crossChainTransferreds
	}

	contractAddr, _ := c.GetContractAddressByName(contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == contract.CrossChainTransferredLogEventName && log.Address == contractAddr {
			crossChainTransferred := new(pb.CrossChainTransferred)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, crossChainTransferred)
			}
			crossChainTransferreds = append(crossChainTransferreds, crossChainTransferred)
		}
	}

	return crossChainTransferreds
}

func (c *AElfClient) GetCrossChainReceived(txId string) []*pb.CrossChainReceived {
	crossChainReceiveds := make([]*pb.CrossChainReceived, 0)
	result, err := c.GetTransactionResult(txId)
	if err != nil || len(result.Logs) == 0 {
		return crossChainReceiveds
	}

	contractAddr, _ := c.GetContractAddressByName(contract.TokenContractSystemName)

	for _, log := range result.Logs {
		if log.Name == contract.CrossChainReceivedLogEventName && log.Address == contractAddr {
			crossChainReceived := new(pb.CrossChainReceived)
			if nonIndexedBytes, err := utils.Base64DecodeBytes(log.NonIndexed); err == nil {
				proto.Unmarshal(nonIndexedBytes, crossChainReceived)
			}
			crossChainReceiveds = append(crossChainReceiveds, crossChainReceived)
		}
	}

	return crossChainReceiveds
}
