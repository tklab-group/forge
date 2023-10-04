package parser

import "io"

// ParseMoldFile parses a MoldFile format file (includes Dockerfile)
func ParseMoldFile(r io.Reader) (MoldFile, error) {
	return parseMoldFile(r)
}

func parseMoldFile(r io.Reader) (MoldFile, error) {
	// TODO: Implement
	return nil, nil
}
