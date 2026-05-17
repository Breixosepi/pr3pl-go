package optimizer

import (
	"fmt"
	"pr3pl/ast"
	"pr3pl/token"
)

func Optimize(node ast.Node) ast.Node {

	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.Program:
		for i, stmt := range n.Statements {
			n.Statements[i] = Optimize(stmt).(ast.Statement)
		}
		return n

	case *ast.ExpressionStatement:
		n.Expression = Optimize(n.Expression).(ast.Expression)
		return n

	case *ast.ValStatement:
		n.Value = Optimize(n.Value).(ast.Expression)
		return n

	case *ast.LetExpression:
		n.Value = Optimize(n.Value).(ast.Expression)
		n.Body = Optimize(n.Body).(ast.Expression)
		return n

	case *ast.FunctionLiteral:
		n.Body = Optimize(n.Body).(ast.Expression)
		return n

	case *ast.CallExpression:
		n.Function = Optimize(n.Function).(ast.Expression)
		n.Argument = Optimize(n.Argument).(ast.Expression)
		return n

	case *ast.PrefixExpression:
		n.Right = Optimize(n.Right).(ast.Expression)
		return n

	case *ast.PairExpression:
		n.Left = Optimize(n.Left).(ast.Expression)
		n.Right = Optimize(n.Right).(ast.Expression)
		return n

	case *ast.IfExpression:
		n.Condition = Optimize(n.Condition).(ast.Expression)
		n.Consequence = Optimize(n.Consequence).(ast.Expression)
		if n.Alternative != nil {
			n.Alternative = Optimize(n.Alternative).(ast.Expression)
		}
		return n

	case *ast.FstExpression:
		n.Argument = Optimize(n.Argument).(ast.Expression)
		return n

	case *ast.SndExpression:
		n.Argument = Optimize(n.Argument).(ast.Expression)
		return n

	case *ast.IsUnitExpression:
		n.Argument = Optimize(n.Argument).(ast.Expression)
		return n

	case *ast.InfixExpression:
		n.Left = Optimize(n.Left).(ast.Expression)
		n.Right = Optimize(n.Right).(ast.Expression)

		leftInt, leftOk := n.Left.(*ast.IntegerLiteral)
		rightInt, rightOk := n.Right.(*ast.IntegerLiteral)

		if leftOk && rightOk {
			var result int64

			switch n.Operator {
			case "+":
				result = leftInt.Value + rightInt.Value
			case "-":
				result = leftInt.Value - rightInt.Value
			case "*":
				result = leftInt.Value * rightInt.Value
			case "/":
				if rightInt.Value != 0 {
					result = leftInt.Value / rightInt.Value
				} else {
					return n
				}
			case "%":
				if rightInt.Value != 0 {
					result = leftInt.Value % rightInt.Value
				} else {
					return n
				}
			default:
				return n
			}

			return &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: fmt.Sprintf("%d", result)},
				Value: result,
			}
		}
		return n
	}

	return node
}
