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
	buildStageInfo optional.Of[*buildStageInfo]
	platformFlag   optional.Of[*platformFlag]
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

type buildStageInfo struct {
	elements []buildStageInfoElement
	asString *buildStageInfoAsString
	name     *buildStageInfoName
}

type buildStageInfoElement interface {
	implBuildStageInfoElement()
}

type buildStageInfoAsString struct {
	rawText string
}

type buildStageInfoName struct {
	rawText string
}

type platformFlag struct {
	rawText string
}

func ParseFromInstruction(r io.Reader) (FromInstruction, error) {
	instruction := &fromInstruction{
		elements:       make([]fromInstructionElement, 0),
		imageInfo:      nil,
		buildStageInfo: optional.Of[*buildStageInfo]{},
		platformFlag:   optional.Of[*platformFlag]{},
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

	if strings.HasPrefix(s, "--platform=") {
		element := &platformFlag{
			rawText: s,
		}

		f.appendElement(element)
		f.platformFlag = optional.NewWithValue(element)

		buffer.Reset()
		appendCurrentByte()

		return nil
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

	// Parse `AS {buildStageName}`
	if strings.ToLower(s) == "as" && !f.buildStageInfo.HasValue() {
		err := f.parseBuildStageInfo(scanner, buffer, currentByte, s)
		if err != nil {
			return fmt.Errorf("failed to parse buildStageInfo: %v", err)
		}
		return nil
	}

	return fmt.Errorf("unexpected format")
}

func (f *fromInstruction) parseBuildStageInfo(scanner *bufio.Scanner, buffer *bytes.Buffer, currentByte []byte, s string) error {
	if !isSpace(currentByte) {
		return fmt.Errorf("unexpected token: %v", currentByte)
	}

	info := &buildStageInfo{
		elements: make([]buildStageInfoElement, 0),
		asString: nil,
		name:     nil,
	}

	asString := &buildStageInfoAsString{rawText: s}
	info.elements = append(info.elements, asString)
	info.asString = asString

	info.elements = append(info.elements, newSpaceFromByte(currentByte))

	buildStageBuff := new(bytes.Buffer)

	// Read all spaces between "AS" and the stage name
	for scanner.Scan() {
		b := scanner.Bytes()

		if isNewlineChar(b) {
			return fmt.Errorf("unexpected newline code")
		}

		if isSpace(b) {
			info.elements = append(info.elements, newNewlineCharFromByte(b))
		} else {
			_, err := buildStageBuff.Write(b)
			if err != nil {
				return fmt.Errorf("failed to write to buffer: %v", err)
			}
			break
		}
	}

	var lastElement buildStageInfoElement
	for scanner.Scan() {
		b := scanner.Bytes()

		if isNewlineChar(b) {
			lastElement = newNewlineCharFromByte(b)
			break
		}
		if isSpace(b) {
			lastElement = newSpaceFromByte(b)
			break
		}

		_, err := buildStageBuff.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write to buffer: %v", err)
		}
	}

	name := &buildStageInfoName{
		rawText: buildStageBuff.String(),
	}
	info.elements = append(info.elements, name)
	info.name = name

	info.elements = append(info.elements, lastElement)

	f.appendElement(info)
	f.buildStageInfo = optional.NewWithValue(info)

	buffer.Reset()
	return nil
}

func (f *fromInstruction) appendElement(element fromInstructionElement) {
	f.elements = append(f.elements, element)
}

func (f *fromString) implFromInstructionElement()     {}
func (i *imageInfo) implFromInstructionElement()      {}
func (b *buildStageInfo) implFromInstructionElement() {}
func (p *platformFlag) implFromInstructionElement()   {}
func (s *space) implFromInstructionElement()          {}
func (n *newlineChar) implFromInstructionElement()    {}

func (b *buildStageInfoAsString) implBuildStageInfoElement() {}
func (b *buildStageInfoName) implBuildStageInfoElement()     {}
func (s *space) implBuildStageInfoElement()                  {}
func (n *newlineChar) implBuildStageInfoElement()            {}
