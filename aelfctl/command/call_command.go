package command

import (
	"ale/pkg/cobrautl"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func NewCallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "call",
		Run: callCommandFunc,
	}
	return cmd
}

func callCommandFunc(cmd *cobra.Command, args []string) {
	method, value := getOp(args)
	resp, err := newSendClientFromCmd(cmd).Call(context.Background(), method, value)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	res, _ := json.Marshal(resp)
	fmt.Println(string(res))
}
