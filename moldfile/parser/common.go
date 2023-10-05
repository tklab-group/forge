package parser

import "bytes"

// space includes single space (" ") or tab ("\t")
type space struct {
	rawText string
}

func isSpace(b []byte) bool {
	return bytes.Equal(b, []byte(" ")) || bytes.Equal(b, []byte("\t"))
}

func newSpaceFromByte(b []byte) *space {
	return &space{
		rawText: string(b),
	}
}

type newlineChar struct {
	rawText string
}

func isNewlineChar(b []byte) bool {
	// TODO: Support new line characters for Windows ("\r\n")
	return bytes.Equal(b, []byte("\n")) || bytes.Equal(b, []byte("\r"))
}
func newNewlineCharFromByte(b []byte) *newlineChar {
	return &newlineChar{
		rawText: string(b),
	}
}
