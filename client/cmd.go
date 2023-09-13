package client

import "C"
import (
	"ale/core/contract"
	"ale/core/types"
	pb "ale/protobuf/generated"
	"ale/utils"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"
	"sync"
)

var (
	contractNameList = []string{
		contract.TokenContractSystemName,
		contract.CrossChainContractSystemName,
		contract.ForestContractSystemName,
		contract.EconomicContractSystemName,
		contract.ProfitContractSystemName,
		contract.TreasuryContractSystemName,
		contract.ElectionContractSystemName,
		contract.VoteContractSystemName,
		contract.ConsensusContractSystemName,
	}
)

type AElfCmd interface {
	Send(ctx context.Context, methodName, inputStr string) (types.SendResult, error)
	Call(ctx context.Context, methodName, inputStr string) (types.CallResult, error)
	GetContracts(ctx context.Context, contractNames []string) (map[string]*types.ContractInfo, error)
}

func (c *Client) Send(ctx context.Context, methodName, inputStr string) (types.SendResult, error) {
	var res types.SendResult
	var err error

	switch strings.ToLower(methodName) {
	case strings.ToLower(contract.TokenContractTransfer):
		res, err = c.sendTransferTransaction(ctx, inputStr)
		break
	default:
		err = errors.New(fmt.Sprintf("No support method:%s", methodName))
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) Call(ctx context.Context, methodName, inputStr string) (types.CallResult, error) {
	var res types.CallResult
	var err error

	switch strings.ToLower(methodName) {
	default:
		err = errors.New(fmt.Sprintf("No support method:%s", methodName))
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) GetContracts(ctx context.Context, contractNames []string) (map[string]*types.ContractInfo, error) {
	var res sync.Map
	var wg sync.WaitGroup
	wg.Add(len(contractNames))

	for i := range contractNames {
		go func(index int) {
			defer wg.Done()

			contractName := contractNames[index]
			addrStr, err := c.AElf.GetContractAddressByName(contractName)
			if err != nil {
				return
			}
			ci, err := c.AElf.GetContractInfoByAddress(addrStr)
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

func (c *Client) sendTransferTransaction(ctx context.Context, inputStr string) (types.SendResult, error) {
	var res types.SendResult
	var rev types.TransferInput

	err := json.Unmarshal([]byte(inputStr), &rev)
	if err != nil {
		return res, errors.New(fmt.Sprintf("Unmarshal send transaction error:%s", err.Error()))
	}

	params := &pb.TransferInput{
		Symbol: rev.Symbol,
		Amount: rev.Amount,
		Memo:   rev.Memo,
	}
	params.To, _ = utils.Base58StringToAddress(rev.To)
	paramsByte, _ := proto.Marshal(params)

	tokenContractAddress, _ := c.AElf.GetContractAddressByName(contract.TokenContractSystemName)
	transaction, _ := c.AElf.CreateTransaction(c.AElf.GetAddressFromPrivateKey(c.AElf.PrivateKey), tokenContractAddress, contract.TokenContractTransfer, paramsByte)
	transaction.Signature, _ = c.AElf.SignTransaction(c.AElf.PrivateKey, transaction)
	transactionByets, _ := proto.Marshal(transaction)
	sendResult, err := c.AElf.SendTransaction(hex.EncodeToString(transactionByets))
	if err != nil {
		return res, nil
	}
	return types.SendResult{TransactionId: sendResult.TransactionID}, nil
}
