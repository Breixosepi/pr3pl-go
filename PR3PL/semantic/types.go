package semantic

import "pr3pl/ast"

type Type interface {
	Signature() string
}

// atomic types

type IntType struct{}

func (i *IntType) Signature() string { return "int" }

type UnitType struct{}

func (u *UnitType) Signature() string { return "unit" }

type BoolType struct{}

func (b *BoolType) Signature() string { return "bool" }

// compound types

type PairType struct {
	First  Type
	Second Type
}

func (p *PairType) Signature() string {
	return "(" + p.First.Signature() + ", " + p.Second.Signature() + ")"
}

type ClosureType struct {
	Env        *Environment
	Function   *ast.FunctionLiteral
	IsChecking bool
	ReturnType Type
}

func (c *ClosureType) Signature() string {
	return "closure"
}
