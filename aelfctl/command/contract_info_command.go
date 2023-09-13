package command

import (
	"ale/core/contract"
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
		contract.TokenContractSystemName,
		contract.CrossChainContractSystemName,
		contract.ForestContractSystemName,
		contract.EconomicContractSystemName,
		contract.ProfitContractSystemName,
		contract.TreasuryContractSystemName,
		contract.ElectionContractSystemName,
		contract.VoteContractSystemName,
		contract.ConsensusContractSystemName,
	}
)

func NewContractInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contracts [options]",
		Short: "Get the default contract ContractName & ContractVersion & Author & ContractAddress",
		Run:   contractInfoCommandFunc,
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
