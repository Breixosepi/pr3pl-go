/* Eugenio Giusepi Montilla Russo */
/* 29958321 */

package lexer

import (
	"pr3pl/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let x = 10 in x + 5 end
/* Este es un comentario principal
   /* Y este es un comentario anidado que el lexer también debe ignorar */
   Seguimos en el comentario principal */
if x < (10, 2) then 1 else 0
(())`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.IN, "in"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.END, "end"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.LPAREN, "("},
		{token.INT, "10"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.INT, "1"},
		{token.ELSE, "else"},
		{token.INT, "0"},
		{token.LPAREN, "("},
		{token.UNIT, "()"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New("test", input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test [%d] - wrong token type. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test [%d] - wrong literal. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
