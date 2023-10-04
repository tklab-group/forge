package parser

type MoldFile interface {
	ToString() string
	TextDiff() string
	Diff() []string // TODO: Change the return value format
}

type moldFile struct {
	instructionStatements []instructionStatement
}

type instructionStatement interface {
	toString() string
	textDiff() string
	diff() string // TODO: Change the return value format
}

type instructionType string

const (
	instructionTypeFrom  instructionType = "FROM"
	instructionTypeRun   instructionType = "RUN"
	instructionTypeOther instructionType = "other" // This parser doesn't parse instruction statements except FROM and RUN
)
