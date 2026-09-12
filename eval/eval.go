package eval

import (
	"monkey/ast"
	"monkey/object"
)

var (
	NULL_OBJECT  = &object.Null{}
	TRUE_OBJECT  = &object.Boolean{Value: true}
	FALSE_OBJECT = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement)
			if result, ok := result.(*object.Return); ok {
				return result.Value
			}
		}
		return result
	case *ast.BlockStatement:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement)
			if result, ok := result.(*object.Return); ok {
				return result
			}
		}
		return result
	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.ValueExpression)}
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
			return newBooleanObject(!isTruthy(right))
		case "-":
			switch right.Type() {
			case object.INTEGER:
				return newIntegerObject(-(right.(*object.Integer).Value))
			default:
				return NULL_OBJECT
			}
		}
		return NULL_OBJECT
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
				return NULL_OBJECT
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
				return NULL_OBJECT
			}
		default:
			return NULL_OBJECT
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
		return NULL_OBJECT
	default:
		return nil
	}
}

func newIntegerObject(value int64) *object.Integer {
	return &object.Integer{Value: value}
}

func newBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE_OBJECT
	}
	return FALSE_OBJECT
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
	case NULL_OBJECT:
		return false
	case FALSE_OBJECT:
		return false
	default:
		return true
	}
}
