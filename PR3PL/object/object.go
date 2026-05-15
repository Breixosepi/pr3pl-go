package object

import (
	"fmt"
	"pr3pl/ast"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	UNIT_OBJ    = "UNIT"
	PAIR_OBJ    = "PAIR"
	CLOSURE_OBJ = "CLOSURE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Unit struct{}

func (u *Unit) Type() ObjectType { return UNIT_OBJ }
func (u *Unit) Inspect() string  { return "()" }

type Pair struct {
	Left  Object
	Right Object
}

func (p *Pair) Type() ObjectType { return PAIR_OBJ }
func (p *Pair) Inspect() string {
	return "(" + p.Left.Inspect() + ", " + p.Right.Inspect() + ")"
}

type Closure struct {
	Env      *Environment
	Function *ast.FunctionLiteral
}

func (c *Closure) Type() ObjectType { return CLOSURE_OBJ }
func (c *Closure) Inspect() string {
	name := "anonima"
	if c.Function.Name != nil {
		name = c.Function.Name.Value
	}
	return fmt.Sprintf("closure(fun %s)", name)
}
