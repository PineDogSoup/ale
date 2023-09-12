package command

import (
	"ale/core/contract"
	"ale/core/types"
	"ale/pkg/cobrautl"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var (
	contractNameList = []string{contract.TokenContractSystemName, contract.CrossChainContractSystemName}
)

func NewContractInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "contracts",
		Run: contractInfoCommandFunc,
	}
	return cmd
}

func contractInfoCommandFunc(cmd *cobra.Command, args []string) {
	ctx, cancel := commandCtx(cmd)
	resp, err := mustClientFromCmd(cmd).GetContracts(ctx, contractNameList)
	cancel()
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	contractInfos(resp)
}

func contractInfos(contractInfos []*types.ContractInfo) {
	var rows [][]string
	hdr := []string{"contract name", "version", "author", "address", "isSystemContract"}
	for _, contract := range contractInfos {
		rows = append(rows, []string{
			contract.ContractName,
			fmt.Sprintf("%x", contract.Info.Version),
			fmt.Sprint(contract.Info.Author),
			fmt.Sprint(contract.Address),
			fmt.Sprint(contract.Info.IsSystemContract),
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(hdr)
	for _, row := range rows {
		table.Append(row)
	}
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.Render()
}
