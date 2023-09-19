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
		Short: "Call portkey contract method name with contract method input by string.",
		Long: `
For example,
$ aelfctl portkey getholderinfo --env=test2 --chainId=tDVW '{"caHash": "a0a96f0c4b45719091ede2634dc05b277df4c68c39e2a3465c0c38f61a7b67fa"}'
will get caHash holder info in env=test chainId=tDVW.
`,
		Run: portkeyCommandFunc,
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
