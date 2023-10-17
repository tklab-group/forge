package parser

import (
	"bytes"
	"github.com/moby/buildkit/frontend/dockerfile/command"
	"io"
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
	// TODO

	// TODO: Separate to multiple stages with FROM instruction

	return nil, nil
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
