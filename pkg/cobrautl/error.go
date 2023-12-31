package cobrautl

import (
	"fmt"
	"os"
)

const (
	ExitSuccess = iota
	ExitError
	ExitBadConnection
	ExitInvalidInput // for txn, watch command
	ExitBadFeature   // provided a valid flag with an unsupported value
	ExitInterrupted
	ExitIO
	ExitBadArgs = 128

	ExitServerError       = 4
	ExitClusterNotHealthy = 5
)

func ExitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}
