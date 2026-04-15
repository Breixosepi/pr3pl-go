package parser

import (
	"fmt"
	"strconv"

	"pr3pl/ast"
	"pr3pl/lexer"
	"pr3pl/token"
)

const (
	_ int = iota
	LOWEST
	LOGICAL_OR
	LOGICAL_AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.OR:       LOGICAL_OR,
	token.AND:      LOGICAL_AND,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.MOD:      PRODUCT,
	token.LPAREN:   CALL,
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.UNIT, p.parseUnitLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.LET, p.parseLetExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.FUN, p.parseFunctionLiteral)
	p.registerPrefix(token.FST, p.parseFstExpression)
	p.registerPrefix(token.SND, p.parseSndExpression)
	p.registerPrefix(token.ISUNIT, p.parseIsUnitExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedOrPairExpression)
	p.registerPrefix(token.LBRACKET, p.parseListLiteral)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("se esperaba que el siguiente token fuera %v, pero se obtuvo %v", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAL:
		return p.parseValStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseValStatement() *ast.ValStatement {
	stmt := &ast.ValStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no se encontró función de prefijo para el token %q", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("no se pudo parsear %q como entero", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseUnitLiteral() ast.Expression {
	return &ast.UnitLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()
	expression.Consequence = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		expression.Alternative = p.parseExpression(LOWEST)
	}

	return expression
}

func (p *Parser) parseLetExpression() ast.Expression {
	expression := &ast.LetExpression{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	expression.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	expression.Body = p.parseExpression(LOWEST)

	if !p.expectPeek(token.END) {
		return nil
	}

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	lit.Parameter = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	lit.Body = p.parseExpression(LOWEST)

	return lit
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	p.nextToken()

	if !p.curTokenIs(token.RPAREN) {
		exp.Argument = p.parseExpression(LOWEST)
		p.nextToken()
	}

	if !p.curTokenIs(token.RPAREN) {
		p.peekError(token.RPAREN)
		return nil
	}

	return exp
}

func (p *Parser) parseGroupedOrPairExpression() ast.Expression {
	lparenToken := p.curToken

	p.nextToken()
	leftExp := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		rightExp := p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		return &ast.PairExpression{
			Token: lparenToken,
			Left:  leftExp,
			Right: rightExp,
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return leftExp
}

func (p *Parser) parseFstExpression() ast.Expression {
	exp := &ast.FstExpression{Token: p.curToken}
	p.nextToken()
	exp.Argument = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseSndExpression() ast.Expression {
	exp := &ast.SndExpression{Token: p.curToken}
	p.nextToken()
	exp.Argument = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseIsUnitExpression() ast.Expression {
	exp := &ast.IsUnitExpression{Token: p.curToken}
	p.nextToken()
	exp.Argument = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseListLiteral() ast.Expression {
	listToken := p.curToken

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return &ast.UnitLiteral{Token: p.curToken}
	}

	p.nextToken()
	var elements []ast.Expression
	elements = append(elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		elements = append(elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	var current ast.Expression = &ast.UnitLiteral{Token: p.curToken}

	for i := len(elements) - 1; i >= 0; i-- {
		current = &ast.PairExpression{
			Token: listToken,
			Left:  elements[i],
			Right: current,
		}
	}

	return current
}
