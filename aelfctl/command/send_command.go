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
		Use: "send",
		Run: sendCommandFunc,
		Long: `
When <value> begins with '-', <value> is interpreted as a flag.
Insert '--' for workaround:
`,
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
