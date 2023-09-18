package client

import "C"
import (
	"ale/core/consts"
	"ale/core/types"
	"ale/core/types/api"
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
	"time"
)

var (
	domainMap = map[string]string{
		"test1":   "https://localtest-applesign.portkey.finance",
		"test2":   "https://localtest-applesign2.portkey.finance",
		"testnet": "https://did-portkey-test.portkey.finance",
		"main":    "https://did-portkey.portkey.finance",
	}
	defaultInterval = 5 * time.Second
)

type Action interface {
	Send(ctx context.Context, methodName, inputStr string) (*types.SendResult, error)
	Call(ctx context.Context, methodName, inputStr string) (*types.CallResult, error)
}

func (c *Client) Send(ctx context.Context, methodName, inputStr string) (*types.SendResult, error) {
	var res *types.SendResult
	var err error

	switch strings.ToLower(methodName) {
	case strings.ToLower(consts.TokenContractTransfer):
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

func (c *Client) Call(ctx context.Context, methodName, inputStr string) (*types.CallResult, error) {
	res := new(types.CallResult)

	switch strings.ToLower(methodName) {
	case strings.ToLower("GetBalance"):
		holderInfoRes, err := c.getBalance(ctx, inputStr)
		if err != nil {
			return &types.CallResult{Message: fmt.Sprintf("%s: %s", methodName, err.Error())}, err
		}
		res.Data = holderInfoRes
		break
	default:
		err := errors.New(fmt.Sprintf("No support method:%s", methodName))
		return &types.CallResult{Message: fmt.Sprintf("%s: %s", methodName, err.Error())}, err
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

func (c *Client) GetChainInfos(ctx context.Context) (map[string]*api.SearchChainInfo, error) {
	var res sync.Map
	var wg sync.WaitGroup
	wg.Add(len(domainMap))

	for env := range domainMap {
		go func(e string) {
			defer wg.Done()

			resBytes, _ := utils.GetRequest(domainMap[e]+consts.CHAININFOURL, "", nil)
			var datas string
			var output = new(api.SearchChainInfo)
			json.Unmarshal(resBytes, &datas)
			json.Unmarshal([]byte(datas), &output)
			res.Store(e, output)
		}(env)

	}

	wg.Wait()

	chainInfo := make(map[string]*api.SearchChainInfo)
	res.Range(func(key, value any) bool {
		chainInfo[key.(string)] = value.(*api.SearchChainInfo)
		return true
	})
	return chainInfo, nil
}

func (c *Client) sendTransferTransaction(ctx context.Context, inputStr string) (*types.SendResult, error) {
	var res *types.SendResult
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

	tokenContractAddress, _ := c.AElf.GetContractAddressByName(consts.TokenContractSystemName)
	transaction, _ := c.AElf.CreateTransaction(c.AElf.GetAddressFromPrivateKey(c.AElf.PrivateKey), tokenContractAddress, consts.TokenContractTransfer, paramsByte)
	transaction.Signature, _ = c.AElf.SignTransaction(c.AElf.PrivateKey, transaction)
	transactionByets, _ := proto.Marshal(transaction)
	sendResult, err := c.AElf.SendTransaction(hex.EncodeToString(transactionByets))
	if err != nil {
		return res, nil
	}
	return &types.SendResult{TransactionId: sendResult.TransactionID}, nil
}

func (c *Client) getBalance(ctx context.Context, inputStr string) (*pb.GetBalanceOutput, error) {
	res := &pb.GetBalanceOutput{}
	var rev types.GetBalanceInput

	err := json.Unmarshal([]byte(inputStr), &rev)
	if err != nil {
		return res, errors.New(fmt.Sprintf("Unmarshal send transaction error:%s", err.Error()))
	}
	ownerAddr, err := utils.Base58StringToAddress(rev.Owner)
	if err != nil {
		return res, errors.New(fmt.Sprintf("String to address error:%s", err.Error()))
	}
	paramsByte, _ := proto.Marshal(&pb.GetBalanceInput{
		Symbol: rev.Symbol,
		Owner:  ownerAddr,
	})

	contractAddr, _ := c.AElf.GetContractAddressByName(consts.TokenContractSystemName)
	transaction, _ := c.AElf.CreateTransaction(c.AElf.GetAddressFromPrivateKey(privateKeyForView), contractAddr, consts.TokenContractGetBalance, paramsByte)
	transaction.Signature, _ = c.AElf.SignTransaction(privateKeyForView, transaction)
	transactionByets, _ := proto.Marshal(transaction)
	callResult, err := c.AElf.ExecuteTransaction(hex.EncodeToString(transactionByets))
	if err != nil {
		return res, errors.New(fmt.Sprintf("Execute transaction error:%s", err.Error()))
	}
	resByte, _ := hex.DecodeString(callResult)
	proto.Unmarshal(resByte, res)
	return res, nil
}

func (c *Client) GetTxResultUntilFinished(ctx context.Context, txId string) (*types.TransactionResult, error) {
	ticker := time.NewTicker(defaultInterval)
	defer ticker.Stop()
	for {
		res, _ := c.getTxResult(ctx, txId)
		if strings.ToUpper(res.Status) == "MINED" || res.Status == "FAILED" {
			return res, nil
		}
		fmt.Printf("txId:%s, status:%s\n", txId, res.Status)
		select {
		case <-ctx.Done():
			return nil, errors.New("get transaction result context done " + ctx.Err().Error())
		case <-ticker.C:
		}
	}
}

func (c *Client) getTxResult(ctx context.Context, txId string) (*types.TransactionResult, error) {
	res, err := c.AElf.GetTransactionResult(txId)
	if err != nil {
		return nil, errors.New("transaction fail, err:" + err.Error())
	}
	return res, err
}
