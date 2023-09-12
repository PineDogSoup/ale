package client

import (
	"ale/core/Contract"
	"ale/core/types"
	pb "ale/protobuf/generated"
	"ale/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/golang/protobuf/proto"
	wrap "github.com/golang/protobuf/ptypes/wrappers"
	secp256 "github.com/haltingstate/secp256k1-go"
)

// AElfClient AElf Client.
type AElfClient struct {
	Host       string
	Version    string
	PrivateKey string
	c          *utils.HttpClient
}

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

	ExamplePrivateKey = "680afd630d82ae5c97942c4141d60b8a9fedfa5b2864fca84072c17ee1f72d9d"
)

func NewAElfClient(privateKey string, c *utils.HttpClient) *AElfClient {
	return &AElfClient{
		PrivateKey: privateKey,
		c:          c,
	}
}

// GetAddressFromPubKey Get the account address through the public key.
func (ac *AElfClient) GetAddressFromPubKey(pubkey string) string {
	bytes, _ := hex.DecodeString(pubkey)
	return utils.GetAddressByBytes(bytes)
}

// GetAddressFromPrivateKey Get the account address through the private key.
func (ac *AElfClient) GetAddressFromPrivateKey(privateKey string) string {
	bytes, _ := hex.DecodeString(privateKey)
	pubkeyBytes := secp256.UncompressedPubkeyFromSeckey(bytes)
	return utils.GetAddressByBytes(pubkeyBytes)
}

// GetFormattedAddress Convert the Address to the displayed string:symbol_base58-string_base58-string-chain-id.
func (ac *AElfClient) GetFormattedAddress(address string) (string, error) {
	chain, _ := ac.GetChainStatus()
	methodName := "GetPrimaryTokenSymbol"
	fromAddress := ac.GetAddressFromPrivateKey(ExamplePrivateKey)
	contractAddress, _ := ac.GetContractAddressByName("AElf.ContractNames.Token")
	transaction, _ := ac.CreateTransaction(fromAddress, contractAddress, methodName, nil)
	signature, _ := ac.SignTransaction(ExamplePrivateKey, transaction)
	transaction.Signature = signature
	transactionBytes, err := proto.Marshal(transaction)
	if err != nil {
		return "", errors.New("proto marshasl transaction error" + err.Error())
	}
	executeResult, _ := ac.ExecuteTransaction(hex.EncodeToString(transactionBytes))
	var symbol = new(wrap.StringValue)
	executeBytes, err := hex.DecodeString(executeResult)
	proto.Unmarshal(executeBytes, symbol)
	return symbol.Value + "_" + address + "_" + chain.ChainId, nil
}

func (ac *AElfClient) GetTokenBalance(symbol, owner string) (*pb.GetBalanceOutput, error) {
	tokenContractAddr, _ := ac.GetContractAddressByName(Contract.TokenContractSystemName)
	addr := ac.GetAddressFromPrivateKey(ac.PrivateKey)
	ownerAddr, err := utils.Base58StringToAddress(owner)
	if err != nil {
		return &pb.GetBalanceOutput{}, err
	}
	inputByte, _ := proto.Marshal(&pb.GetBalanceInput{
		Symbol: symbol,
		Owner:  ownerAddr,
	})

	tx, _ := ac.CreateTransaction(addr, tokenContractAddr, Contract.TokenContractGetBalance, inputByte)
	sign, _ := ac.SignTransaction(ac.PrivateKey, tx)
	tx.Signature = sign

	txByets, _ := proto.Marshal(tx)
	re, _ := ac.ExecuteTransaction(hex.EncodeToString(txByets))

	balance := &pb.GetBalanceOutput{}
	bytes, _ := hex.DecodeString(re)
	proto.Unmarshal(bytes, balance)

	return balance, nil
}

func (ac *AElfClient) GetTokenInfo(symbol string) (*pb.TokenInfo, error) {
	tokenContractAddr, _ := ac.GetContractAddressByName(Contract.TokenContractSystemName)
	addr := ac.GetAddressFromPrivateKey(ac.PrivateKey)
	inputByte, _ := proto.Marshal(&pb.TokenInfo{
		Symbol: symbol,
	})

	tx, _ := ac.CreateTransaction(addr, tokenContractAddr, Contract.TokenContractGetTokenInfo, inputByte)
	sign, _ := ac.SignTransaction(ac.PrivateKey, tx)
	tx.Signature = sign

	txBytes, _ := proto.Marshal(tx)
	re, _ := ac.ExecuteTransaction(hex.EncodeToString(txBytes))

	tokenInfo := &pb.TokenInfo{}
	bytes, _ := hex.DecodeString(re)
	proto.Unmarshal(bytes, tokenInfo)

	return tokenInfo, nil
}

// GetContractAddressByName Get  contract address by contract name.
func (ac *AElfClient) GetContractAddressByName(contractName string) (string, error) {
	fromAddress := ac.GetAddressFromPrivateKey(ExamplePrivateKey)
	toAddress, err := ac.GetGenesisContractAddress()
	if err != nil {
		return "", errors.New("Get Genesis Contract Address error")
	}
	contractNameBytes := utils.GetBytesSha256(contractName)
	var hash = new(pb.Hash)
	hash.Value = contractNameBytes
	hashBytes, _ := proto.Marshal(hash)

	transaction, _ := ac.CreateTransaction(fromAddress, toAddress, "GetContractAddressByName", hashBytes)
	signature, _ := ac.SignTransaction(ExamplePrivateKey, transaction)
	transaction.Signature = signature
	transactionBytes, err := proto.Marshal(transaction)
	if err != nil {
		return "", errors.New("proto marshasl transaction error" + err.Error())
	}
	result, _ := ac.ExecuteTransaction(hex.EncodeToString(transactionBytes))
	var address = new(pb.Address)
	resultBytes, err := hex.DecodeString(result)
	proto.Unmarshal(resultBytes, address)
	return utils.EncodeCheck(address.Value), nil
}

// SignTransaction Sign a transaction using private key.
func (ac *AElfClient) SignTransaction(privateKey string, transaction *pb.Transaction) ([]byte, error) {
	transactionBytes, _ := proto.Marshal(transaction)
	txDataBytes := sha256.Sum256(transactionBytes)
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	signatureBytes := secp256.Sign(txDataBytes[:], privateKeyBytes)
	return signatureBytes, nil
}

// CreateTransaction create a transaction from the input parameters.
func (ac *AElfClient) CreateTransaction(from, to, method string, params []byte) (*pb.Transaction, error) {
	chainStatus, err := ac.GetChainStatus()
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

// GetGenesisContractAddress Get the address of genesis contract.
func (ac *AElfClient) GetGenesisContractAddress() (string, error) {
	chainStatus, err := ac.GetChainStatus()
	if err != nil {
		return "", errors.New("Get Genesis Contract Address error:" + err.Error())
	}
	address := chainStatus.GenesisContractAddress
	return address, nil
}

// IsConnected Verify whether this sdk successfully connects the chain.
func (ac *AElfClient) IsConnected() bool {
	data, err := ac.GetChainStatus()
	if err != nil || data == nil {
		return false
	}
	return true
}

// GenerateKeyPairInfo Generate KeyPair Info.
func (ac *AElfClient) GenerateKeyPairInfo() *types.KeyPairInfo {
	publicKeyBytes, privateKeyBytes := secp256.GenerateKeyPair()
	publicKey := hex.EncodeToString(publicKeyBytes)
	privateKey := hex.EncodeToString(privateKeyBytes)
	privateKeyAddress := ac.GetAddressFromPrivateKey(privateKey)
	var keyPair = &types.KeyPairInfo{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    privateKeyAddress,
	}
	return keyPair
}

// GetSignatureWithPrivateKey Get Signature With PrivateKey.
func GetSignatureWithPrivateKey(privateKey string, txData []byte) (string, error) {
	privateKeyBytes, _ := hex.DecodeString(privateKey)
	txDataBytes := sha256.Sum256(txData)
	signatureBytes := secp256.Sign(txDataBytes[:], privateKeyBytes)
	return hex.EncodeToString(signatureBytes), nil
}
