package eval

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL_OBJECT  = &object.Null{}
	TRUE_OBJECT  = &object.Boolean{Value: true}
	FALSE_OBJECT = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement, env)
			switch result := result.(type) {
			case *object.Return:
				return result.Value
			case *object.Error:
				return result
			}
		}
		return result
	case *ast.BlockStatement:
		var result object.Object
		for _, statement := range node.Statements {
			result = Eval(statement, env)
			switch result := result.(type) {
			case *object.Return:
				return result
			case *object.Error:
				return result
			}
		}
		return result
	case *ast.ReturnStatement:
		val := &object.Return{Value: Eval(node.ValueExpression, env)}
		if isError(val) {
			return val
		}
		return val
	case *ast.LetStatement:
		val := Eval(node.ValueExpression, env)
		if isError(val) {
			return val
		}
		env.Set(node.IdentifierExpression.Value, val)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IdentifierExpression:
		val, ok := env.Get(node.Value)
		if !ok {
			return newErrorObject("identifier not found: %s", node.Value)
		}
		return val
	case *ast.FunctionExpression:
		val := &object.Function{
			Parameters: node.ParameterExpressions,
			Body:       node.BodyStatement,
			Env:        env,
		}
		return val
	case *ast.CallExpression:
		var args []object.Object
		for _, argumentExpression := range node.ArgumentExpressions {
			val := Eval(argumentExpression, env)
			if isError(val) {
				return val
			}
			args = append(args, val)
		}
		function := Eval(node.FunctionExpression, env)
		if isError(function) {
			return function
		}
		functionObject, ok := function.(*object.Function)
		if !ok {
			return newErrorObject("identifier is not function: %s", node.FunctionExpression.String())
		}
		for i, arg := range args {
			name := functionObject.Parameters[i]
			env.Set(name.Value, arg)
		}
		return Eval(functionObject.Body, env)
	case *ast.IntegerExpression:
		return newIntegerObject(node.Value)
	case *ast.BooleanExpression:
		return newBooleanObject(node.Value)
	case *ast.UnaryExpression:
		right := Eval(node.RightExpression, env)
		if isError(right) {
			return right
		}
		switch node.Operator {
		case "!":
			return newBooleanObject(!isTruthy(right))
		case "-":
			switch right.Type() {
			case object.INTEGER:
				return newIntegerObject(-(right.(*object.Integer).Value))
			default:
				return newErrorObject("unknown operator: %s%s", node.Operator, right.Type())
			}
		}
		return newErrorObject("unknown operator: %s%s", node.Operator, right.Type())
	case *ast.BinaryExpression:
		left := Eval(node.LeftExpression, env)
		if isError(left) {
			return left
		}
		right := Eval(node.RightExpression, env)
		if isError(right) {
			return right
		}
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
			}
		case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
			leftValue := left.(*object.Boolean).Value
			rightValue := right.(*object.Boolean).Value
			switch node.Operator {
			case "==":
				return newBooleanObject(leftValue == rightValue)
			case "!=":
				return newBooleanObject(leftValue != rightValue)
			}
		case left.Type() != right.Type():
			return newErrorObject("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
		}
		return newErrorObject("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	case *ast.IfExpression:
		condition := Eval(node.ConditionExpression, env)
		if isError(condition) {
			return condition
		}
		if isTruthy(condition) {
			return Eval(node.ConsequenceStatement, env)
		} else {
			if node.AlternativeStatement != nil {
				return Eval(node.AlternativeStatement, env)
			}
		}
	}
	return NULL_OBJECT
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

func newErrorObject(format string, values ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, values...)}
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

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR
}
