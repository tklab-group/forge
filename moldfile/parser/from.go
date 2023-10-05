package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tklab-group/forge/util/optional"
)

type FromInstruction interface {
}

type fromInstruction struct {
	elements       []fromInstructionElement
	fromString     *fromString
	imageInfo      *imageInfo
	buildStageName optional.Of[*buildStageName]
	platformArg    optional.Of[*platformArg]
}

type fromInstructionElement interface {
	implFromInstructionElement()
}

type fromString struct {
	rawText string
}

type imageInfo struct {
	name   string
	tag    optional.Of[string]
	digest optional.Of[string]
}

type buildStageName struct {
}

type platformArg struct{}

func ParseFromInstruction(r io.Reader) (FromInstruction, error) {
	instruction := &fromInstruction{
		elements:       make([]fromInstructionElement, 0),
		imageInfo:      nil,
		buildStageName: optional.Of[*buildStageName]{},
		platformArg:    optional.Of[*platformArg]{},
	}

	buffer := new(bytes.Buffer)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		b := scanner.Bytes()
		if isSpace(b) {
			err := instruction.treatFromInstructionElement(scanner, buffer, b)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
			}
			continue
		}

		if isNewlineChar(b) {
			err := instruction.treatFromInstructionElement(scanner, buffer, b)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
			}

			// FROM instruction ends with newline
			break
		}

		_, err := buffer.Write(b)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %v", err)
		}
	}

	err := instruction.treatFromInstructionElement(scanner, buffer, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
	}

	return instruction, nil
}

func (f *fromInstruction) treatFromInstructionElement(scanner *bufio.Scanner, buffer *bytes.Buffer, currentByte []byte) error {
	appendCurrentByte := func() {
		if isSpace(currentByte) {
			f.appendElement(newSpaceFromByte(currentByte))
		}
		if isNewlineChar(currentByte) {
			f.appendElement(newNewlineCharFromByte(currentByte))
		}
	}

	if buffer.Len() == 0 {
		appendCurrentByte()
		return nil
	}

	s := buffer.String()

	if strings.ToLower(s) == "from" && f.fromString == nil {
		element := &fromString{
			rawText: s,
		}
		f.appendElement(element)
		f.fromString = element

		buffer.Reset()
		appendCurrentByte()

		return nil
	}

	if strings.HasPrefix(buffer.String(), "--platform=") {
		// TODO
	}

	if f.imageInfo == nil {
		image := &imageInfo{
			name:   "",
			tag:    optional.Of[string]{},
			digest: optional.Of[string]{},
		}

		switch {
		case strings.Contains(s, "@"):
			split := strings.Split(s, "@")
			image.name = split[0]
			image.digest = optional.NewWithValue(split[1])
		case strings.Contains(s, ":"):
			split := strings.Split(s, ":")
			image.name = split[0]
			image.tag = optional.NewWithValue(split[1])
		default:
			image.name = s
		}

		f.appendElement(image)
		f.imageInfo = image

		buffer.Reset()
		appendCurrentByte()

		return nil
	}

	// TODO: buildStageName

	return fmt.Errorf("unexpected format")
}

func (f *fromInstruction) appendElement(element fromInstructionElement) {
	f.elements = append(f.elements, element)
}

func (f *fromString) implFromInstructionElement() {}

func (i *imageInfo) implFromInstructionElement() {}

func (s *space) implFromInstructionElement() {}

func (n *newlineChar) implFromInstructionElement() {}
