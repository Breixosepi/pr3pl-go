/* Eugenio Giusepi Montilla Russo */
/*29958321*/

package token

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// identifiers + literals
	IDENT
	INT
	UNIT

	// Operators
	PLUS
	MINUS
	ASTERISK
	SLASH
	MOD
	LT
	GT
	EQ
	NEQ
	ASSIGN

	// Delimitadores
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	COMMA

	// conditionals
	IF
	THEN
	ELSE
	AND
	OR
	NOT

	// Keywords
	ISUNIT
	FST
	SND
	VAL
	LET
	IN
	END
	FUN
	TRUE
	FALSE
)

var tokenNames = map[TokenType]string{
	ILLEGAL:  "ILLEGAL",
	EOF:      "EOF",
	IDENT:    "IDENT",
	INT:      "INT",
	UNIT:     "UNIT",
	PLUS:     "PLUS",
	MINUS:    "MINUS",
	ASTERISK: "ASTERISK",
	SLASH:    "SLASH",
	MOD:      "MOD",
	LT:       "LT",
	GT:       "GT",
	EQ:       "EQ",
	NEQ:      "NEQ",
	ASSIGN:   "ASSIGN",
	LPAREN:   "LPAREN",
	RPAREN:   "RPAREN",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",
	COMMA:    "COMMA",
	IF:       "IF",
	THEN:     "THEN",
	ELSE:     "ELSE",
	AND:      "AND",
	OR:       "OR",
	NOT:      "NOT",
	ISUNIT:   "ISUNIT",
	FST:      "FST",
	SND:      "SND",
	VAL:      "VAL",
	LET:      "LET",
	IN:       "IN",
	END:      "END",
	FUN:      "FUN",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"isunit": ISUNIT,
	"fst":    FST,
	"snd":    SND,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"and":    AND,
	"or":     OR,
	"not":    NOT,
	"val":    VAL,
	"let":    LET,
	"in":     IN,
	"end":    END,
	"fun":    FUN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	name, ok := tokenNames[t.Type]
	if !ok {
		name = "UNKNOWN"
	}
	return fmt.Sprintf("%s(%q)", name, t.Literal)
}
