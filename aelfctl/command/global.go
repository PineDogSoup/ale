package command

import (
	"ale/client"
	"ale/pkg/cobrautl"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

const (
	DefaultContextTimeOut = 15 * time.Second
)

type GlobalFlags struct {
	Endpoints     []string
	WalletAddress string
}

func mustClientFromCmd(cmd *cobra.Command) *client.AElfClient {
	cfg := &client.Config{}

	eps, err := cmd.Flags().GetStringSlice("endpoints")
	if err == nil {
		for i, ip := range eps {
			eps[i] = strings.TrimSpace(ip)
		}
	}
	cfg.Endpoints = eps
	c, err := client.New(cfg)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadConnection, err)
	}

	return c
}
