package parser

import (
	"bytes"
	"strings"
)

// space includes single space (" ") or tab ("\t")
type space struct {
	rawTextContainer
}

func isSpace(b []byte) bool {
	return bytes.Equal(b, []byte(" ")) || bytes.Equal(b, []byte("\t"))
}

func newSpaceFromByte(b []byte) *space {
	return &space{
		newRawTextContainer(string(b)),
	}
}

type newlineChar struct {
	rawTextContainer
}

func isNewlineChar(b []byte) bool {
	// TODO: Support new line characters for Windows ("\r\n")
	return bytes.Equal(b, []byte("\n")) || bytes.Equal(b, []byte("\r"))
}
func newNewlineCharFromByte(b []byte) *newlineChar {
	return &newlineChar{
		newRawTextContainer(string(b)),
	}
}

type backslash struct {
	rawTextContainer
}

func isBackslash(b []byte) bool {
	return bytes.Equal(b, []byte("\\"))
}

func isBackslashString(s string) bool {
	return s == "\\"
}

func newBackslashFromByte(b []byte) *backslash {
	return &backslash{
		newRawTextContainer(string(b)),
	}
}

func newBackslash(s string) *backslash {
	return &backslash{
		newRawTextContainer(s),
	}
}

type comment struct {
	rawTextContainer
}

func isCommentSharp(b []byte) bool {
	return bytes.Equal(b, []byte("#"))
}

func newComment(s string) *comment {
	return &comment{
		newRawTextContainer(s),
	}
}

type rawTextContainer struct {
	rawText string
}

func newRawTextContainer(s string) rawTextContainer {
	return rawTextContainer{
		rawText: s,
	}
}

func (r *rawTextContainer) toString() string {
	return r.rawText
}

type stringfy interface {
	toString() string
}

func joinStringfys[T stringfy](stringfys []T) string {
	list := make([]string, len(stringfys))
	for i, element := range stringfys {
		list[i] = element.toString()
	}

	return strings.Join(list, "")
}
