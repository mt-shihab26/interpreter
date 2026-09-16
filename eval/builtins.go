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
}
