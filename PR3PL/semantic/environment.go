package semantic

type Environment struct {
	store map[string]Type
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Type),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Type, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Type) Type {
	e.store[name] = val
	return val
}
