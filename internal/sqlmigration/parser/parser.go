// Package parser provides functionality to parse SQL migrations, divide it to up and down SQL statements.
// See testdata for examples.
package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNoSemicolon         = errors.New("statement must be ended by a semicolon")
	ErrIncompleteCommand   = errors.New("incomplete migration command")
	ErrUnknownCommand      = errors.New("unknown migration command after prefix")
	ErrStatementNotEnded   = errors.New("statement was started but not ended")
	ErrStatementNotStarted = errors.New("statement was ended but not started")
	ErrNoUpDownCommands    = errors.New("no Up and Down commands found during parsing")
)

// ParsedMigration describes up and down SQL statements.
type ParsedMigration struct {
	UpStatements   []string
	DownStatements []string

	DisableTransactionUp   bool
	DisableTransactionDown bool
}

// ParseMigration parses SQL migration scripts into up and down statements, handling specific commands and identifiers.
// It returns a ParsedMigration with the parsed statements and transactions configuration or an error on failure.
func ParseMigration(r io.Reader) (*ParsedMigration, error) {
	p := newParser(r)
	if err := p.parseLines(); err != nil {
		return nil, fmt.Errorf("failed to parse lines: %w", err)
	}

	result, err := p.getResult()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize parsing: %w", err)
	}

	return result, nil
}

type parser struct {
	scanner *bufio.Scanner
	buffer  *bytes.Buffer
	state   *parsingState
	result  *ParsedMigration
}

func newParser(r io.Reader) *parser {
	return &parser{
		scanner: bufio.NewScanner(r),
		buffer:  &bytes.Buffer{},
		state:   newParsingState(),
		result:  &ParsedMigration{},
	}
}

func (p *parser) parseLines() error {
	for p.scanner.Scan() {
		line := p.scanner.Text()

		if isEmpty(line) || isSQLComment(line) {
			continue
		}

		if isCommand(line) {
			if err := p.handleCommand(line); err != nil {
				return err
			}
		} else {
			if err := p.writeToBuffer(line); err != nil {
				return err
			}
		}

		if p.state.isStatementEnded() || (p.state.isStatementNone() && endsWithSemicolon(line)) {
			if p.state.direction == directionUp {
				p.result.UpStatements = append(p.result.UpStatements, p.buffer.String())
			} else {
				p.result.DownStatements = append(p.result.DownStatements, p.buffer.String())
			}

			p.state.setStatementNone()
			p.buffer.Reset()
		}
	}

	if err := p.scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan strings: %w", err)
	}

	return nil
}

func (p *parser) handleCommand(line string) error {
	cmd, err := newCommand(line)
	if err != nil {
		return err
	}

	switch cmd.body {
	case commandUp:
		if p.buffer.Len() > 0 {
			return ErrNoSemicolon
		}
		p.state.setDirectionUp()
		if cmd.hasOption(optionNoTransaction) {
			p.result.DisableTransactionUp = true
		}

	case commandDown:
		if p.buffer.Len() > 0 {
			return ErrNoSemicolon
		}
		p.state.setDirectionDown()
		if cmd.hasOption(optionNoTransaction) {
			p.result.DisableTransactionDown = true
		}

	case commandStatementBegin:
		p.state.setStatementStarted()

	case commandStatementEnd:
		return p.state.setStatementEnded()

	default:
		return ErrUnknownCommand
	}

	return nil
}

func (p *parser) getResult() (*ParsedMigration, error) {
	if p.state.statement == statementStarted {
		return nil, ErrStatementNotEnded
	}
	if p.state.direction == directionNone {
		return nil, ErrNoUpDownCommands
	}
	if p.buffer.Len() > 0 {
		return nil, ErrNoSemicolon
	}
	return p.result, nil
}

func (p *parser) writeToBuffer(line string) error {
	if _, err := p.buffer.WriteString(line + "\n"); err != nil {
		return fmt.Errorf("failed to write string to buffer: %w", err)
	}
	return nil
}
