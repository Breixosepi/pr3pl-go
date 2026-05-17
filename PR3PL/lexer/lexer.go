package lexer

import (
	"fmt"
	"pr3pl/token"
	"strings"
	"unicode/utf8"
)

const eof = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	tokens chan token.Token
	state  stateFn
}

func New(name, input string) *Lexer {

	l := &Lexer{
		name:   name,
		input:  input,
		state:  lexTopLevel,
		tokens: make(chan token.Token, 2),
	}
	return l
}

func (l *Lexer) getLineAndCol(pos int) (int, int) {

	text := l.input[:pos]
	line := strings.Count(text, "\n") + 1
	lastNL := strings.LastIndex(text, "\n")
	col := pos - lastNL
	return line, col
}

func (l *Lexer) NextToken() token.Token {

	for {
		select {
		case tok := <-l.tokens:
			return tok
		default:
			if l.state == nil {
				line, col := l.getLineAndCol(l.pos)
				return token.Token{Type: token.EOF, Literal: "", Line: line, Column: col}
			}
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) emit(t token.TokenType) {

	line, col := l.getLineAndCol(l.start)

	l.tokens <- token.Token{
		Type:    t,
		Literal: l.input[l.start:l.pos],
		Line:    line,
		Column:  col,
	}
	l.start = l.pos
}

func (l *Lexer) next() rune {

	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	runeValue, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return runeValue
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {

	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {

	line, col := l.getLineAndCol(l.start)

	l.tokens <- token.Token{
		Type:    token.ILLEGAL,
		Literal: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
	return nil
}
