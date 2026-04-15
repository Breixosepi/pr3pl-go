package parser

import (
	"pr3pl/ast"
	"pr3pl/lexer"
	"testing"
)

func TestIdentifierExpression(t *testing.T) {
	input := "mi_variable"

	l := lexer.New("test", input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("El programa no tiene suficientes statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] no es un ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("La expresión no es un *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "mi_variable" {
		t.Errorf("ident.Value incorrecto. expected=%s, got=%s", "mi_variable", ident.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New("test", input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("El programa no tiene suficientes statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] no es un ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("La expresión no es un *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value incorrecto. expected=%d, got=%d", 5, literal.Value)
	}
}

func TestUnitLiteralExpression(t *testing.T) {
	input := "()"

	l := lexer.New("test", input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("El programa no tiene suficientes statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] no es un ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	_, ok = stmt.Expression.(*ast.UnitLiteral)
	if !ok {
		t.Fatalf("La expresión no es un *ast.UnitLiteral. got=%T", stmt.Expression)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("El parser tiene %d errores", len(errors))
	for _, msg := range errors {
		t.Errorf("Error de parser: %q", msg)
	}
	t.FailNow()
}
