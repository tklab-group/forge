package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/check"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/cli/mold"
	"github.com/tklab-group/forge/cli/vdiff"
)

var logLevel string

func newRootCmd(config config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "forge",
		Short:        "", // TODO
		Long:         "", // TODO,
		SilenceUsage: true,
		Version:      "v0.0.5",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := settingLog(logLevel)
			if err != nil {
				return err
			}

			return nil
		},
	}
	rootCmd.SetIn(config.In)
	rootCmd.SetOut(config.Out)
	rootCmd.SetErr(config.Err)

	rootCmd.AddCommand(
		mold.Cmd(config),
		check.Cmd(config),
		vdiff.Cmd(config),
	)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", `Set the logging level ("debug", "info", "warn", "error") (default "info")`)

	return rootCmd
}

func Execute() error {
	conf := config.Config{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}

	err := newRootCmd(conf).Execute()
	if err != nil {
		return err
	}

	return nil
}

func settingLog(logLevelStr string) error {
	logLevels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level, ok := logLevels[logLevelStr]
	if !ok {
		return fmt.Errorf("unexpected log level: %s", logLevelStr)
	}

	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Format time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(time.DateTime))
		}
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replace,
	}))

	slog.SetDefault(l)

	return nil
}
