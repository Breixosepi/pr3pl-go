package transpiler

import (
	"fmt"
	"pr3pl/ast"
)

func ToOriginalPR3PL(node ast.Node) string {

	switch n := node.(type) {

	case *ast.Program:
		var out string
		for _, stmt := range n.Statements {
			out += ToOriginalPR3PL(stmt) + "\n"
		}
		return out

	case *ast.ExpressionStatement:
		return ToOriginalPR3PL(n.Expression)

	case *ast.IntegerLiteral:
		return fmt.Sprintf("(%d)", n.Value)

	case *ast.Identifier:
		return fmt.Sprintf("(var (%s))", n.Value)

	case *ast.ValStatement:
		if fn, ok := n.Value.(*ast.FunctionLiteral); ok {
			return fmt.Sprintf("(fun (%s) (%s) %s)",
				n.Name.Value, fn.Parameter.Value, ToOriginalPR3PL(fn.Body))
		}
		return fmt.Sprintf("(val (%s) %s)", n.Name.Value, ToOriginalPR3PL(n.Value))

	case *ast.LetExpression:
		return fmt.Sprintf("(let (%s) %s %s)",
			n.Name.Value, ToOriginalPR3PL(n.Value), ToOriginalPR3PL(n.Body))

	case *ast.FunctionLiteral:
		name := "anon"
		if n.Name != nil {
			name = n.Name.Value
		}
		return fmt.Sprintf("(fun (%s) (%s) %s)",
			name, n.Parameter.Value, ToOriginalPR3PL(n.Body))

	case *ast.CallExpression:
		return fmt.Sprintf("(call (%s) %s)",
			n.Function.String(), ToOriginalPR3PL(n.Argument))

	case *ast.PrefixExpression:
		switch n.Operator {

		case "-":
			return fmt.Sprintf("(- %s)", ToOriginalPR3PL(n.Right))

		case "not":
			return fmt.Sprintf("(iflesser (0) %s (0) (1))", ToOriginalPR3PL(n.Right))
		}

	case *ast.InfixExpression:

		l := ToOriginalPR3PL(n.Left)
		r := ToOriginalPR3PL(n.Right)

		switch n.Operator {

		case "+", "*", "/", "%":
			return fmt.Sprintf("(%s %s %s)", n.Operator, l, r)

		case "-":
			return fmt.Sprintf("(+ %s (- %s))", l, r)

		case "<":
			return fmt.Sprintf("(iflesser %s %s (1) (0))", l, r)

		case ">":
			return fmt.Sprintf("(iflesser %s %s (1) (0))", r, l)

		case "==":
			return fmt.Sprintf("(iflesser %s %s (0) (iflesser %s %s (0) (1)))", l, r, r, l)

		case "!=":
			return fmt.Sprintf("(iflesser %s %s (1) (iflesser %s %s (1) (0)))", l, r, r, l)

		case "and":
			return fmt.Sprintf("(iflesser (0) %s (iflesser (0) %s (1) (0)) (0))", l, r)

		case "or":
			return fmt.Sprintf("(iflesser (0) %s (1) (iflesser (0) %s (1) (0)))", l, r)
		}

	case *ast.PairExpression:
		return fmt.Sprintf("(pair %s %s)", ToOriginalPR3PL(n.Left), ToOriginalPR3PL(n.Right))

	case *ast.FstExpression:
		return fmt.Sprintf("(fst %s)", ToOriginalPR3PL(n.Argument))

	case *ast.SndExpression:
		return fmt.Sprintf("(snd %s)", ToOriginalPR3PL(n.Argument))

	case *ast.IsUnitExpression:
		return fmt.Sprintf("(isunit %s)", ToOriginalPR3PL(n.Argument))

	case *ast.IfExpression:
		return toIfLesser(n)

	case *ast.UnitLiteral:
		return "()"

	case *ast.Boolean:
		if n.Value {
			return "(1)"
		}
		return "(0)"
	}
	return ""
}

func toIfLesser(n *ast.IfExpression) string {

	alt := "()"

	if n.Alternative != nil {
		alt = ToOriginalPR3PL(n.Alternative)
	}

	if infix, ok := n.Condition.(*ast.InfixExpression); ok && infix.Operator == "<" {
		return fmt.Sprintf("(iflesser %s %s %s %s)",
			ToOriginalPR3PL(infix.Left), ToOriginalPR3PL(infix.Right),
			ToOriginalPR3PL(n.Consequence), alt)
	}

	return fmt.Sprintf("(iflesser (0) %s %s %s)",
		ToOriginalPR3PL(n.Condition), ToOriginalPR3PL(n.Consequence), alt)
}
