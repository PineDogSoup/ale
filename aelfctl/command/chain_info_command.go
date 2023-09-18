package command

import (
	"ale/core/types/api"
	"ale/pkg/cobrautl"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func NewChainInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chains [options]",
		Short: "Get all chain env chainId chainName portkeyContractAddress endpoint",
		Run:   chainInfoCommandFunc,
	}
	return cmd
}

func chainInfoCommandFunc(cmd *cobra.Command, args []string) {
	resp, err := mustClientFromCmd(cmd).GetChainInfos(context.Background())
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	chainInfos(resp)
}

func chainInfos(contractInfos map[string]*api.SearchChainInfo) {
	var rows [][]string
	hdr := []string{"env", "chainId", "chainName", "portkey contract address", "endpoint"}
	for env, Info := range contractInfos {
		for _, chainInfo := range Info.Items {
			rows = append(rows, []string{
				env,
				fmt.Sprintf(chainInfo.ChainId),
				fmt.Sprint(chainInfo.ChainName),
				fmt.Sprint(chainInfo.CaContractAddress),
				fmt.Sprint(chainInfo.Endpoint),
			})

		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(hdr)
	for _, row := range rows {
		table.Append(row)
	}
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Render()
}
