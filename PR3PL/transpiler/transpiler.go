package transpiler

import (
	"fmt"
	"pr3pl/ast"
	"strings"
)

func indent(depth int) string {
	return strings.Repeat("    ", depth)
}

func ToOriginalPR3PL(node ast.Node) string {
	return toPR3PL(node, 0)
}

func toPR3PL(node ast.Node, depth int) string {

	switch n := node.(type) {

	case *ast.Program:
		var out string
		for _, stmt := range n.Statements {
			out += toPR3PL(stmt, depth) + "\n\n"
		}
		return strings.TrimSpace(out)

	case *ast.ExpressionStatement:
		return toPR3PL(n.Expression, depth)

	case *ast.IntegerLiteral:
		return fmt.Sprintf("(%d)", n.Value)

	case *ast.Identifier:
		return fmt.Sprintf("(var (%s))", n.Value)

	case *ast.ValStatement:
		if fn, ok := n.Value.(*ast.FunctionLiteral); ok {
			return fmt.Sprintf("(fun (%s) (%s)\n%s%s)",
				n.Name.Value, fn.Parameter.Value, indent(depth+1), toPR3PL(fn.Body, depth+1))
		}
		return fmt.Sprintf("(val (%s)\n%s%s)", n.Name.Value, indent(depth+1), toPR3PL(n.Value, depth+1))

	case *ast.LetExpression:
		return fmt.Sprintf("(let (%s)\n%s%s\n%s%s)",
			n.Name.Value,
			indent(depth+1), toPR3PL(n.Value, depth+1),
			indent(depth+1), toPR3PL(n.Body, depth+1))

	case *ast.FunctionLiteral:
		name := "anon"
		if n.Name != nil {
			name = n.Name.Value
		}
		return fmt.Sprintf("(fun (%s) (%s)\n%s%s)",
			name, n.Parameter.Value, indent(depth+1), toPR3PL(n.Body, depth+1))

	case *ast.CallExpression:
		return fmt.Sprintf("(call (%s) %s)",
			n.Function.String(), toPR3PL(n.Argument, depth))

	case *ast.PrefixExpression:
		switch n.Operator {
		case "-":
			return fmt.Sprintf("(- %s)", toPR3PL(n.Right, depth))
		case "not":
			return fmt.Sprintf("(iflesser (0) %s\n%s(0)\n%s(1))",
				toPR3PL(n.Right, depth), indent(depth+1), indent(depth+1))
		}

	case *ast.InfixExpression:
		l := toPR3PL(n.Left, depth)
		r := toPR3PL(n.Right, depth)

		switch n.Operator {
		case "+", "*", "/":
			return fmt.Sprintf("(%s %s %s)", n.Operator, l, r)
		case "-":
			return fmt.Sprintf("(+ %s (- %s))", l, r)
		case "%":
			return fmt.Sprintf("(+ %s (- (* (/ %s %s) %s)))", l, l, r, r)
		case "<":
			return fmt.Sprintf("(iflesser %s %s (1) (0))", l, r)
		case ">":
			return fmt.Sprintf("(iflesser %s %s (1) (0))", r, l)
		case "==":
			return fmt.Sprintf("(iflesser %s %s\n%s(0)\n%s(iflesser %s %s (0) (1)))",
				l, r, indent(depth+1), indent(depth+1), r, l)
		case "!=":
			return fmt.Sprintf("(iflesser %s %s\n%s(1)\n%s(iflesser %s %s (1) (0)))",
				l, r, indent(depth+1), indent(depth+1), r, l)
		case "and":
			return fmt.Sprintf("(iflesser (0) %s\n%s(iflesser (0) %s (1) (0))\n%s(0))",
				l, indent(depth+1), r, indent(depth+1))
		case "or":
			return fmt.Sprintf("(iflesser (0) %s\n%s(1)\n%s(iflesser (0) %s (1) (0)))",
				l, indent(depth+1), indent(depth+1), r)
		}

	case *ast.PairExpression:
		return fmt.Sprintf("(pair %s %s)", toPR3PL(n.Left, depth), toPR3PL(n.Right, depth))

	case *ast.FstExpression:
		return fmt.Sprintf("(fst %s)", toPR3PL(n.Argument, depth))

	case *ast.SndExpression:
		return fmt.Sprintf("(snd %s)", toPR3PL(n.Argument, depth))

	case *ast.IsUnitExpression:
		return fmt.Sprintf("(isunit %s)", toPR3PL(n.Argument, depth))

	case *ast.IfExpression:
		return toIfLesser(n, depth)

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

func toIfLesser(n *ast.IfExpression, depth int) string {
	alt := "()"
	if n.Alternative != nil {
		alt = toPR3PL(n.Alternative, depth+1)
	}

	if infix, ok := n.Condition.(*ast.InfixExpression); ok && infix.Operator == "<" {
		return fmt.Sprintf("(iflesser %s %s\n%s%s\n%s%s)",
			toPR3PL(infix.Left, depth), toPR3PL(infix.Right, depth),
			indent(depth+1), toPR3PL(n.Consequence, depth+1),
			indent(depth+1), alt)
	}

	return fmt.Sprintf("(iflesser (0) %s\n%s%s\n%s%s)",
		toPR3PL(n.Condition, depth),
		indent(depth+1), toPR3PL(n.Consequence, depth+1),
		indent(depth+1), alt)
}
