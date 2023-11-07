package check

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/moldfile/parser"
	"io"
	"log/slog"
	"os"
)

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check FILE_PATH",
		Short: "Check if the Dockerfile/Moldfile is available format.",
		Long: `Check if the Dockerfile/Moldfile is appropriate format to parse with forge.
If FILE_PATH is "-", the file content is read from stdin.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			var r io.Reader
			if filePath == "-" {
				slog.Info("reading the file content from stdin")
				r = cmd.InOrStdin()
			} else {
				f, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer f.Close()
				r = f
			}

			moldfile, err := parser.ParseMoldFile(r)
			if err != nil {
				return fmt.Errorf("failed to parse: %v", err)
			}

			if moldfile.BuildStageCount() == 0 {
				return fmt.Errorf("unexpected parsing result. Is the file empty?")
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "ok")
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
