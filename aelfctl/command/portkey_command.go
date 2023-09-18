package command

import (
	"ale/pkg/cobrautl"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	environment = ""
	chainId     = ""
)

func NewPortkeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portkey [options] <methodName> <inputString>",
		Short: "Call contract method name with contract method input by string.",
		Run:   portkeyCommandFunc,
	}

	cmd.PersistentFlags().StringVar(&environment, "env", "test1", "query environment")
	cmd.PersistentFlags().StringVar(&chainId, "chainId", "AELF", "query chainId")

	return cmd
}

func portkeyCommandFunc(cmd *cobra.Command, args []string) {
	method, value := getOp(args)
	resp, err := mustPortkeyFromCmd(cmd).Call(context.Background(), method, value)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	res, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(res))
}
