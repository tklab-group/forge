package mold

import (
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/moldfile/generator"
	"path"
)

var dockerfilePath string
var moldfilePath string

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mold [OPTIONS] PATH",
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

			err := generator.GenerateMoldfile(dockerfilePath, buildContext, moldfilePath)
			if err != nil {
				return err
			}

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
