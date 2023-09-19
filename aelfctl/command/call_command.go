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
		Use:   "call [options] <methodName> <inputString>",
		Short: "Call contract method name with contract method input by string.",
		Long: `
For example,
$ call getBalance --endpoint="http://127.0.0.1:8000" '{"symbol":"ELF", "owner":"wWnNNjrUiveHWrmogtQnQvzJBbiqerwndzM7WuvU1kxLfZC6Z"}'
will get <owner> balance of ELF in your local chain.
`,
		Run: callCommandFunc,
	}

	return cmd
}

func callCommandFunc(cmd *cobra.Command, args []string) {
	method, value := getOp(args)
	resp, err := mustClientFromCmd(cmd).Call(context.Background(), method, value)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	res, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(res))
}
