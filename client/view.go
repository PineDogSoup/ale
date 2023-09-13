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
	"sync"
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
	fromAddress := c.GetAddressFromPrivateKey(privateKeyForView)
	contractAddress, _ := c.GetContractAddressByName("AElf.ContractNames.Token")
	transaction, _ := c.CreateTransaction(fromAddress, contractAddress, methodName, nil)
	signature, _ := c.SignTransaction(privateKeyForView, transaction)
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
	toAddress, err := c.GetGenesisContractAddress()
	if err != nil {
		return "", errors.New("Get Genesis contract Address error")
	}

	hashBytes, _ := proto.Marshal(&pb.Hash{Value: utils.GetBytesSha256(contractName)})

	transaction, _ := c.CreateTransaction(utils.GetAddressFromPrivateKey(privateKeyForView), toAddress, "GetContractAddressByName", hashBytes)
	signature, _ := c.SignTransaction(privateKeyForView, transaction)
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

func (c *AElfClient) GetContractInfoByAddress(address string) (*pb.ContractInfo, error) {
	res := new(pb.ContractInfo)
	toAddress, err := c.GetGenesisContractAddress()
	if err != nil {
		return res, errors.New("Get Genesis contract Address error")
	}

	addr, _ := utils.Base58StringToAddress(address)
	addrBytes, _ := proto.Marshal(addr)
	transaction, err := c.CreateTransaction(utils.GetAddressFromPrivateKey(privateKeyForView), toAddress, "GetContractInfo", addrBytes)
	if err != nil {
		return res, errors.New("Create Transaction error" + err.Error())
	}

	signature, _ := c.SignTransaction(privateKeyForView, transaction)
	if err != nil {
		return res, errors.New("Sign Transaction error" + err.Error())
	}
	transaction.Signature = signature
	transactionBytes, err := proto.Marshal(transaction)
	if err != nil {
		return res, errors.New("proto marshasl transaction error" + err.Error())
	}

	result, err := c.ExecuteTransaction(hex.EncodeToString(transactionBytes))
	if err != nil {
		return res, errors.New("Execute Transaction error" + err.Error())
	}

	resultBytes, err := hex.DecodeString(result)
	if err != nil {
		return res, errors.New("Decode error" + err.Error())
	}

	err = proto.Unmarshal(resultBytes, res)
	if err != nil {
		return res, errors.New("Unmarshal error" + err.Error())
	}
	return res, nil
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

func (c *AElfClient) GetContracts(ctx context.Context, contractNames []string) (map[string]*types.ContractInfo, error) {
	var res sync.Map
	var wg sync.WaitGroup
	wg.Add(len(contractNames))

	for i := range contractNames {
		go func(index int) {
			defer wg.Done()

			contractName := contractNames[index]
			addrStr, err := c.GetContractAddressByName(contractName)
			if err != nil {
				return
			}
			ci, err := c.GetContractInfoByAddress(addrStr)
			if err != nil {
				return
			}
			res.Store(contractName, &types.ContractInfo{
				Info:    ci,
				Address: addrStr,
			})
		}(i)
	}

	wg.Wait()

	contracts := make(map[string]*types.ContractInfo)
	res.Range(func(key, value any) bool {
		contracts[key.(string)] = value.(*types.ContractInfo)
		return true
	})
	return contracts, nil
}

func (c *AElfClient) GetContractInfo(contractName string) *types.ContractInfo {
	if name, ok := c.ContractInfo.Load(contractName); ok {
		return name.(*types.ContractInfo)
	}

	if contractName == contract.GenesisContractSystemName {
		genesisContractAddr, _ := c.GetGenesisContractAddress()
		c.ContractInfo.Store(contractName, &types.ContractInfo{Address: genesisContractAddr})
	}

	addr, _ := c.GetContractAddressByName(contractName)
	info, _ := c.GetContractInfoByAddress(addr)
	c.ContractInfo.Store(contractName, &types.ContractInfo{
		Info:    info,
		Address: addr,
	})

	res, _ := c.ContractInfo.Load(contractName)
	return res.(*types.ContractInfo)
}
