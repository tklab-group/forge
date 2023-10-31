package parser

import (
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/command"
	"io"
	"log/slog"
	"strings"
)

type MoldFile interface {
	stringfy
	ToString() string
	//TextDiff() string
	//Diff() []string // TODO: Change the return value format

	BuildStageCount() int
	GetBuildStage(index int) (BuildStage, error)
}

type moldFile struct {
	buildStages []BuildStage
}

type BuildStage interface {
	implBuildStage()
	stringfy
	ToString() string
	GetFromInstruction() (FromInstruction, error)
	UpdatePackageInfos(reference packageVersions)
}

type buildStage struct {
	instructions []instruction
}

type instruction interface {
	implInstruction()
	stringfy
}

// ParseMoldFile parses a MoldFile format file (includes Dockerfile)
func ParseMoldFile(ior io.Reader) (MoldFile, error) {
	r, err := newReader(ior)
	if err != nil {
		return nil, err
	}

	instructions := make([]instruction, 0)
	for !r.Empty() {
		nextType, err := checkNextInstructionType(r)
		if err != nil {
			return nil, fmt.Errorf("faled to check next instruction type: %v", err)
		}

		switch nextType {
		case command.Run:
			slog.Debug("parsing RUN instruction")
			runInstruction, err := ParseRunInstruction(r)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as RUN: %v", err)
			}

			instructions = append(instructions, runInstruction)
		case command.From:
			slog.Debug("parsing FROM instruction")
			fromInstruction, err := ParseFromInstruction(r)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as FROM: %v", err)
			}

			instructions = append(instructions, fromInstruction)
		case "unknown":
			return nil, fmt.Errorf("unknown instruction found")
		default:
			slog.Debug(fmt.Sprintf("parsing %s instruction with ParseOtherInstruction", nextType))

			enableMultiline := false
			if nextType == command.Healthcheck {
				enableMultiline = true
			}

			otherInstruction, err := ParseOtherInstruction(r, enableMultiline)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as OtherInstruction: %v", err)
			}

			instructions = append(instructions, otherInstruction)
		}

	}

	buildStages := make([]BuildStage, 0)
	current := make([]instruction, 0)

	for _, istr := range instructions {
		_, isFromInstruction := istr.(*fromInstruction)
		if isFromInstruction && len(current) != 0 {
			currentStage := &buildStage{
				instructions: current,
			}
			buildStages = append(buildStages, currentStage)

			current = make([]instruction, 0)
		}

		current = append(current, istr)
	}

	if len(current) != 0 {
		currentStage := &buildStage{
			instructions: current,
		}
		buildStages = append(buildStages, currentStage)
	}

	parsed := &moldFile{
		buildStages: buildStages,
	}

	return parsed, nil
}

func checkNextInstructionType(r reader) (string, error) {
	builder := new(strings.Builder)

	for !r.Empty() {
		b, err := r.ReadBytes()
		if err != nil {
			return "", err
		}

		if isSpace(b) || isNewlineChar(b) {
			_, err = r.Seek(-1, io.SeekCurrent)
			if err != nil {
				return "", err
			}

			break
		}

		_, err = builder.Write(b)
		if err != nil {
			return "", err
		}
	}

	// Reset the offset of the reading buffer
	_, err := r.Seek(int64(-builder.Len()), io.SeekCurrent)
	if err != nil {
		return "", err
	}

	s := strings.ToLower(builder.String())
	_, isDockerfileInstruction := command.Commands[s]
	if isDockerfileInstruction {
		return s, nil
	}

	switch s {
	case "#":
		return "comment", nil
	case "":
		return "blank", nil
	default:
		return "unknown", nil
	}
}

func (m *moldFile) toString() string {
	return joinStringfys(m.buildStages)
}

func (m *moldFile) ToString() string {
	return m.toString()
}

func (m *moldFile) BuildStageCount() int {
	return len(m.buildStages)
}

func (m *moldFile) GetBuildStage(index int) (BuildStage, error) {
	actualCount := m.BuildStageCount()
	if index >= actualCount {
		return nil, fmt.Errorf("build stage counts is %d. Index %d is out of range", actualCount, index)
	}

	return m.buildStages[index], nil
}

func (b *buildStage) ToString() string {
	return b.toString()
}

func (b *buildStage) GetFromInstruction() (FromInstruction, error) {
	if len(b.instructions) == 0 {
		return nil, fmt.Errorf("this build stage has no instruction")
	}

	firstIstr := b.instructions[0]
	fromIstr, ok := firstIstr.(FromInstruction)
	if !ok {
		return nil, fmt.Errorf("expected first instruction is FromInstruction but %T", firstIstr)
	}

	return fromIstr, nil
}

func (b *buildStage) UpdatePackageInfos(reference packageVersions) {
	for _, istr := range b.instructions {
		runIstr, ok := istr.(RunInstruction)
		if ok {
			runIstr.UpdatePackageInfos(reference)
		}
	}

}

func (b *buildStage) toString() string {
	return joinStringfys(b.instructions)
}

func (b *buildStage) implBuildStage() {}
