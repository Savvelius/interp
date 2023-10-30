package evaluator

import (
	"fmt"

	"github.com/Savvelius/go-interp/ast"
	"github.com/Savvelius/go-interp/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		obj := Eval(node.Function, env)
		if isError(obj) {
			return obj
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(obj, args)

	case *ast.FunctionLiteral:
		body := node.Body
		params := node.Parameters
		return &object.Function{Body: body, Parameters: params, Env: object.NewEnclosedEnvironment(env)}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func applyFunction(fn object.Object, arguments []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(arguments) {
			return newError("mismatched number of arguments. expected %d got %d",
				len(fn.Parameters), len(arguments))
		}
		extendedEnv := extendFunctionEnv(fn, arguments)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(arguments...)

	default:
		return newError("not a function: %s", fn.Type())
	}

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if val, ok := builtins[node.Value]; ok {
		return val
	}

	return newError("identifier not found: %s", node.Value)
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		// don't unwrap here, unwrap in evalProgram
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		evaled := Eval(expr, env)
		if isError(evaled) {
			return []object.Object{evaled}
		}
		result = append(result, evaled)
	}
	return result
}

func evalIfExpression(node ast.Node, env *object.Environment) object.Object {
	ifNode := node.(*ast.IfExpression)
	condititon := Eval(ifNode.Condition, env)
	if isError(condititon) {
		return condititon
	}

	if isTruthy(condititon) {
		return Eval(ifNode.Consequence, env)
	}
	if ifNode.Alternative == nil {
		return NULL
	}

	return Eval(ifNode.Alternative, env)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	strl := left.(*object.String).Value
	strr := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: strl + strr}
	case "==":
		return nativeBoolToBooleanObject(strl == strr)
	case "!=":
		return nativeBoolToBooleanObject(strl != strr)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	//compare booleans by their address
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(op string, arg object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(arg)
	case "-":
		return evalMinusPrefixOperatorExpression(arg)
	default:
		return newError("unknown operator: %s%s", op, arg.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(arg object.Object) object.Object {
	switch arg {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(arg object.Object) object.Object {
	if arg.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", arg.Type())
	}

	value := arg.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
