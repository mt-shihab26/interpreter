package object

import "fmt"

// Type identifies the kind of an Object.
type Type string

// The Type value for each kind of Object.
const (
	NULL    = "NULL"
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	RETURN  = "RETURN"
	ERROR   = "ERROR"
)

// Object is implemented by every value the Monkey evaluator produces.
type Object interface {
	// Type returns the object's Type.
	Type() Type
	// Inspect returns a human-readable representation of the object's value.
	Inspect() string
}

// Null implements the Object interface.
type Null struct {
}

// Type returns NULL.
func (i *Null) Type() Type {
	return NULL
}

// Inspect returns "null".
func (i *Null) Inspect() string {
	return "null"
}

// Integer implements the Object interface.
type Integer struct {
	Value int64
}

// Type returns INTEGER.
func (i *Integer) Type() Type {
	return INTEGER
}

// Inspect returns the integer's value as a string.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%v", i.Value)
}

// Boolean implements the Object interface.
type Boolean struct {
	Value bool
}

// Type returns BOOLEAN.
func (b *Boolean) Type() Type {
	return BOOLEAN
}

// Inspect returns the boolean's value as a string.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

// Return implements the Object interface.
type Return struct {
	Value Object
}

// Type returns BOOLEAN.
func (b *Return) Type() Type {
	return RETURN
}

// Inspect returns the return's value as a string.
func (b *Return) Inspect() string {
	return b.Value.Inspect()
}

// Error implements the Object interface.
type Error struct {
	Message string
}

// Type returns ERROR.
func (e *Error) Type() Type {
	return ERROR
}

// Inspect returns the error's value as a string.
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
