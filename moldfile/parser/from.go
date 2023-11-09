package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tklab-group/forge/util/optional"
)

type FromInstruction interface {
	implInstruction()
	implFromInstruction()
	stringfy
	ToString() string
	ImageInfoString() string
	UpdateImageInfo(digest string)
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
	stringfy
}

type fromString struct {
	rawTextContainer
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
	stringfy
}

type buildStageInfoAsString struct {
	rawTextContainer
}

type buildStageInfoName struct {
	rawTextContainer
}

type platformFlag struct {
	rawTextContainer
}

func ParseFromInstruction(r reader) (FromInstruction, error) {
	instruction := &fromInstruction{
		elements:       make([]fromInstructionElement, 0),
		imageInfo:      nil,
		buildStageInfo: optional.Of[*buildStageInfo]{},
		platformFlag:   optional.Of[*platformFlag]{},
	}

	buffer := new(bytes.Buffer)

	for !r.Empty() {
		b, err := r.ReadBytes()
		if err != nil {
			return nil, err
		}

		if isSpace(b) {
			err := instruction.treatFromInstructionElement(r, buffer, b)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
			}
			continue
		}

		if isNewlineChar(b) {
			err := instruction.treatFromInstructionElement(r, buffer, b)
			if err != nil {
				return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
			}

			// FROM instruction ends with newline
			break
		}

		_, err = buffer.Write(b)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer: %v", err)
		}
	}

	err := instruction.treatFromInstructionElement(r, buffer, []byte(""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse as fromInstructionElement: %v", err)
	}

	return instruction, nil
}

func (f *fromInstruction) treatFromInstructionElement(r reader, buffer *bytes.Buffer, currentByte []byte) error {
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
			newRawTextContainer(s),
		}
		f.appendElement(element)
		f.fromString = element

		buffer.Reset()
		appendCurrentByte()

		return nil
	}

	if strings.HasPrefix(s, "--platform=") {
		element := &platformFlag{
			newRawTextContainer(s),
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
		err := f.parseBuildStageInfo(r, buffer, currentByte, s)
		if err != nil {
			return fmt.Errorf("failed to parse buildStageInfo: %v", err)
		}
		return nil
	}

	return fmt.Errorf("unexpected format")
}

func (f *fromInstruction) parseBuildStageInfo(r reader, buffer *bytes.Buffer, currentByte []byte, s string) error {
	if !isSpace(currentByte) {
		return fmt.Errorf("unexpected token: %v", currentByte)
	}

	info := &buildStageInfo{
		elements: make([]buildStageInfoElement, 0),
		asString: nil,
		name:     nil,
	}

	asString := &buildStageInfoAsString{
		newRawTextContainer(s),
	}
	info.elements = append(info.elements, asString)
	info.asString = asString

	info.elements = append(info.elements, newSpaceFromByte(currentByte))

	buildStageBuff := new(bytes.Buffer)

	// Read all spaces between "AS" and the stage name
	for !r.Empty() {
		b, err := r.ReadBytes()
		if err != nil {
			return err
		}

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
	for !r.Empty() {
		b, err := r.ReadBytes()
		if err != nil {
			return err
		}

		if isNewlineChar(b) {
			lastElement = newNewlineCharFromByte(b)
			break
		}
		if isSpace(b) {
			lastElement = newSpaceFromByte(b)
			break
		}

		_, err = buildStageBuff.Write(b)
		if err != nil {
			return fmt.Errorf("failed to write to buffer: %v", err)
		}
	}

	name := &buildStageInfoName{
		newRawTextContainer(buildStageBuff.String()),
	}
	info.elements = append(info.elements, name)
	info.name = name

	info.elements = append(info.elements, lastElement)

	f.appendElement(info)
	f.buildStageInfo = optional.NewWithValue(info)

	buffer.Reset()
	return nil
}

func (f *fromInstruction) ImageInfoString() string {
	return f.imageInfo.toString()
}

func (f *fromInstruction) UpdateImageInfo(digest string) {
	f.imageInfo.updateWithDigest(digest)
}

func (f *fromInstruction) appendElement(element fromInstructionElement) {
	f.elements = append(f.elements, element)
}

func (f *fromInstruction) ToString() string {
	return f.toString()
}

func (f *fromInstruction) toString() string {
	return joinStringfys(f.elements)
}

func (i *imageInfo) toString() string {
	if i.tag.HasValue() || i.digest.HasValue() {
		return i.name + i.getTagOrDigest(true)
	}

	return i.name
}

func (i *imageInfo) updateWithDigest(digest string) {
	if i.tag.HasValue() {
		i.tag = optional.Of[string]{}
	}

	i.digest = optional.NewWithValue(digest)
}

func (i *imageInfo) getTagOrDigest(setPrefix bool) string {
	if i.tag.HasValue() {
		if setPrefix {
			return fmt.Sprintf(":%s", i.tag.ValueOrZero())
		} else {
			return i.tag.ValueOrZero()
		}
	} else {
		if setPrefix {
			return fmt.Sprintf("@%s", i.digest.ValueOrZero())
		} else {
			return i.digest.ValueOrZero()
		}
	}
}

func (b *buildStageInfo) toString() string {
	return joinStringfys(b.elements)
}

func (f *fromInstruction) implInstruction()     {}
func (f *fromInstruction) implFromInstruction() {}

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
