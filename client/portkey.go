package client

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
)

type Portkey struct {
	AElf              AElfClient
	CaContractAddress string
}

func NewPortkey(env, chainId string) (*Portkey, error) {
	if domain, ok := domainMap[env]; ok {
		resBytes, _ := utils.GetRequest(domain+consts.CHAININFOURL, "", nil)
		var datas string
		var output = new(api.SearchChainInfo)
		json.Unmarshal(resBytes, &datas)
		json.Unmarshal([]byte(datas), &output)
		for _, chain := range output.Items {
			if chain.ChainId == chainId {
				return &Portkey{CaContractAddress: chain.CaContractAddress, AElf: AElfClient{
					Host: chain.Endpoint,
				}}, nil
			}
		}
		return nil, errors.New("not support chainId")
	}
	return nil, errors.New("not support environment")
}

func (c *Portkey) Call(ctx context.Context, methodName, inputStr string) (*types.CallResult, error) {
	res := new(types.CallResult)

	switch strings.ToLower(methodName) {
	case strings.ToLower("GetHolderInfo"):
		holderInfoRes, err := c.getHolderInfo(ctx, inputStr)
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

func (c *Portkey) Send(ctx context.Context, methodName, inputStr string) (*types.SendResult, error) {
	var res *types.SendResult
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

func (c *Portkey) getHolderInfo(ctx context.Context, inputStr string) (*pb.GetHolderInfoOutput, error) {
	res := &pb.GetHolderInfoOutput{}
	var rev types.GetHolderInfoInput

	err := json.Unmarshal([]byte(inputStr), &rev)
	if err != nil {
		return res, errors.New(fmt.Sprintf("Unmarshal send transaction error:%s", err.Error()))
	}

	paramsByte, _ := proto.Marshal(&pb.GetHolderInfoInput{
		CaHash: &pb.Hash{Value: utils.HexStringToByteArray(rev.CaHash)},
	})
	transaction, _ := c.AElf.CreateTransaction(c.AElf.GetAddressFromPrivateKey(privateKeyForView), c.CaContractAddress, consts.PortkeyContractGetHolderInfo, paramsByte)
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
