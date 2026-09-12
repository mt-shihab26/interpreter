package eval

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
		return evalStatements(node.Statements)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerExpression:
		return newIntegerObject(node.Value)
	case *ast.BooleanExpression:
		return newBooleanObject(node.Value)
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
	case *ast.BinaryExpression:
		left := Eval(node.LeftExpression)
		right := Eval(node.RightExpression)
		switch {
		case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value
			switch node.Operator {
			case "+":
				return newIntegerObject(leftValue + rightValue)
			case "-":
				return newIntegerObject(leftValue - rightValue)
			case "*":
				return newIntegerObject(leftValue * rightValue)
			case "/":
				return newIntegerObject(leftValue / rightValue)
			case "<":
				return newBooleanObject(leftValue < rightValue)
			case ">":
				return newBooleanObject(leftValue > rightValue)
			case "==":
				return newBooleanObject(leftValue == rightValue)
			case "!=":
				return newBooleanObject(leftValue != rightValue)
			default:
				return NULL
			}
		case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
			leftValue := left.(*object.Boolean).Value
			rightValue := right.(*object.Boolean).Value
			switch node.Operator {
			case "==":
				return newBooleanObject(leftValue == rightValue)
			case "!=":
				return newBooleanObject(leftValue != rightValue)
			default:
				return NULL
			}
		default:
			return NULL
		}
	case *ast.IfExpression:
		condition := Eval(node.ConditionExpression)
		if isTruthy(condition) {
			return Eval(node.ConsequenceStatement)
		} else {
			if node.AlternativeStatement != nil {
				return Eval(node.AlternativeStatement)
			}
		}
		return NULL
	default:
		return nil
	}
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
	}
	return result
}

func newIntegerObject(value int64) *object.Integer {
	return &object.Integer{Value: value}
}

func newBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	if obj.Type() == object.INTEGER {
		if obj.(*object.Integer).Value == 0 {
			return false
		} else {
			return true
		}
	}
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}
