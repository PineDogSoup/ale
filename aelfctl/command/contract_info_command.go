package command

import (
	"ale/core/consts"
	"ale/core/types"
	"ale/pkg/cobrautl"
	"ale/utils"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var (
	contractNameList = []string{
		consts.TokenContractSystemName,
		consts.CrossChainContractSystemName,
		consts.EconomicContractSystemName,
		consts.ProfitContractSystemName,
		consts.TreasuryContractSystemName,
		consts.ElectionContractSystemName,
		consts.VoteContractSystemName,
		consts.ConsensusContractSystemName,
	}
)

func NewContractInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contracts [options]",
		Short: "Get the default contract ContractName & ContractVersion & Author & ContractAddress",
		Long: `
For example,
$ contracts --endpoint="http://127.0.0.1:8000" 
will get all contract info in your chain.
`,
		Run: contractInfoCommandFunc,
	}
	return cmd
}

func contractInfoCommandFunc(cmd *cobra.Command, args []string) {
	resp, err := mustClientFromCmd(cmd).GetContracts(context.Background(), contractNameList)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}

	contractInfos(resp)
}

func contractInfos(contractInfos map[string]*types.ContractInfo) {
	var rows [][]string
	hdr := []string{"contract name", "version", "author", "contract address"}
	for name, contract := range contractInfos {
		rows = append(rows, []string{
			name,
			fmt.Sprintf(contract.Info.ContractVersion),
			fmt.Sprint(utils.AddressToBase58String(contract.Info.Author)),
			fmt.Sprint(contract.Address),
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
