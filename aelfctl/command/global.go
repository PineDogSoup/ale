package command

import (
	"ale/client"
	"ale/pkg/cobrautl"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type GlobalFlags struct {
	Endpoints     []string
	WalletAddress string
	PrivateKey    string
}

func newSendClientFromCmd(cmd *cobra.Command) *client.Client {
	c := mustClientFromCmd(cmd)
	pk, err := cmd.Flags().GetString("privateKey")
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadConnection, err)
	}
	c.AElf.PrivateKey = pk
	return c
}

func mustClientFromCmd(cmd *cobra.Command) *client.Client {
	cfg := &client.Config{}

	Endpoints, err := cmd.Flags().GetStringSlice("endpoints")
	if err == nil {
		for i, ip := range Endpoints {
			Endpoints[i] = strings.TrimSpace(ip)
		}
		cfg.Endpoints = Endpoints
	}
	c, err := client.New(cfg)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadConnection, err)
	}

	return c
}

func mustPortkeyFromCmd(cmd *cobra.Command) *client.Portkey {
	env, err := cmd.Flags().GetString("env")
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitInvalidInput, err)
	}
	chainId, err := cmd.Flags().GetString("chainId")
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitInvalidInput, err)
	}

	p, err := client.NewPortkey(env, chainId)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadConnection, err)
	}

	return p
}

func argOrStdin(args []string, stdin io.Reader, i int) (string, error) {
	if i < len(args) {
		return args[i], nil
	}
	bytes, err := ioutil.ReadAll(stdin)
	if string(bytes) == "" || err != nil {
		return "", errors.New("no available argument and stdin")
	}
	return string(bytes), nil
}

func getOp(args []string) (string, string) {
	if len(args) == 0 {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("command needs 1 argument and input from stdin or 2 arguments"))
	}

	value, err := argOrStdin(args, os.Stdin, 1)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("command needs 1 argument and input from stdin or 2 arguments"))
	}

	return args[0], value
}
