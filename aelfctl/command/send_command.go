package command

import (
	"ale/pkg/cobrautl"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
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
	resp, err := newSendClientFromCmd(cmd).Send(context.Background(), key, value)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	res, _ := json.Marshal(resp)
	fmt.Println(string(res))
}
