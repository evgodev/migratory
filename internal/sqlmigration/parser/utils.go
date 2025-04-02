package parser

import (
	"strings"
)

const (
	sqlCommentPrefix = "--"
	commandPrefix    = "-- +migrate"
)

// endsWithSemicolon checks if the given line ends with a semicolon, ignoring SQL comments and trailing whitespace.
func endsWithSemicolon(line string) bool {
	if idx := strings.Index(line, sqlCommentPrefix); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	return len(line) > 0 && line[len(line)-1] == ';'
}

func isCommand(line string) bool {
	return strings.HasPrefix(line, commandPrefix)
}

func isSQLComment(line string) bool {
	return strings.HasPrefix(line, sqlCommentPrefix) && !isCommand(line)
}

func isEmpty(line string) bool {
	return strings.TrimSpace(line) == ""
}
