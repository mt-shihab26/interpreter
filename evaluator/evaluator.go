package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement)
		}
		return result
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerExpression:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.UnaryExpression:
		right := Eval(node.RightExpression)
		switch node.Operator {
		case "!":
			switch right {
			case TRUE:
				return FALSE
			case FALSE:
				return TRUE
			case NULL:
				return TRUE
			default:
				return FALSE
			}
		case "-":
			switch right.Type() {
			case object.INTEGER:
				return &object.Integer{Value: -(right.(*object.Integer).Value)}
			default:
				return NULL
			}
		}
		return NULL
	default:
		return nil
	}
}
