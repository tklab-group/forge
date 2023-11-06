package mold

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/moldfile/generator"
	"log/slog"
	"os"
	"path"
)

var dockerfilePath string
var moldfilePath string

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mold PATH",
		Short: `Generate moldfile.`,
		Long:  `Generate moldfile from existing Dockerfile and build context. Using PATH as build context.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildContext := args[0]

			if dockerfilePath == "" {
				dockerfilePath = path.Join(buildContext, "Dockerfile")
			}
			if moldfilePath == "" {
				moldfilePath = path.Join(buildContext, "Dockerfile.mold")
			}

			moldfile, err := generator.GenerateMoldfile(dockerfilePath, buildContext)
			if err != nil {
				return err
			}

			slog.Info(fmt.Sprintf("write to %s", moldfilePath))

			f, err := os.Create(moldfilePath)
			if err != nil {
				return fmt.Errorf("failed to create `%s`: %v", moldfilePath, err)
			}
			defer f.Close()

			_, err = f.WriteString(moldfile.ToString())
			if err != nil {
				return fmt.Errorf("failed to write to `%s`: %v", moldfilePath, err)
			}

			slog.Info("Done")

			return nil
		},
	}

	cmd.SetIn(config.In)
	cmd.SetOut(config.Out)
	cmd.SetErr(config.Err)

	cmd.Flags().StringVarP(&dockerfilePath, "dockerfile", "d", "", `Name of the Dockerfile (default: "PATH/Dockerfile")`)
	cmd.Flags().StringVarP(&moldfilePath, "moldfile", "m", "", `Name of the Moldfile (default: "PATH/Dockerfile.mold")`)

	return cmd
}
