package command

import (
	"ale/pkg/cobrautl"
	"context"
	"github.com/spf13/cobra"
)

func commandCtx(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	timeOut, err := cmd.Flags().GetDuration("command-timeout")
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	return context.WithTimeout(context.Background(), timeOut)
}
