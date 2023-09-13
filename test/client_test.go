package test

import (
	"ale/client"
	"ale/core/contract"
	"ale/core/types"
	pb "ale/protobuf/generated"
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	secp256 "github.com/haltingstate/secp256k1-go"
	"google.golang.org/protobuf/proto"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mainClient, _ = client.New(&client.Config{
		Endpoints:  []string{"http://192.168.66.61:8000"},
		Version:    "1.0",
		PrivateKey: "cd86ab6347d8e52bbbe8532141fc59ce596268143a308d1d40fedf385528b458",
		Timeout:    0,
	})

	sideClient, _ = client.New(&client.Config{
		Endpoints:  []string{"http://127.0.0.1:8000"},
		Version:    "1.0",
		PrivateKey: "cd86ab6347d8e52bbbe8532141fc59ce596268143a308d1d40fedf385528b458",
		Timeout:    0,
	})

	networkInfo, _ = mainClient.GetNetworkInfo()
	peers, _       = mainClient.GetPeers(true)
)

var _address = mainClient.GetAddressFromPrivateKey(mainClient.PrivateKey)

func TestGetAddressFromPubKey(t *testing.T) {
	privateKeyBytes, _ := hex.DecodeString(mainClient.PrivateKey)
	pubkeyBytes := secp256.UncompressedPubkeyFromSeckey(privateKeyBytes)
	pubKeyAddress := mainClient.GetAddressFromPubKey(hex.EncodeToString(pubkeyBytes))
	assert.Equal(t, _address, pubKeyAddress)
	spew.Dump("Get Address From Public Key", pubKeyAddress)
}

func TestGetFormattedAddress(t *testing.T) {
	formattedAddress, err := mainClient.GetFormattedAddress(_address)
	assert.NoError(t, err)
	spew.Dump("Get Formatted Address result", formattedAddress, err)

	privateKeyaddress := mainClient.GetAddressFromPrivateKey(mainClient.PrivateKey)
	assert.Equal(t, formattedAddress, ("ELF_" + privateKeyaddress + "_AELF"))
}

func TestGenerateKeyPairInfo(t *testing.T) {
	keyPair := mainClient.GenerateKeyPairInfo()
	spew.Dump("Generate KeyPair Info Result", keyPair)
}

func TestGetContractAddressByName(t *testing.T) {
	contractAddress, err := mainClient.GetContractAddressByName(contract.TokenContractSystemName)
	assert.NoError(t, err)
	spew.Dump("Get ContractAddress By Name Result", contractAddress)
}

func TestClient(t *testing.T) {
	//Test IsConnected
	isConnected := mainClient.IsConnected()
	assert.True(t, isConnected)

	wrongClient := client.AElfClient{
		Host:       "http://127.0.0.1:8008",
		Version:    "1.0",
		PrivateKey: "cd86ab6347d8e52bbbe8532141fc59ce596268143a308d1d40fedf385528b458",
	}

	isConnected = wrongClient.IsConnected()
	assert.False(t, isConnected)
}

func TestGetTransactionFee(t *testing.T) {
	var result types.TransactionResult
	var logEventDto types.LogEvent
	logEventDto.Name = "TransactionFeeCharged"
	var param = &pb.TransactionFeeCharged{Symbol: "ELF", Amount: 1000}
	paramBytes, _ := proto.Marshal(param)
	logEventDto.NonIndexed = base64.StdEncoding.EncodeToString(paramBytes)
	result.Logs = append(result.Logs, logEventDto)

	logEventDto.Name = "ResourceTokenCharged"
	var params = &pb.ResourceTokenCharged{Symbol: "READ", Amount: 800}
	paramsBytes, _ := proto.Marshal(params)
	logEventDto.NonIndexed = base64.StdEncoding.EncodeToString(paramsBytes)
	result.Logs = append(result.Logs, logEventDto)

	logEventDto.Name = "ResourceTokenCharged"
	params = &pb.ResourceTokenCharged{Symbol: "WRITE", Amount: 600}
	paramsBytes, _ = proto.Marshal(params)
	logEventDto.NonIndexed = base64.StdEncoding.EncodeToString(paramsBytes)
	result.Logs = append(result.Logs, logEventDto)

	logEventDto.Name = "ResourceTokenCharged"
	params = &pb.ResourceTokenCharged{Symbol: "READ", Amount: 200}
	paramsBytes, _ = proto.Marshal(params)
	logEventDto.NonIndexed = base64.StdEncoding.EncodeToString(paramsBytes)
	result.Logs = append(result.Logs, logEventDto)

	res, _ := mainClient.GetTransactionFees(result)
	assert.Equal(t, int64(1000), res["TransactionFeeCharged"][0]["ELF"])
	assert.Equal(t, int64(800), res["ResourceTokenCharged"][0]["READ"])
	assert.Equal(t, int64(600), res["ResourceTokenCharged"][1]["WRITE"])
	assert.Equal(t, int64(200), res["ResourceTokenCharged"][2]["READ"])
}

func TestGetTokenInfo(t *testing.T) {
	tokenInfo, err := mainClient.GetTokenInfo(DefaultTestSymbol)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenInfo)
	assert.Equal(t, DefaultTestSymbol, tokenInfo.Symbol)
	assert.Equal(t, DefaultTestTokenTotalSupply, tokenInfo.TotalSupply)
	assert.Equal(t, DefaultTestTokenDecimals, tokenInfo.Decimals)
	assert.Equal(t, DefaultTestTokenIsBurnable, tokenInfo.IsBurnable)
	assert.Equal(t, DefaultTestTokenIssueChainId, tokenInfo.IssueChainId)
}

func TestGetContracts(t *testing.T) {
	contractNames := []string{contract.TokenContractSystemName, contract.CrossChainContractSystemName}
	contracts, _ := mainClient.GetContracts(context.Background(), contractNames)
	assert.Equal(t, len(contracts), len(contractNames))
}
