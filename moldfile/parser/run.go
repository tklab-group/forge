package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
)

type RunInstruction interface {
}

type runInstruction struct {
	elements  []runInstructionElement
	runString *runString
}

type runInstructionElement interface {
	implRunInstructionElement()
}

type runString struct {
	rawTextContainer
}

type packageManagerCmd[T packageInfo] struct {
	elements []packageManagerCmdElement
	mainCmd  *packageManagerMainCmd
	options  []*packageManagerOption
	subCmd   *packageManagerSubCmd
	args     []T
}

type packageManagerCmdElement interface {
	implPackageManagerCmdElement()
}

type packageManagerMainCmd struct {
	rawTextContainer
}

type packageManagerOption struct {
	rawTextContainer
}

type packageManagerSubCmd struct {
	rawTextContainer
}

type packageInfo interface {
	implPackageInfo()
	stringfy
}

var supportedPackageManagerCmd = []string{"apt", "apt-get"}

type otherCmd struct {
	rawTextContainer
}

func ParseRunInstruction(r io.Reader) (RunInstruction, error) {
	instruction := &runInstruction{
		elements:  make([]runInstructionElement, 0),
		runString: nil,
	}

	buffer := new(bytes.Buffer)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)

	// Parse "RUN"
	for scanner.Scan() {
		b := scanner.Bytes()

		if isSpace(b) {
			s := buffer.String()
			if strings.ToLower(s) != "run" {
				return nil, fmt.Errorf("RunInstruction must start with `RUN`. got: %s", s)
			}

			element := &runString{
				newRawTextContainer(s),
			}
			instruction.appendElement(element)
			instruction.runString = element

			instruction.appendElement(newSpaceFromByte(b))

			buffer.Reset()

			break
		}

		_, err := buffer.Write(b)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %v", err)
		}
	}

	// Parse commands in RUN statements
	var noStrInCurrentLine bool = false // The first line starts with RUN
	for scanner.Scan() {
		b := scanner.Bytes()

		if isCommentSharp(b) && noStrInCurrentLine {
			comment, err := parseCommentLine(scanner, b)
			if err != nil {
				return nil, fmt.Errorf("faild to parse comment line: %v", err)
			}

			instruction.appendElement(comment)
		}

		if isSpace(b) {
			if buffer.Len() == 0 {
				instruction.appendElement(newSpaceFromByte(b))
			} else {
				// TODO: parse command
				if slices.Contains(supportedPackageManagerCmd, buffer.String()) {

				} else {
					
				}
			}
		}

		_, err := buffer.Write(b)
		if err != nil {
			return nil, err
		}
		noStrInCurrentLine = false
	}

	return instruction, nil
}

func parseCommentLine(scanner *bufio.Scanner, prevByte []byte) (*comment, error) {
	commentBuffer := new(bytes.Buffer)
	_, err := commentBuffer.Write(prevByte)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		b := scanner.Bytes()
		if isNewlineChar(b) {
			_, err = commentBuffer.Write(b)
			if err != nil {
				return nil, err
			}

			return newComment(commentBuffer.String()), nil
		}

		_, err := commentBuffer.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("comment must end with newline. actual: %s", commentBuffer.String())
}

func (r *runInstruction) appendElement(element runInstructionElement) {
	r.elements = append(r.elements, element)
}

func (r *runString) implRunInstructionElement() {}
func (c *comment) implRunInstructionElement()   {}
func (s *space) implRunInstructionElement()     {}
