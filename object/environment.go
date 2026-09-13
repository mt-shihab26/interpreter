package object

import "maps"

type Environment struct {
	store map[string]Object
	outer map[string]Object
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

func NewEnclosedEnvironment(env *Environment) *Environment {
	store := make(map[string]Object)
	outer := maps.Clone(env.store)
	return &Environment{store: store, outer: outer}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if ok {
		return obj, ok
	}
	obj, ok = e.outer[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
