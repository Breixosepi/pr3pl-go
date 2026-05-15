package semantic

import (
	"fmt"
	"pr3pl/ast"
	"strings"
)

func TypeCheck(node ast.Node, env *Environment) (Type, error) {

	switch node := node.(type) {

	case *ast.Program:
		return typeCheckProgram(node, env)

	case *ast.ExpressionStatement:
		return TypeCheck(node.Expression, env)

	case *ast.IntegerLiteral:
		return &IntType{}, nil

	case *ast.UnitLiteral:
		return &UnitType{}, nil

	case *ast.Boolean:
		return &BoolType{}, nil

	case *ast.InfixExpression:
		return typeCheckInfixExpression(node, env)

	case *ast.PrefixExpression:
		return typeCheckPrefixExpression(node, env)

	case *ast.Identifier:
		return typeCheckIdentifier(node, env)

	case *ast.ValStatement:
		return typeCheckValStatement(node, env)

	case *ast.LetExpression:
		return typeCheckLetExpression(node, env)

	case *ast.FunctionLiteral:
		return typeCheckFunctionLiteral(node, env)

	case *ast.CallExpression:
		return typeCheckCallExpression(node, env)

	case *ast.IfExpression:
		return typeCheckIfExpression(node, env)

	case *ast.PairExpression:
		return typeCheckPairExpression(node, env)

	case *ast.FstExpression:
		return typeCheckFstExpression(node, env)

	case *ast.SndExpression:
		return typeCheckSndExpression(node, env)

	case *ast.IsUnitExpression:
		return typeCheckIsUnitExpression(node, env)

	default:
		return nil, fmt.Errorf("error semántico: nodo no soportado %T", node)
	}
}

func typeCheckProgram(program *ast.Program, env *Environment) (Type, error) {

	var result Type
	var err error

	for _, statement := range program.Statements {
		result, err = TypeCheck(statement, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func typeCheckInfixExpression(node *ast.InfixExpression, env *Environment) (Type, error) {

	leftType, err := TypeCheck(node.Left, env)
	if err != nil {
		return nil, err
	}

	rightType, err := TypeCheck(node.Right, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "+", "-", "*", "/", "%":
		if leftType.Signature() != "int" || rightType.Signature() != "int" {
			return nil, fmt.Errorf("TypeError: %s requiere int, obtuvo %s y %s", node.Operator, leftType.Signature(), rightType.Signature())
		}
		return &IntType{}, nil

	case "<", ">", "==", "!=":
		if leftType.Signature() != "int" || rightType.Signature() != "int" {
			return nil, fmt.Errorf("TypeError: %s requiere int", node.Operator)
		}
		return &BoolType{}, nil

	case "and", "or":
		if leftType.Signature() != "bool" || rightType.Signature() != "bool" {
			return nil, fmt.Errorf("TypeError: %s requiere bool", node.Operator)
		}
		return &BoolType{}, nil
	}
	return nil, fmt.Errorf("operador desconocido %s", node.Operator)
}

func typeCheckPrefixExpression(node *ast.PrefixExpression, env *Environment) (Type, error) {

	rightType, err := TypeCheck(node.Right, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "-":
		if rightType.Signature() != "int" {
			return nil, fmt.Errorf("TypeError: negación requiere int")
		}
		return &IntType{}, nil
	case "not":
		if rightType.Signature() != "bool" {
			return nil, fmt.Errorf("TypeError: not requiere bool")
		}
		return &BoolType{}, nil
	}
	return nil, fmt.Errorf("operador desconocido %s", node.Operator)
}

func typeCheckIdentifier(node *ast.Identifier, env *Environment) (Type, error) {

	if val, ok := env.Get(node.Value); ok {
		return val, nil
	}
	return nil, fmt.Errorf("NotFoundError: la variable %s no existe", node.Value)
}

func typeCheckValStatement(node *ast.ValStatement, env *Environment) (Type, error) {

	valType, err := TypeCheck(node.Value, env)
	if err != nil {
		return nil, err
	}
	env.Set(node.Name.Value, valType)

	return &UnitType{}, nil
}

func typeCheckLetExpression(node *ast.LetExpression, env *Environment) (Type, error) {

	valType, err := TypeCheck(node.Value, env)
	if err != nil {
		return nil, err
	}
	newEnv := NewEnclosedEnvironment(env)
	newEnv.Set(node.Name.Value, valType)

	return TypeCheck(node.Body, newEnv)
}

func typeCheckFunctionLiteral(node *ast.FunctionLiteral, env *Environment) (Type, error) {

	closure := &ClosureType{
		Env:      env,
		Function: node,
	}

	if node.Name != nil && node.Name.Value != "" {
		env.Set(node.Name.Value, closure)
	}

	return closure, nil
}

// optimizar esta
func typeCheckCallExpression(node *ast.CallExpression, env *Environment) (Type, error) {

	funcType, err := TypeCheck(node.Function, env)
	if err != nil {
		return nil, err
	}

	closure, ok := funcType.(*ClosureType)
	if !ok {
		return nil, fmt.Errorf("TypeError: call aplicado a un no-closure")
	}

	argType, err := TypeCheck(node.Argument, env)
	if err != nil {
		return nil, err
	}

	if closure.IsChecking {
		if closure.ReturnType != nil {
			return closure.ReturnType, nil
		}
		return &IntType{}, nil
	}

	closure.IsChecking = true
	functionNode := closure.Function

	closureEnv := NewEnclosedEnvironment(closure.Env)
	closureEnv.Set(functionNode.Parameter.Value, argType)
	closureEnv.Set(functionNode.Name.Value, closure)

	result, err := TypeCheck(functionNode.Body, closureEnv)
	closure.ReturnType = result
	closure.IsChecking = false

	return result, err
}

func areTypesCompatible(t1 Type, t2 Type) bool {

	sig1 := t1.Signature()
	sig2 := t2.Signature()

	if sig1 == sig2 {
		return true
	}

	isList1 := strings.HasPrefix(sig1, "(") || sig1 == "unit" || sig1 == "int"
	isList2 := strings.HasPrefix(sig2, "(") || sig2 == "unit" || sig2 == "int"

	return isList1 && isList2
}

func typeCheckIfExpression(node *ast.IfExpression, env *Environment) (Type, error) {

	condType, err := TypeCheck(node.Condition, env)
	if err != nil {
		return nil, err
	}

	if condType.Signature() != "bool" {
		return nil, fmt.Errorf("TypeError: la condicion del if debe ser booleana , got %s", condType.Signature())
	}

	consType, err := TypeCheck(node.Consequence, env)
	if err != nil {
		return nil, err
	}

	if node.Alternative != nil {
		altType, err := TypeCheck(node.Alternative, env)
		if err != nil {
			return nil, err
		}

		if !areTypesCompatible(consType, altType) {
			return nil, fmt.Errorf("TypeError: discrepancia de tipos en condicional. Ramas devuelven %s y %s", consType.Signature(), altType.Signature())
		}

		if strings.HasPrefix(consType.Signature(), "(") {
			return consType, nil
		}
		return altType, nil
	}

	return consType, nil
}

func typeCheckPairExpression(node *ast.PairExpression, env *Environment) (Type, error) {

	leftType, err := TypeCheck(node.Left, env)
	if err != nil {
		return nil, err
	}

	rightType, err := TypeCheck(node.Right, env)
	if err != nil {
		return nil, err
	}

	return &PairType{First: leftType, Second: rightType}, nil
}

func typeCheckFstExpression(node *ast.FstExpression, env *Environment) (Type, error) {

	argType, err := TypeCheck(node.Argument, env)
	if err != nil {
		return nil, err
	}

	pair, ok := argType.(*PairType)
	if !ok {
		if argType.Signature() == "int" {
			return &IntType{}, nil
		}
		return nil, fmt.Errorf("TypeError: el operador fst requiere un operando PairType, se obtuvo %s", argType.Signature())
	}
	return pair.First, nil
}

func typeCheckSndExpression(node *ast.SndExpression, env *Environment) (Type, error) {

	argType, err := TypeCheck(node.Argument, env)
	if err != nil {
		return nil, err
	}

	pair, ok := argType.(*PairType)
	if !ok {
		if argType.Signature() == "int" {
			return &IntType{}, nil
		}
		return nil, fmt.Errorf("TypeError: el operador snd requiere un operando PairType, se obtuvo %s", argType.Signature())
	}
	return pair.Second, nil
}

func typeCheckIsUnitExpression(node *ast.IsUnitExpression, env *Environment) (Type, error) {

	_, err := TypeCheck(node.Argument, env)
	if err != nil {
		return nil, err
	}
	return &BoolType{}, nil
}
