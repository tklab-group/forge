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
	stringfy
	ToString() string
}

type runInstruction struct {
	elements  []runInstructionElement
	runString *runString
}

type runInstructionElement interface {
	implRunInstructionElement()
	stringfy
}

type runString struct {
	rawTextContainer
}

type packageManagerCmd struct {
	elements       []packageManagerCmdElement
	mainCmd        *packageManagerMainCmd
	mainCmdOptions []*packageManagerOption
	subCmd         *packageManagerSubCmd
	subCmdOptions  []*packageManagerOption
	packages       []packageInfo
}

type packageManagerCmdElement interface {
	implPackageManagerCmdElement()
	stringfy
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

type packageManagerArg struct {
	packageInfo packageInfo
}

type packageInfo interface {
	implPackageInfo()
	stringfy
}

var supportedPackageManagerCmd = []string{"apt", "apt-get"}
var packageInfoParseFuncs = map[string]func(s string) packageInfo{
	"apt":     parseAptPackageInfo,
	"apt-get": parseAptPackageInfo,
}

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
			continue
		}

		if isSpace(b) {
			if buffer.Len() == 0 {
				instruction.appendElement(newSpaceFromByte(b))
			} else {
				if slices.Contains(supportedPackageManagerCmd, buffer.String()) {
					err := instruction.parsePackageManagerCmd(scanner, buffer, b)
					if err != nil {
						return nil, fmt.Errorf("failed to parse as a package manager command: %v", err)
					}
				} else {
					err := instruction.parseOtherCmd(scanner, buffer, b)
					if err != nil {
						return nil, fmt.Errorf("failed to pares as an other command: %v", err)
					}
				}
			}
			continue
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

func (r *runInstruction) parsePackageManagerCmd(scanner *bufio.Scanner, buffer *bytes.Buffer, currentByte []byte) error {
	if !(slices.Contains(supportedPackageManagerCmd, buffer.String()) && isSpace(currentByte)) {
		return fmt.Errorf("unexpected input for parsePackageManagerCmd: `%s%s`", buffer.String(), string(currentByte))
	}

	packageInfoParser, ok := packageInfoParseFuncs[buffer.String()]
	if !ok {
		return fmt.Errorf("no function for parse its args: `%s`", buffer.String())
	}

	managerCmd := &packageManagerCmd{
		elements:       make([]packageManagerCmdElement, 0),
		mainCmd:        nil,
		mainCmdOptions: make([]*packageManagerOption, 0),
		subCmd:         nil,
		subCmdOptions:  make([]*packageManagerOption, 0),
		packages:       make([]packageInfo, 0),
	}
	r.appendElement(managerCmd)

	mainCmd := &packageManagerMainCmd{newRawTextContainer(buffer.String())}
	buffer.Reset()

	managerCmd.appendElement(mainCmd)
	managerCmd.mainCmd = mainCmd
	managerCmd.appendElement(newSpaceFromByte(currentByte))

	var noStrInCurrentLine bool = false // The first line starts with a package manager command
	for scanner.Scan() {
		b := scanner.Bytes()
		if isSpace(b) {
			// TODO: Handle `&&`

			if buffer.Len() != 0 {
				managerCmd.parseElement(buffer.String(), packageInfoParser)
			}

			managerCmd.appendElement(newSpaceFromByte(b))
			buffer.Reset()
			continue
		}

		if isNewlineChar(b) {
			if isBackslashString(buffer.String()) {
				managerCmd.appendElement(newBackslash(buffer.String()))
				managerCmd.appendElement(newNewlineCharFromByte(b))
				buffer.Reset()

				noStrInCurrentLine = true
				continue
			}

			managerCmd.parseElement(buffer.String(), packageInfoParser)
			managerCmd.appendElement(newNewlineCharFromByte(b))
			buffer.Reset()
			return nil // RUN Instruction must end here
		}

		if isCommentSharp(b) && noStrInCurrentLine {
			comment, err := parseCommentLine(scanner, b)
			if err != nil {
				return fmt.Errorf("failed to parse as a comment line: %v", err)
			}

			managerCmd.appendElement(comment)
			continue
		}

		_, err := buffer.Write(b)
		if err != nil {
			return err
		}
		noStrInCurrentLine = false
	}

	// File ends with the RUN instruction
	if buffer.Len() != 0 {
		managerCmd.parseElement(buffer.String(), packageInfoParser)
		buffer.Reset()
	}

	return nil
}

func (pmc *packageManagerCmd) parseElement(s string, packageInfoParser func(s string) packageInfo) {
	if strings.HasPrefix(s, "-") {
		option := &packageManagerOption{rawTextContainer{s}}
		pmc.appendElement(option)

		if pmc.subCmd == nil {
			pmc.mainCmdOptions = append(pmc.mainCmdOptions, option)
		} else {
			pmc.subCmdOptions = append(pmc.subCmdOptions, option)
		}
		return
	}

	if pmc.subCmd == nil {
		subCmd := &packageManagerSubCmd{newRawTextContainer(s)}
		pmc.appendElement(subCmd)
		pmc.subCmd = subCmd

		return
	}

	packageInfo := packageInfoParser(s)
	arg := &packageManagerArg{packageInfo: packageInfo}
	pmc.appendElement(arg)
	pmc.packages = append(pmc.packages, packageInfo)
}

func (r *runInstruction) parseOtherCmd(scanner *bufio.Scanner, buffer *bytes.Buffer, currentByte []byte) error {
	builder := new(strings.Builder)

	_, err := builder.WriteString(buffer.String())
	if err != nil {
		return err
	}
	buffer.Reset()

	_, err = builder.Write(currentByte)
	if err != nil {
		return err
	}

	var noStrInCurrentLine bool = false // The first line starts with a command
	for scanner.Scan() {
		b := scanner.Bytes()

		if isSpace(b) {
			// TODO: Handle `&&`

			if buffer.Len() != 0 {
				_, err = builder.WriteString(buffer.String())
				if err != nil {
					return err
				}

				buffer.Reset()
			}

			_, err = builder.Write(b)
			if err != nil {
				return err
			}

			continue
		}

		if isNewlineChar(b) {
			instructionMustEnd := !isBackslashString(buffer.String())

			_, err = builder.WriteString(buffer.String())
			if err != nil {
				return err
			}

			buffer.Reset()

			_, err = builder.Write(b)
			if err != nil {
				return err
			}

			noStrInCurrentLine = true

			if instructionMustEnd {
				break
			} else {
				continue
			}
		}

		if isCommentSharp(b) && noStrInCurrentLine {
			comment, err := parseCommentLine(scanner, b)
			if err != nil {
				return fmt.Errorf("failed to parse as a comment in a command: %v", err)
			}

			// In otherCmd, comment is also parsed as just a string element in a command
			_, err = builder.WriteString(comment.toString())
			if err != nil {
				return err
			}
			continue
		}

		_, err = buffer.Write(b)
		if err != nil {
			return err
		}
		noStrInCurrentLine = false
	}

	// File ends with the RUN instruction
	if buffer.Len() != 0 {
		_, err = builder.WriteString(buffer.String())
		if err != nil {
			return err
		}

		buffer.Reset()
	}

	cmd := &otherCmd{newRawTextContainer(builder.String())}
	r.appendElement(cmd)

	return nil
}

func (r *runInstruction) appendElement(element runInstructionElement) {
	r.elements = append(r.elements, element)
}

func (p *packageManagerCmd) appendElement(element packageManagerCmdElement) {
	p.elements = append(p.elements, element)
}

func (r *runInstruction) toString() string {
	return joinStringfys(r.elements)
}

func (r *runInstruction) ToString() string {
	return r.toString()
}

func (pmc *packageManagerCmd) toString() string {
	return joinStringfys(pmc.elements)
}

func (p *packageManagerArg) toString() string {
	return p.packageInfo.toString()
}

func (r *runString) implRunInstructionElement()         {}
func (p *packageManagerCmd) implRunInstructionElement() {}
func (o *otherCmd) implRunInstructionElement()          {}
func (c *comment) implRunInstructionElement()           {}
func (s *space) implRunInstructionElement()             {}

func (p *packageManagerMainCmd) implPackageManagerCmdElement() {}
func (p *packageManagerOption) implPackageManagerCmdElement()  {}
func (p *packageManagerSubCmd) implPackageManagerCmdElement()  {}
func (p *packageManagerArg) implPackageManagerCmdElement()     {}
func (s *space) implPackageManagerCmdElement()                 {}
func (b *backslash) implPackageManagerCmdElement()             {}
func (n *newlineChar) implPackageManagerCmdElement()           {}
func (c *comment) implPackageManagerCmdElement()               {}

func (a *aptPackageInfo) implPackageInfo() {}
