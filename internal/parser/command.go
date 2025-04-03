package parser

import "strings"

type (
	commandBody   = string
	commandOption = string
)

const (
	commandUp             commandBody   = "up"
	commandDown           commandBody   = "down"
	commandStatementBegin commandBody   = "statement_begin"
	commandStatementEnd   commandBody   = "statement_end"
	optionNoTransaction   commandOption = "no_transaction"
)

type command struct {
	body    commandBody
	options []commandOption
}

func newCommand(line string) (*command, error) {
	fields := strings.Fields(strings.TrimPrefix(line, commandPrefix))
	if len(fields) == 0 {
		return nil, ErrIncompleteCommand
	}

	return &command{
		body:    fields[0],
		options: fields[1:],
	}, nil
}

func (c *command) hasOption(option string) bool {
	for _, o := range c.options {
		if o == option {
			return true
		}
	}

	return false
}
