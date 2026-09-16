package eval

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return newIntegerObject(int64(len(arg.Elements)))
			case *object.String:
				return newIntegerObject(int64(len(arg.Value)))
			default:
				return newErrorObject("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) >= 1 {
					return arg.Elements[0]
				}
				return NULL_OBJECT
			default:
				return newErrorObject("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) >= 1 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL_OBJECT
			default:
				return newErrorObject("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorObject("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return NULL_OBJECT
			default:
				return newErrorObject("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorObject("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newErrorObject("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
}
