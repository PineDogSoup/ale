package client

import (
	"ale/core/contract"
	"ale/core/types"
	pb "ale/protobuf/generated"
	"ale/utils"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	wrap "github.com/golang/protobuf/ptypes/wrappers"
	secp256 "github.com/haltingstate/secp256k1-go"
	"google.golang.org/protobuf/proto"
)

// GetAddressFromPubKey Get the account address through the public key.
func (c *AElfClient) GetAddressFromPubKey(pubkey string) string {
	bytes, _ := hex.DecodeString(pubkey)
	return utils.GetAddressByBytes(bytes)
}

// GetAddressFromPrivateKey Get the account address through the private key.
func (c *AElfClient) GetAddressFromPrivateKey(privateKey string) string {
	bytes, _ := hex.DecodeString(privateKey)
	pubkeyBytes := secp256.UncompressedPubkeyFromSeckey(bytes)
	return utils.GetAddressByBytes(pubkeyBytes)
}

// GetFormattedAddress Convert the Address to the displayed string:symbol_base58-string_base58-string-chain-id.
func (c *AElfClient) GetFormattedAddress(address string) (string, error) {
	chain, _ := c.GetChainStatus()
	methodName := "GetPrimaryTokenSymbol"
	fromAddress := c.GetAddressFromPrivateKey(ExamplePrivateKey)
	contractAddress, _ := c.GetContractAddressByName("AElf.ContractNames.Token")
	transaction, _ := c.CreateTransaction(fromAddress, contractAddress, methodName, nil)
	signature, _ := c.SignTransaction(ExamplePrivateKey, transaction)
	transaction.Signature = signature
	transactionBytes, err := proto.Marshal(transaction)
	if err != nil {
		return "", errors.New("proto marshasl transaction error" + err.Error())
	}
	executeResult, _ := c.ExecuteTransaction(hex.EncodeToString(transactionBytes))
	var symbol = new(wrap.StringValue)
	executeBytes, err := hex.DecodeString(executeResult)
	proto.Unmarshal(executeBytes, symbol)
	return symbol.Value + "_" + address + "_" + chain.ChainId, nil
}

func (c *AElfClient) GetTokenBalance(symbol, owner string) (*pb.GetBalanceOutput, error) {
	tokenContractAddr, _ := c.GetContractAddressByName(contract.TokenContractSystemName)
	addr := c.GetAddressFromPrivateKey(c.PrivateKey)
	ownerAddr, err := utils.Base58StringToAddress(owner)
	if err != nil {
		return &pb.GetBalanceOutput{}, err
	}
	inputByte, _ := proto.Marshal(&pb.GetBalanceInput{
		Symbol: symbol,
		Owner:  ownerAddr,
	})

	tx, _ := c.CreateTransaction(addr, tokenContractAddr, contract.TokenContractGetBalance, inputByte)
	sign, _ := c.SignTransaction(c.PrivateKey, tx)
	tx.Signature = sign

	txByets, _ := proto.Marshal(tx)
	re, _ := c.ExecuteTransaction(hex.EncodeToString(txByets))

	balance := &pb.GetBalanceOutput{}
	bytes, _ := hex.DecodeString(re)
	proto.Unmarshal(bytes, balance)

	return balance, nil
}

func (c *AElfClient) GetTokenInfo(symbol string) (*pb.TokenInfo, error) {
	tokenContractAddr, _ := c.GetContractAddressByName(contract.TokenContractSystemName)
	addr := c.GetAddressFromPrivateKey(c.PrivateKey)
	inputByte, _ := proto.Marshal(&pb.TokenInfo{
		Symbol: symbol,
	})

	tx, _ := c.CreateTransaction(addr, tokenContractAddr, contract.TokenContractGetTokenInfo, inputByte)
	sign, _ := c.SignTransaction(c.PrivateKey, tx)
	tx.Signature = sign

	txBytes, _ := proto.Marshal(tx)
	re, _ := c.ExecuteTransaction(hex.EncodeToString(txBytes))

	tokenInfo := &pb.TokenInfo{}
	bytes, _ := hex.DecodeString(re)
	proto.Unmarshal(bytes, tokenInfo)

	return tokenInfo, nil
}

// GetContractAddressByName Get  contract address by contract name.
func (c *AElfClient) GetContractAddressByName(contractName string) (string, error) {
	fromAddress := c.GetAddressFromPrivateKey(ExamplePrivateKey)
	toAddress, err := c.GetGenesisContractAddress()
	if err != nil {
		return "", errors.New("Get Genesis contract Address error")
	}
	contractNameBytes := utils.GetBytesSha256(contractName)
	var hash = new(pb.Hash)
	hash.Value = contractNameBytes
	hashBytes, _ := proto.Marshal(hash)

	transaction, _ := c.CreateTransaction(fromAddress, toAddress, "GetContractAddressByName", hashBytes)
	signature, _ := c.SignTransaction(ExamplePrivateKey, transaction)
	transaction.Signature = signature
	transactionBytes, err := proto.Marshal(transaction)
	if err != nil {
		return "", errors.New("proto marshasl transaction error" + err.Error())
	}
	result, _ := c.ExecuteTransaction(hex.EncodeToString(transactionBytes))
	var address = new(pb.Address)
	resultBytes, err := hex.DecodeString(result)
	proto.Unmarshal(resultBytes, address)
	return utils.EncodeCheck(address.Value), nil
}

// GetGenesisContractAddress Get the address of genesis contract.
func (c *AElfClient) GetGenesisContractAddress() (string, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return "", errors.New("Get Genesis contract Address error:" + err.Error())
	}
	address := chainStatus.GenesisContractAddress
	return address, nil
}

// GetTransactionPoolStatus Get information about the current transaction pool.
func (c *AElfClient) GetTransactionPoolStatus() (*types.TransactionPoolStatusOutput, error) {
	url := c.Host + TRANSACTIONPOOLSTATUS
	transactionPoolBytes, err := utils.GetRequest("GET", url, c.Version, nil)
	if err != nil {
		return nil, errors.New("Get Transaction Pool Status error:" + err.Error())
	}
	var transactionPool = new(types.TransactionPoolStatusOutput)
	json.Unmarshal(transactionPoolBytes, &transactionPool)
	return transactionPool, nil
}

// GetTransactionResult Gets the result of transaction execution by the given transactionId.
func (c *AElfClient) GetTransactionResult(transactionID string) (*types.TransactionResult, error) {
	url := c.Host + TRANSACTIONRESULT
	_, err := hex.DecodeString(transactionID)
	if err != nil {
		return nil, errors.New("transactionID hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{"transactionId": transactionID}
	transactionBytes, err := utils.GetRequest("GET", url, c.Version, params)
	if err != nil {
		return nil, errors.New("Get Transaction Result error:" + err.Error())
	}
	var transaction = new(types.TransactionResult)
	json.Unmarshal(transactionBytes, &transaction)
	return transaction, nil
}

// GetTransactionResults Get results of multiple transactions by specified blockHash.
func (c *AElfClient) GetTransactionResults(blockHash string, offset, limit int) ([]*types.TransactionResult, error) {
	url := c.Host + TRANSACTIONRESULTS
	_, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, errors.New("blockHash hex to []byte error:" + err.Error())
	}
	params := map[string]interface{}{
		"blockHash": blockHash,
		"offset":    offset,
		"limit":     limit,
	}
	transactionsBytes, err := utils.GetRequest("GET", url, c.Version, params)
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

func (c *AElfClient) GetContracts(ctx context.Context, contractNames []string) ([]*types.ContractInfo, error) {
	var res []*types.ContractInfo
	fromAddress := c.GetAddressFromPrivateKey(ExamplePrivateKey)
	toAddress, err := c.GetGenesisContractAddress()
	if err != nil {
		return res, errors.New("Get Genesis contract Address error")
	}

	for _, contractName := range contractNames {
		go func(contractName string) {

		}(contractName)
		var hash = new(pb.Hash)
		hash.Value = utils.GetBytesSha256(contractName)
		hashBytes, _ := proto.Marshal(hash)

		transaction, _ := c.CreateTransaction(fromAddress, toAddress, "GetContractAddressByName", hashBytes)
		transaction.Signature, _ = c.SignTransaction(ExamplePrivateKey, transaction)

		transactionBytes, err := proto.Marshal(transaction)
		if err != nil {
			return res, errors.New("proto marshasl transaction error" + err.Error())
		}
		result, _ := c.ExecuteTransaction(hex.EncodeToString(transactionBytes))
		var contractInfo = new(pb.ContractInfo)
		resultBytes, err := hex.DecodeString(result)
		proto.Unmarshal(resultBytes, contractInfo)

		addr, _ := c.GetContractAddressByName(contractName)

		res = append(res, &types.ContractInfo{
			ContractName: contractName,
			Info:         contractInfo,
			Address:      addr,
		})
	}
	return res, nil
}
