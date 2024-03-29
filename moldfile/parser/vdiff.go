package parser

import (
	"fmt"
	"log/slog"
)

type VDiff struct {
	BuildStages []VDiffBuildStage `json:"buildStages"`
}

type VDiffBuildStage struct {
	StageName string             `json:"stageName"`
	BaseImage VDiffBaseImage     `json:"baseImage"`
	Packages  []VDiffPackageInfo `json:"packages"`
}

type VDiffBaseImage struct {
	Name      string `json:"name"`
	Moldfile1 string `json:"moldfile1"`
	Moldfile2 string `json:"moldfile2"`
}

type VDiffPackageInfo struct {
	PackageManager string `json:"packageManager"`
	Name           string `json:"name"`
	Moldfile1      string `json:"moldfile1"`
	Moldfile2      string `json:"moldfile2"`
}

func VDiffMoldfiles(moldfile1 MoldFile, moldfile2 MoldFile) (VDiff, error) {
	slog.Info("start vdiff moldfiles")

	if moldfile1.BuildStageCount() != moldfile2.BuildStageCount() {
		return VDiff{}, fmt.Errorf("build stage counts don't match")
	}

	m1 := moldfile1.(*moldFile)
	m2 := moldfile2.(*moldFile)

	buildStages := make([]VDiffBuildStage, 0)
	for i := 0; i < m1.BuildStageCount(); i++ {
		slog.Info(fmt.Sprintf("vdiff build stage index=%d", i))
		buildStage1 := m1.buildStages[i]
		buildStage2 := m2.buildStages[i]

		vdiff, err := vdiffBuildStages(buildStage1, buildStage2)
		if err != nil {
			return VDiff{}, fmt.Errorf("failed to vdiff build stage index %d: %v", i, err)
		}

		buildStages = append(buildStages, vdiff)
	}

	return VDiff{
		BuildStages: buildStages,
	}, nil
}

func vdiffBuildStages(buildStage1 BuildStage, buildStage2 BuildStage) (VDiffBuildStage, error) {
	bs1 := buildStage1.(*buildStage)
	bs2 := buildStage2.(*buildStage)

	slog.Debug(fmt.Sprintf("instruction counts of buildStage1: %d", bs1.instructions))
	slog.Debug(fmt.Sprintf("instruction counts of buildStage2: %d", bs2.instructions))

	buildStageName, err := vdiffBuildStageName(bs1, bs1)
	if err != nil {
		return VDiffBuildStage{}, fmt.Errorf("faild to get build stage name: %v", err)
	}

	baseImage, err := vdiffBaseImages(bs1, bs2)
	if err != nil {
		return VDiffBuildStage{}, fmt.Errorf("failed to vdiff base images: %v", err)
	}

	packages, err := vdiffPackages(bs1, bs2)
	if err != nil {
		return VDiffBuildStage{}, fmt.Errorf("failed to vdiff packages: %v", err)
	}

	return VDiffBuildStage{
		StageName: buildStageName,
		BaseImage: baseImage,
		Packages:  packages,
	}, nil
}

func vdiffBuildStageName(bs1 *buildStage, bs2 *buildStage) (string, error) {
	if len(bs1.instructions) == 0 {
		return "", fmt.Errorf("BuildStage should have more than 0 instruction, but buildStage1 is 0")
	}
	if len(bs2.instructions) == 0 {
		return "", fmt.Errorf("BuildStage should have more than 0 instruction, but buildStage2 is 0")
	}

	fromInstruction1, err := bs1.getFromInstruction()
	if err != nil {
		return "", fmt.Errorf("failed to get FROM instruction from buildStage1: %v", err)
	}

	fromInstruction2, err := bs2.getFromInstruction()
	if err != nil {
		return "", fmt.Errorf("failed to get FROM instruction from buildStage2: %v", err)
	}

	var buildStageName string
	if fromInstruction1.buildStageInfo.HasValue() {
		buildStageName = fromInstruction1.buildStageInfo.ValueOrZero().name.toString()
	} else if fromInstruction2.buildStageInfo.HasValue() {
		buildStageName = fromInstruction2.buildStageInfo.ValueOrZero().name.toString()
	}

	return buildStageName, nil
}

func vdiffBaseImages(bs1 *buildStage, bs2 *buildStage) (VDiffBaseImage, error) {
	slog.Debug("vdiff base images")

	if len(bs1.instructions) == 0 {
		return VDiffBaseImage{}, fmt.Errorf("BuildStage should have more than 0 instruction, but buildStage1 is 0")
	}
	if len(bs2.instructions) == 0 {
		return VDiffBaseImage{}, fmt.Errorf("BuildStage should have more than 0 instruction, but buildStage2 is 0")
	}

	fromInstruction1, err := bs1.getFromInstruction()
	if err != nil {
		return VDiffBaseImage{}, fmt.Errorf("failed to get FROM instruction from buildStage1: %v", err)
	}

	fromInstruction2, err := bs2.getFromInstruction()
	if err != nil {
		return VDiffBaseImage{}, fmt.Errorf("failed to get FROM instruction from buildStage2: %v", err)
	}

	imageInfo1 := fromInstruction1.imageInfo
	imageInfo2 := fromInstruction2.imageInfo

	if imageInfo1.name != imageInfo2.name {
		return VDiffBaseImage{}, fmt.Errorf("base image's name must equal but different: %s, %s", imageInfo1.name, imageInfo2.name)
	}

	return VDiffBaseImage{
		Name:      imageInfo1.name,
		Moldfile1: imageInfo1.getTagOrDigest(true),
		Moldfile2: imageInfo2.getTagOrDigest(true),
	}, nil
}

func vdiffPackages(bs1 *buildStage, bs2 *buildStage) ([]VDiffPackageInfo, error) {
	slog.Debug("vdiff packages")

	reader1 := newVDiffInstructionsReader(bs1.instructions)
	reader2 := newVDiffInstructionsReader(bs2.instructions)

	packageInfos := make([]VDiffPackageInfo, 0)
	for {
		runIstr1 := reader1.nextRunInstruction()
		runIstr2 := reader2.nextRunInstruction()

		if runIstr1 == nil && runIstr2 == nil {
			break
		}

		slog.Debug(fmt.Sprintf("runIstr1 found at index %d", reader1.currentIndex))
		slog.Debug(fmt.Sprintf("runIstr2 found at index %d", reader2.currentIndex))

		if runIstr1 == nil || runIstr2 == nil {
			return nil, fmt.Errorf("RunInstruction counts are different")
		}

		pmcs1 := runIstr1.getPackageManagerCmds()
		pmcs2 := runIstr2.getPackageManagerCmds()
		if len(pmcs1) != len(pmcs2) {
			return nil, fmt.Errorf("RunInstructions must have same number of PackageManagerCmd, but differnt: %d, %d", len(pmcs1), len(pmcs2))
		}

		for i := 0; i < len(pmcs1); i++ {
			pmc1 := pmcs1[i]
			pmc2 := pmcs2[i]

			if len(pmc1.packages) == 0 && len(pmc2.packages) == 0 {
				continue
			}

			if pmc1.mainCmd.toString() != pmc2.mainCmd.toString() {
				return nil, fmt.Errorf("PackageManager's MainCommand must equal, but different: %s, %s", pmc1.mainCmd, pmc2.mainCmd)
			}
			if pmc1.subCmd.toString() != pmc2.subCmd.toString() {
				return nil, fmt.Errorf("PackageManager's SubCommand must equal, but different: %s, %s", pmc1.subCmd, pmc2.subCmd)
			}

			slog.Debug(fmt.Sprintf("vdiff package manager command: `%s %s`", pmc1.mainCmd, pmc1.subCmd))

			packageInfos1 := pmc1.packages
			packageInfos2 := pmc2.packages
			if len(packageInfos1) != len(packageInfos2) {
				return nil, fmt.Errorf("PackageManagerCommand must have same number of packages, but differnt: %d, %d", len(packageInfos1), len(packageInfos2))
			}

			for j := 0; j < len(packageInfos1); j++ {
				vdiff, err := vdiffPackageInfos(packageInfos1[j], packageInfos2[j])
				if err != nil {
					return nil, err
				}

				packageInfos = append(packageInfos, vdiff)
			}
		}
	}

	return packageInfos, nil
}

func vdiffPackageInfos(packageInfo1 packageInfo, packageInfo2 packageInfo) (VDiffPackageInfo, error) {
	// TODO: Support other package managers
	apt1 := packageInfo1.(*aptPackageInfo)
	apt2 := packageInfo2.(*aptPackageInfo)

	if apt1.name != apt2.name {
		return VDiffPackageInfo{}, fmt.Errorf("package name must equal, but differnt: %s, %s", apt1.name, apt2.name)
	}

	return VDiffPackageInfo{
		PackageManager: "apt",
		Name:           apt1.name,
		Moldfile1:      apt1.version.ValueOrZero(),
		Moldfile2:      apt2.version.ValueOrZero(),
	}, nil
}

type vdiffInstructionsReader struct {
	currentIndex int
	instructions []instruction
}

func newVDiffInstructionsReader(instructions []instruction) *vdiffInstructionsReader {
	return &vdiffInstructionsReader{
		currentIndex: -1,
		instructions: instructions,
	}
}

func (v *vdiffInstructionsReader) nextRunInstruction() *runInstruction {
	v.currentIndex++
	for v.currentIndex < len(v.instructions) {
		istr := v.instructions[v.currentIndex]
		runIstr, ok := istr.(*runInstruction)
		if ok {
			return runIstr
		}

		v.currentIndex++
	}

	return nil
}
