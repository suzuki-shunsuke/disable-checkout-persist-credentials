package main

import (
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/disable-checkout-persist-credentials/pkg/cli"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "disable-checkout-persist-credentials",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	runner := cli.Runner{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		Logger:      logger,
		LogLevelVar: logLevelVar,
	}
	if err := runner.Run(); err != nil {
		slogerr.WithError(logger, err).Error("disable-checkout-persist-credentials failed")
		return 1
	}
	return 0
}
