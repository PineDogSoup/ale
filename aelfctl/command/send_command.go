package command

import "github.com/spf13/cobra"

var (
	privateKey    string
	walletAddress string
)

func NewSendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "send",
		Run: sendCommandFunc,
	}
	cmd.Flags().StringVar(&walletAddress, "wallet", "", "wallet address")
	cmd.Flags().StringVar(&privateKey, "privateKey", "", "wallet privateKey")
	return cmd
}

func sendCommandFunc(cmd *cobra.Command, args []string) {

}
