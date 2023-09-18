package command

import (
	"ale/pkg/cobrautl"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [options] <methodName> <inputString>",
		Short: "Send transaction with contract method input by string.",
		Run:   sendCommandFunc,
	}
	return cmd
}

func sendCommandFunc(cmd *cobra.Command, args []string) {
	key, value := getOp(args)
	sendClient := newSendClientFromCmd(cmd)
	resp, err := sendClient.Send(context.Background(), key, value)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	res, _ := json.Marshal(resp)
	fmt.Println(string(res))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	txRes, _ := sendClient.GetTxResultUntilFinished(ctx, resp.TransactionId)
	tx, _ := json.MarshalIndent(txRes, "", "  ")
	fmt.Println(string(tx))
}
