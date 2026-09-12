package evaluator

import (
	"monkey/ast"
	"monkey/object"
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
	}
	return nil
}
