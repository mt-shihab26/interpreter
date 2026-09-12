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
	}
	return nil
}
