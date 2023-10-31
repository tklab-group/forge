package cli

import (
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/cli/mold"
	"os"
)

func newRootCmd(config config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "", // TODO
		Long:  "", // TODO
	}
	rootCmd.SetIn(config.In)
	rootCmd.SetOut(config.Out)
	rootCmd.SetErr(config.Err)
	
	rootCmd.AddCommand(
		mold.Cmd(config),
	)

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
