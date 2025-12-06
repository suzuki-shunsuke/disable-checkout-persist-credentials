package cli

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	flag "github.com/spf13/pflag"
	"github.com/suzuki-shunsuke/disable-checkout-persist-credentials/pkg/controller"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

const help = `disable-checkout-persist-credentials - Disable actions/checkout persist-credentials.
https://github.com/suzuki-shunsuke/disable-checkout-persist-credentials

Usage:
	disable-checkout-persist-credentials [<options>] [file ...]

files: GitHub Actions files. By default, \.github/workflows/*\.ya?ml$

Options:
	--help, -h       Show help
	--version, -v    Show version
	--log-level      The log level (default: info, debug, info, warn, error)
	--log-color      The log color (always, never, auto)
	--dry-run        Dry Run`

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
	Logger  *slogutil.Logger
}

type LDFlags struct {
	Version string
}

func (r *Runner) Run(logger *slogutil.Logger) error {
	flg := &Flag{}
	parseFlags(flg)
	if flg.Version {
		fmt.Fprintln(r.Stdout, r.LDFlags.Version)
		return nil
	}
	if flg.Help {
		fmt.Fprintln(r.Stdout, help)
		return nil
	}
	if err := logger.SetLevel(flg.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	ctrl := &controller.Controller{}
	ctrl.Init(afero.NewOsFs(), r.Stdout, r.Stderr)
	return ctrl.Run(logger.Logger, &controller.Input{ //nolint:wrapcheck
		DryRun: flg.DryRun,
		Args:   flg.Args,
	})
}

type Flag struct {
	LogLevel string
	LogColor string
	Args     []string
	Help     bool
	Version  bool
	DryRun   bool
}

func parseFlags(f *Flag) {
	flag.StringVar(&f.LogLevel, "log-level", "info", "The log level")
	flag.StringVar(&f.LogColor, "log-color", "auto", "The log color")
	flag.BoolVarP(&f.Help, "help", "h", false, "Show help")
	flag.BoolVarP(&f.Version, "version", "v", false, "Show version")
	flag.BoolVar(&f.DryRun, "dry-run", false, "Dry Run")
	flag.Parse()
	f.Args = flag.Args()
}
