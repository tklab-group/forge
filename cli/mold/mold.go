package mold

import (
	"github.com/spf13/cobra"
	"github.com/tklab-group/forge/cli/config"
	"github.com/tklab-group/forge/moldfile/generator"
)

var dockerfilePath string

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mold [-F --file] PATH",
		Short: `Generate moldfile.`,
		Long:  `Generate moldfile from existing Dockerfile and build context. Using PATH as build context.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			buildContext := args[0]

			if dockerfilePath == "" {
				dockerfilePath = "Dockerfile"
			}

			err := generator.GenerateMoldfile(dockerfilePath, buildContext)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.SetIn(config.In)
	cmd.SetOut(config.Out)
	cmd.SetErr(config.Err)

	cmd.Flags().StringVarP(&dockerfilePath, "file", "F", "", `Name of the Dockerfile (default: "PATH/Dockerfile")`)

	return cmd
}
