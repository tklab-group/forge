package parser

type runStatement struct {
	cmds []cmdInstruction
}

type cmdInstruction interface{} // TODO: Rename

type cmdSeparator struct{}

type cmd struct{}

type cmdNode interface{}

type cmdName struct{}

type cmdOption struct{}

type cmdArg struct{}

type cmdInnerComment struct{}

type cmdInnerBackSlash struct{}

// cmdInnerSpace includes space, \r and \t
type cmdInnerSpace struct{}
