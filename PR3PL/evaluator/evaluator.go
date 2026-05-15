package evaluator

import (
	"pr3pl/ast"
	"pr3pl/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	UNIT  = &object.Unit{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.UnitLiteral:
		return UNIT

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ValStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
		return UNIT

	case *ast.LetExpression:
		val := Eval(node.Value, env)
		newEnv := object.NewEnclosedEnvironment(env)
		newEnv.Set(node.Name.Value, val)
		return Eval(node.Body, newEnv)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.FunctionLiteral:

		closure := &object.Closure{Env: env, Function: node}
		if node.Name != nil && node.Name.Value != "" {
			env.Set(node.Name.Value, closure)
		}

		return closure

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		arg := Eval(node.Argument, env)
		return applyFunction(function, arg)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PairExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return &object.Pair{Left: left, Right: right}

	case *ast.FstExpression:
		val := Eval(node.Argument, env)
		return val.(*object.Pair).Left

	case *ast.SndExpression:
		val := Eval(node.Argument, env)
		return val.(*object.Pair).Right

	case *ast.IsUnitExpression:
		val := Eval(node.Argument, env)
		if val.Type() == object.UNIT_OBJ {
			return TRUE
		}
		return FALSE
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {

	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {

	switch operator {
	case "-":
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}

	case "not":
		if right == TRUE {
			return FALSE
		}
		return TRUE

	default:
		return nil
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}
	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}
	return nil
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}

	case "-":
		return &object.Integer{Value: leftVal - rightVal}

	case "*":
		return &object.Integer{Value: leftVal * rightVal}

	case "/":
		if rightVal == 0 {
			panic("OperationError: division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}

	case "%":
		if rightVal == 0 {
			panic("OperationError: division by zero")
		}
		return &object.Integer{Value: leftVal % rightVal}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)

	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)

	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)

	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return nil
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "and":
		return nativeBoolToBooleanObject(leftVal && rightVal)

	case "or":
		return nativeBoolToBooleanObject(leftVal || rightVal)

	default:
		return nil
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {

	if input {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {

	condition := Eval(node.Condition, env)
	if condition == TRUE {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}
	return UNIT
}

func applyFunction(fn object.Object, arg object.Object) object.Object {

	closure := fn.(*object.Closure)
	functionNode := closure.Function
	extendedEnv := object.NewEnclosedEnvironment(closure.Env)
	extendedEnv.Set(functionNode.Parameter.Value, arg)
	extendedEnv.Set(functionNode.Name.Value, closure)
	return Eval(functionNode.Body, extendedEnv)
}
