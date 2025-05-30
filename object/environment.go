package object

// Enclosed Env (Outer)
// ┌────────────────────────┐
// │ outer variables        │
// │                        │
// │   Extended Env (Inner) │
// │   ┌─────────────────┐  │
// │   │ function args   │  │
// │   │ local variables │  │
// │   └─────────────────┘  │
// └────────────────────────┘

type Environment struct {
	store map[string]Object
	outer *Environment // The enclosing env of the current env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		// Recursively get the outermost environment
		// Until we find an associated value or there is no outer environment
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
