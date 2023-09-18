package test

import (
	"ale/core/consts"
	"ale/utils"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContracts(t *testing.T) {
	contractNames := []string{consts.TokenContractSystemName, consts.CrossChainContractSystemName}
	contracts, _ := mainClient.GetContracts(context.Background(), contractNames)
	assert.Equal(t, len(contracts), len(contractNames))
}

func TestSendTransfer(t *testing.T) {
	paramsJson := fmt.Sprintf(`{"to":"%s","symbol":"%s","amount":%d,"memo":"%s"}`, utils.AddressToBase58String(defaultTestHolder.Address), DefaultTestSymbol, DefaultTransferTestAmount, DefaultTransferTestMemo)
	res, err := mainClient.Send(context.Background(), consts.TokenContractTransfer, paramsJson)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.NotEmpty(t, res.TransactionId)
}
