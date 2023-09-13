package main

import (
	"ale/aelfctl/command"
	"ale/pkg/cobrautl"
	"github.com/spf13/cobra"
)

var (
	globalFlags = command.GlobalFlags{}
)

func initCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "aelfctl",
		Short: "aelf ctl",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"http://127.0.0.1:8000"}, "aelf node endpoints")

	root.AddCommand(
		//command.NewCallCommand(),
		command.NewSendCommand(),
		command.NewContractInfoCommand(),
	)
	return root
}

func main() {
	rootCmd := initCmd()
	rootCmd.SetUsageFunc(usageFunc)

	if c, err := rootCmd.ExecuteC(); err != nil {
		rootCmd.Println(c.UsageString())
	}
}

func usageFunc(c *cobra.Command) error {
	return cobrautl.UsageFunc(c)
}
