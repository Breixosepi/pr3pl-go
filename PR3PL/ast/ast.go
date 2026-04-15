package ast

import (
	"bytes"
	"pr3pl/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type UnitLiteral struct {
	Token token.Token
}

func (ul *UnitLiteral) expressionNode()      {}
func (ul *UnitLiteral) TokenLiteral() string { return ul.Token.Literal }
func (ul *UnitLiteral) String() string       { return "()" }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)

	if pe.Operator == "not" {
		out.WriteString(" ")
	}

	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type ValStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *ValStatement) statementNode()       {}
func (vs *ValStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *ValStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" then ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type LetExpression struct {
	Token token.Token
	Name  *Identifier
	Value Expression
	Body  Expression
}

func (le *LetExpression) expressionNode()      {}
func (le *LetExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LetExpression) String() string {
	var out bytes.Buffer
	out.WriteString("let ")
	out.WriteString(le.Name.String())
	out.WriteString(" = ")
	out.WriteString(le.Value.String())
	out.WriteString(" in ")
	out.WriteString(le.Body.String())
	out.WriteString(" end")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type FunctionLiteral struct {
	Token     token.Token
	Name      *Identifier
	Parameter *Identifier
	Body      Expression
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("fun ")
	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}
	out.WriteString("(")
	if fl.Parameter != nil {
		out.WriteString(fl.Parameter.String())
	}
	out.WriteString(") = ")
	if fl.Body != nil {
		out.WriteString(fl.Body.String())
	}
	return out.String()
}

type CallExpression struct {
	Token    token.Token
	Function Expression
	Argument Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	if ce.Argument != nil {
		out.WriteString(ce.Argument.String())
	}
	out.WriteString(")")
	return out.String()
}

type PairExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (pe *PairExpression) expressionNode()      {}
func (pe *PairExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PairExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(", ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type FstExpression struct {
	Token    token.Token
	Argument Expression
}

func (fe *FstExpression) expressionNode()      {}
func (fe *FstExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *FstExpression) String() string {
	return "fst(" + fe.Argument.String() + ")"
}

type SndExpression struct {
	Token    token.Token
	Argument Expression
}

func (se *SndExpression) expressionNode()      {}
func (se *SndExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SndExpression) String() string {
	return "snd(" + se.Argument.String() + ")"
}

type IsUnitExpression struct {
	Token    token.Token
	Argument Expression
}

func (iue *IsUnitExpression) expressionNode()      {}
func (iue *IsUnitExpression) TokenLiteral() string { return iue.Token.Literal }
func (iue *IsUnitExpression) String() string {
	return "isunit(" + iue.Argument.String() + ")"
}
