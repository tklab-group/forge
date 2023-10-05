package parser

import "bytes"

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
