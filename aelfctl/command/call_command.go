package command

import "github.com/spf13/cobra"

func NewCallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "call",
		Run: callCommandFunc,
	}
	return cmd
}

func callCommandFunc(cmd *cobra.Command, args []string) {}
