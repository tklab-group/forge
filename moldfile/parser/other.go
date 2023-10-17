package parser

import (
	"bufio"
	"io"
	"strings"
)

// OtherInstruction holds all instructions exclude `FROM` and `RUN`.
// It also holds a comment line and a blank line.
//
// TODO: Consider to separate to other interface for holding a comment and a blank line.
type OtherInstruction interface {
	implOtherInstruction()
	stringfy
	ToString() string
}

type otherInstruction struct {
	rawTextContainer
}

func ParseOtherInstruction(r io.Reader, enableMultiline bool) (OtherInstruction, error) {
	builder := new(strings.Builder)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		b := scanner.Bytes()
		if isNewlineChar(b) {
			if enableMultiline && strings.HasSuffix(builder.String(), " \\") {
				_, err := builder.Write(b)
				if err != nil {
					return nil, err
				}

				continue
			}

			_, err := builder.Write(b)
			if err != nil {
				return nil, err
			}

			break // Instruction must end here
		}

		_, err := builder.Write(b)
		if err != nil {
			return nil, err
		}
	}

	instruction := &otherInstruction{
		newRawTextContainer(builder.String()),
	}

	return instruction, nil
}

func (o *otherInstruction) ToString() string {
	return o.toString()
}

func (o *otherInstruction) implOtherInstruction() {}
