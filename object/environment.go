package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return NewEnclosedEnvironment(nil)
}

func NewEnclosedEnvironment(env *Environment) *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: env}
}

func (e *Environment) Get(name string) (Object, bool) {
	if obj, ok := e.store[name]; ok {
		return obj, ok
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false

}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
