package main

import (
	"os"

	"github.com/suzuki-shunsuke/disable-checkout-persist-credentials/pkg/cli"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

var version = ""

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "disable-checkout-persist-credentials",
		Version: version,
	})
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
		},
		Logger: logger,
	}
	if err := runner.Run(logger); err != nil {
		slogerr.WithError(logger.Logger, err).Error("disable-checkout-persist-credentials failed")
		return 1
	}
	return 0
}
