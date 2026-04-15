/* Eugenio Giusepi Montilla Russo */
/* 29958321 */

package lexer

import (
	"pr3pl/token"
	"strings"
)

func lexTopLevel(l *Lexer) stateFn {
	for {
		r := l.next()

		if isSpace(r) {
			l.ignore()
			continue
		}

		if r == eof {
			l.emit(token.EOF)
			return nil
		}

		switch r {
		case '+':
			l.emit(token.PLUS)
		case '-':
			l.emit(token.MINUS)
		case '*':
			l.emit(token.ASTERISK)
		case '/':
			if l.peek() == '*' {
				l.next()
				return nestedComments
			}

			l.emit(token.SLASH)
		case '%':
			l.emit(token.MOD)
		case '<':
			l.emit(token.LT)
		case '>':
			l.emit(token.GT)
		case '!':
			if l.peek() == '=' {
				l.emit(token.NEQ)
			} else {
				l.emit(token.ILLEGAL)
			}
		case '=':
			if l.peek() == '=' {
				l.next()
				l.emit(token.EQ)
			} else {
				l.emit(token.ASSIGN)
			}
		case '(':
			if l.peek() == ')' {
				l.next()
				l.emit(token.UNIT)
			} else {
				l.emit(token.LPAREN)
			}
		case ')':
			l.emit(token.RPAREN)
		case '[':
			l.emit(token.LBRACKET)
		case ']':
			l.emit(token.RBRACKET)
		case ',':
			l.emit(token.COMMA)
		default:
			if isAlpha(r) {
				l.backup()
				return lexIdentifier
			} else if isDigit(r) {
				l.backup()
				return lexNumber
			}
			l.emit(token.ILLEGAL)
		}
		return lexTopLevel
	}
}

func lexIdentifier(l *Lexer) stateFn {
	for {
		r := l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	word := l.input[l.start:l.pos]
	l.emit(token.LookupIdent(word))
	return lexTopLevel
}

func lexNumber(l *Lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	l.emit(token.INT)
	return lexTopLevel
}

func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}

func nestedComments(l *Lexer) stateFn {
	depth := 1

	for {
		r := l.next()
		if r == eof {
			return l.errorf("unclosed comment")
		}
		if r == '/' && l.peek() == '*' {
			l.next()
			depth++
			continue
		}

		if r == '*' && l.peek() == '/' {
			l.next()
			depth--

			if depth == 0 {
				l.ignore()
				return lexTopLevel
			}
		}
	}
}
