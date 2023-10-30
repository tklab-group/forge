package parser

import (
	"bytes"
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/command"
	"io"
	"log/slog"
	"strings"
)

type MoldFile interface {
	stringfy
	//ToString() string
	//TextDiff() string
	//Diff() []string // TODO: Change the return value format
}

type moldFile struct {
	buildStages []*buildStage
}

type buildStage struct {
	instructions []instruction
}

type instruction interface {
	implInstruction()
	stringfy
}

// ParseMoldFile parses a MoldFile format file (includes Dockerfile)
func ParseMoldFile(r io.Reader) (MoldFile, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read all: %v", err)
	}

	br := bytes.NewReader(b)

	instructions := make([]instruction, 0)
	for br.Len() != 0 {
		nextType, err := checkNextInstructionType(br)
		if err != nil {
			return nil, fmt.Errorf("faled to check next instruction type: %v", err)
		}

		switch nextType {
		case command.Run:
			slog.Debug("parsing RUN instruction")
			runInstruction, err := ParseRunInstruction(br)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as RUN: %v", err)
			}

			instructions = append(instructions, runInstruction)
		case command.From:
			slog.Debug("parsing FROM instruction")
			fromInstruction, err := ParseFromInstruction(br)
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

			otherInstruction, err := ParseOtherInstruction(br, enableMultiline)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as OtherInstruction: %v", err)
			}

			instructions = append(instructions, otherInstruction)
		}

	}

	// TODO: Separate to multiple stages with FROM instruction
	parsed := &moldFile{
		buildStages: []*buildStage{
			{
				instructions: instructions,
			},
		},
	}

	return parsed, nil
}

func checkNextInstructionType(br *bytes.Reader) (string, error) {
	builder := new(strings.Builder)

	for br.Len() != 0 {
		bUnit, err := br.ReadByte()
		if err != nil {
			return "", err
		}
		b := []byte{bUnit}

		if isSpace(b) || isNewlineChar(b) {
			err = br.UnreadByte()
			if err != nil {
				return "", err
			}

			break
		}

		err = builder.WriteByte(bUnit)
		if err != nil {
			return "", err
		}
	}

	// Reset the offset of the reading buffer
	_, err := br.Seek(int64(-builder.Len()), io.SeekCurrent)
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

func (b *buildStage) toString() string {
	return joinStringfys(b.instructions)
}
