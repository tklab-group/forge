package generator

import (
	"fmt"
	"github.com/tklab-group/forge/docker/image"
	"github.com/tklab-group/forge/moldfile/parser"
	"log/slog"
	"os"
)

func GenerateMoldfile(dockerfilePath string, buildContext string) error {
	slog.Info("reading Dockerfile")
	dockerfile, err := os.Open(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to open Dockerfile: %v", err)
	}
	defer dockerfile.Close()

	slog.Info("parsing Dockerfile")
	moldfile, err := parser.ParseMoldFile(dockerfile)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile: %v", err)
	}

	for i := 0; i < moldfile.BuildStageCount(); i++ {
		slog.Info(fmt.Sprintf("Molding build stage index=%d", i))

		err = moldPerBuildStage(moldfile, buildContext, i)
		if err != nil {
			return fmt.Errorf("failed to mold build stage index=%d: %v", i, err)
		}
	}

	return nil
}

func moldPerBuildStage(moldFile parser.MoldFile, buildContext string, stageIndex int) error {
	target, err := moldFile.GetBuildStage(stageIndex)
	if err != nil {
		return fmt.Errorf("failed to get target BuildStage: %v", err)
	}

	err = moldBaseImage(target)
	if err != nil {
		return fmt.Errorf("failed to mold base image: %v", err)
	}

	// TODO: Mold package version

	return nil
}

func moldBaseImage(buildStage parser.BuildStage) error {
	fromInstruction, err := buildStage.GetFromInstruction()
	if err != nil {
		return fmt.Errorf("failed to get FROM instruction: %v", err)
	}

	latestDigest, err := image.GetLatestDigest(fromInstruction.ImageInfoString())
	if err != nil {
		return fmt.Errorf("failed to get latest digest of `%s`: %v", fromInstruction.ImageInfoString(), latestDigest)
	}

	fromInstruction.UpdateImageInfo(latestDigest)

	return nil
}
