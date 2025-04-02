package parser

// parsingState describes state of a migration parsing.
type parsingState struct {
	direction migrationDirection
	statement statementState
}

type (
	migrationDirection int
	statementState     int
)

const (
	directionNone migrationDirection = iota
	directionUp
	directionDown
)

const (
	statementNone statementState = iota
	statementStarted
	statementEnded
)

func newParsingState() *parsingState {
	return &parsingState{
		direction: directionNone,
		statement: statementNone,
	}
}

func (s *parsingState) setDirectionUp() {
	s.direction = directionUp
}

func (s *parsingState) setDirectionDown() {
	s.direction = directionDown
}

func (s *parsingState) setStatementStarted() {
	s.statement = statementStarted
}

func (s *parsingState) setStatementEnded() error {
	if s.statement != statementStarted {
		return ErrStatementNotStarted
	}
	s.statement = statementEnded
	return nil
}

func (s *parsingState) setStatementNone() {
	s.statement = statementNone
}

func (s *parsingState) isStatementNone() bool {
	return s.statement == statementNone
}

func (s *parsingState) isStatementEnded() bool {
	return s.statement == statementEnded
}
