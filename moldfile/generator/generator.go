package generator

import (
	"fmt"
	"github.com/tklab-group/forge/docker/image"
	"github.com/tklab-group/forge/moldfile/parser"
	"log/slog"
	"os"

	"github.com/tklab-group/docker-image-disassembler/disassembler"
)

func GenerateMoldfile(dockerfilePath string, buildContext string) (parser.MoldFile, error) {
	slog.Info("reading Dockerfile")
	dockerfile, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Dockerfile: %v", err)
	}
	defer dockerfile.Close()

	slog.Info("parsing Dockerfile")
	moldfile, err := parser.ParseMoldFile(dockerfile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Dockerfile: %v", err)
	}

	for i := 0; i < moldfile.BuildStageCount(); i++ {
		slog.Info(fmt.Sprintf("molding build stage index=%d", i))

		err = moldPerBuildStage(moldfile, buildContext, i)
		if err != nil {
			return nil, fmt.Errorf("failed to mold build stage index=%d: %v", i, err)
		}
	}

	slog.Info("finish molding")

	return moldfile, nil
}

func moldPerBuildStage(moldFile parser.MoldFile, buildContext string, stageIndex int) error {
	target, err := moldFile.GetBuildStage(stageIndex)
	if err != nil {
		return fmt.Errorf("failed to get target BuildStage: %v", err)
	}

	slog.Info("molding base image")
	err = moldBaseImage(target)
	if err != nil {
		return fmt.Errorf("failed to mold base image: %v", err)
	}

	priorBuildStages := make([]parser.BuildStage, 0)
	for i := 0; i < stageIndex; i++ {
		priorBuildStage, err := moldFile.GetBuildStage(i)
		if err != nil {
			return fmt.Errorf("failed to get prior BuildStage index=%d: %v", i, err)
		}

		priorBuildStages = append(priorBuildStages, priorBuildStage)
	}

	slog.Info("molding package versions")
	err = moldPackageVersion(target, priorBuildStages, buildContext)
	if err != nil {
		return fmt.Errorf("failed to mold package versions: %v", err)
	}

	return nil
}

func moldBaseImage(buildStage parser.BuildStage) error {
	fromInstruction, err := buildStage.GetFromInstruction()
	if err != nil {
		return fmt.Errorf("failed to get FROM instruction: %v", err)
	}

	currentBaseImage := fromInstruction.ImageInfoString()
	if currentBaseImage == "scratch" {
		slog.Info("scratch has no digest")
		return nil
	}

	latestDigest, err := image.GetLatestDigest(currentBaseImage)
	if err != nil {
		return fmt.Errorf("failed to get latest digest of `%s`: %v", currentBaseImage, latestDigest)
	}

	fromInstruction.UpdateImageInfo(latestDigest)

	return nil
}

func moldPackageVersion(buildStage parser.BuildStage, priorBuildStages []parser.BuildStage, buildContext string) error {
	tmpDockerfile, err := os.CreateTemp("/tmp", "Dockerfile")
	if err != nil {
		return err
	}
	defer os.Remove(tmpDockerfile.Name())

	for _, priorBuildStage := range priorBuildStages {
		_, err = tmpDockerfile.WriteString(priorBuildStage.ToString())
		if err != nil {
			return err
		}
	}

	_, err = tmpDockerfile.WriteString(buildStage.ToString())
	if err != nil {
		return err
	}

	err = tmpDockerfile.Close()
	if err != nil {
		return err
	}

	iid, err := image.BuildImageWithCLI(tmpDockerfile.Name(), buildContext)
	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("iid: %s", iid))

	// TODO: Support other package managers

	aptPkgInfo, err := disassembler.GetAptPkgInfoInImageFromImageID(iid)
	if err != nil {
		return fmt.Errorf("failed to disassemble: %v", err)
	}

	packageVersions := map[string]map[string]string{
		"apt": aptPkgInfo,
	}

	buildStage.UpdatePackageInfos(packageVersions)

	return nil
}
