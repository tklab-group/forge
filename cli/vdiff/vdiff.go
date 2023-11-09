package vdiff

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/moldfile/parser"
	"log/slog"
	"os"
)

func Cmd(config config.Config) *cobra.Command {
	// TODO: Support other options of output format

	cmd := &cobra.Command{
		Use:   "vdiff FILE_PATH1 FILE_PATH2",
		Short: "Extract version differences between two files",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath1 := args[0]
			filePath2 := args[1]

			slog.Info(fmt.Sprintf("parsing %s as Moldfile", filePath1))
			moldfile1, err := getMoldfileFromFilePath(filePath1)
			if err != nil {
				return fmt.Errorf("failed to parse %s as Moldfile: %v", filePath1, err)
			}

			slog.Info(fmt.Sprintf("parsing %s as Moldfile", filePath2))
			moldfile2, err := getMoldfileFromFilePath(filePath2)
			if err != nil {
				return fmt.Errorf("failed to parse %s as Moldfile: %v", filePath2, err)
			}

			slog.Info("starting vdiff")
			vdiff, err := parser.VDiffMoldfiles(moldfile1, moldfile2)
			if err != nil {
				return fmt.Errorf("failed to vdiff: %v", err)
			}

			rawJson, err := json.Marshal(vdiff)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(rawJson))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetIn(config.In)
	cmd.SetOut(config.Out)
	cmd.SetErr(config.Err)

	return cmd
}

func getMoldfileFromFilePath(filePath string) (parser.MoldFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	moldfile, err := parser.ParseMoldFile(f)
	if err != nil {
		return nil, err
	}

	return moldfile, nil
}
