package evaluator

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/ast"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/file"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/lexer"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/parser"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		res := evalProgram(node, env)
		if isError(res) {
			fmt.Fprintf(os.Stderr, "%s\n", res.Inspect())
			return NULL
		} else {
			return res
		}
	case *ast.ExpressionStatement:
		res := Eval(node.Expression, env)
		if isError(res) {
			fmt.Fprintf(os.Stderr, "%s\n", res.Inspect())
			return NULL
		} else {
			return res
		}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.AssignStatement:
		return evalAssignStatement(node, env, node.Token)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env, node.Token)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		if node.Function.TokenLiteral() == "quote" {
			if len(node.Arguments) != 1 {
				return newError(
					"wrong number of arguments. got=%d, want=1",
					node.Token,
					len(node.Arguments),
				)
			}
			return quote(node.Arguments[0], env)
		}

		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, node.Token)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, node.Token)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, node.Token)
	case *ast.ForLoopExpression:
		return evalForLoopExpression(node, env)
	case *ast.PostfixExpression:
		return evalPostfixExpression(env, node.Operator.Value, node, node.Token)
	case *ast.ImportStatement:
		evalImportStatement(node, env)
	case *ast.NullLiteral:
		return NULL
	case *ast.SwitchExpression:
		return evalSwitchStatement(node, env)
	case *ast.ForeachStatement:
		return evalForeachExpression(node, env, node.Token)
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.ObjectCallExpression:
		res := evalObjectCallExpression(node, env, node.Token)
		return res
	}
	return nil
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
	token token.Token,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: "+node.Value, token)
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
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

func evalPrefixExpression(operator string, right object.Object, token token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, token)
	default:
		return newError("unknown operator: %s %s", token, operator, right.Type())
	}
}

func evalPostfixExpression(
	env *object.Environment,
	operator string,
	node *ast.PostfixExpression,
	token token.Token,
) object.Object {
	switch operator {
	case "++":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError("%s is unknown", token, node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v + 1})
			return arg
		default:
			return newError("%s is not an int", token, node.Token.Literal)
		}
	case "--":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError("%s is unknown", token, node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v - 1})
			return arg
		default:
			return newError("%s is not an int", token, node.Token.Literal)
		}
	default:
		return newError("unknown operator: %s", token, operator)
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
	token token.Token,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, token)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, token)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatIntegerInfixExpression(operator, left, right, token)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalIntegerFloatInfixExpression(operator, left, right, token)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, token)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case operator == "&&":
		return nativeBoolToBooleanObject(coerceObjectToNativeBool(left) && coerceObjectToNativeBool(right))
	case operator == "||":
		return nativeBoolToBooleanObject(coerceObjectToNativeBool(left) || coerceObjectToNativeBool(right))
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			token, left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			token, left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
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
}

func evalMinusPrefixOperatorExpression(right object.Object, token token.Token) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return newError("unknown operator: -%s", token, right.Type())
	}
}

func coerceObjectToNativeBool(o object.Object) bool {
	if rv, ok := o.(*object.ReturnValue); ok {
		o = rv.Value
	}

	switch obj := o.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return obj.Value != ""
	case *object.Null:
		return false
	case *object.Integer:
		return obj.Value != 0
	case *object.Array:
		return len(obj.Elements) > 0
	case *object.Hash:
		return len(obj.Pairs) > 0
	default:
		return true
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
	token token.Token,
) object.Object {
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
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "^":
		return &object.Integer{Value: powInt(leftVal, rightVal)}
	default:
		return newError("unknown operator: %s %s %s", token,
			left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, token token.Token) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError(
			"unknown operator: %s %s %s",
			token,
			left.Type(),
			operator,
			right.Type(),
		)
	}
}

func evalFloatIntegerInfixExpression(operator string, left, right object.Object, token token.Token) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := float64(right.(*object.Integer).Value)
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError(
			"unknown operator: %s %s %s",
			token,
			left.Type(),
			operator,
			right.Type(),
		)
	}
}

func evalIntegerFloatInfixExpression(operator string, left, right object.Object, token token.Token) object.Object {
	leftVal := float64(left.(*object.Integer).Value)
	rightVal := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError(
			"unknown operator: %s %s %s",
			token,
			left.Type(),
			operator,
			right.Type(),
		)
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
	token token.Token,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", token,
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else {
		if ie.Elif != nil {
			for _, elifExpr := range ie.Elif {
				elifCondition := Eval(elifExpr.Condition, env)
				if isError(elifCondition) {
					return elifCondition
				}
				if isTruthy(elifCondition) {
					return Eval(elifExpr.Consequence, env)
				}
			}
		}
		if ie.Else != nil {
			return Eval(ie.Else, env)
		}
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, token token.Token, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...), FileName: file.GetFileName(), Token: token}
}

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}

func powInt(x, y int64) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func applyFunction(fn object.Object, args []object.Object, token token.Token) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(token, args...)
	default:
		return newError("not a function: %s", token, fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
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

func evalIndexExpression(left, index object.Object, token token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index, token)
	default:
		return newError("index operator not supported: %s", token, left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
	token token.Token,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", token, key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object, token token.Token) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", token, index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalForLoopExpression(fle *ast.ForLoopExpression, env *object.Environment) object.Object {
	var rt object.Object
	for {
		condition := Eval(fle.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			rt := Eval(fle.Consequence, env)
			if !isError(rt) && (rt.Type() == object.RETURN_VALUE_OBJ || rt.Type() == object.ERROR_OBJ) {
				return rt
			}
		} else {
			break
		}
	}
	return rt
}

func evalImportStatement(is *ast.ImportStatement, env *object.Environment) {
	filePath := is.Path.Value
	file.SetFileName(filePath)
	contents, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("Failure to read file '%s'. Err: %s", string(contents), err)
		return
	}

	l := lexer.New(string(contents))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		PrintParserErrors(os.Stdout, p.Errors())
		return
	}

	Eval(program, env)
	file.SetFileName(file.GetMainFileName())
}

func evalAssignStatement(vs *ast.AssignStatement, env *object.Environment, token_ token.Token) object.Object {
	val := Eval(vs.Value, env)
	if isError(val) {
		return val
	}

	if vs.Token.Type == token.VAR {
		if _, ok := env.Get(vs.Name.Value); ok {
			return newError("Variable `%s` already defined", token_, vs.Name.Value)
		}
	} else if vs.Token.Type == token.MUTATE {
		if _, ok := env.Get(vs.Name.Value); !ok {
			return newError("Variable `%s` not defined", token_, vs.Name.Value)
		}
	}

	env.Set(vs.Name.Value, val)
	return NULL
}

func evalSwitchStatement(se *ast.SwitchExpression, env *object.Environment) object.Object {
	obj := Eval(se.Value, env)

	for _, opt := range se.Choices {
		if opt.Default.Token.Type == token.TRUE {
			continue
		}

		val := Eval(opt.Expr, env)

		if obj.Type() == val.Type() &&
			(obj.Inspect() == val.Inspect()) {

			out := evalBlockStatement(opt.Block, env)
			return out
		}
	}

	for _, opt := range se.Choices {
		if opt.Default.Token.Type == token.TRUE {
			out := evalBlockStatement(opt.Block, env)
			return out
		}
	}

	return nil
}

func evalForeachExpression(fle *ast.ForeachStatement, env *object.Environment, token token.Token) object.Object {
	val := Eval(fle.Value, env)

	helper, ok := val.(object.Iterable)
	if !ok {
		return newError("%s object doesn't implement the Iterable interface", token, val.Type())
	}

	child := object.NewEnclosedEnvironment(env)
	helper.Reset()
	ret, idx, ok := helper.Next()

	for ok {
		child.Set(fle.Identifier.Value, ret)
		if fle.Index != nil && fle.Index.Value != "" {
			child.Set(fle.Index.Value, idx)
		}

		rt := Eval(fle.Body, child)
		if !isError(rt) && (rt.Type() == object.RETURN_VALUE_OBJ || rt.Type() == object.ERROR_OBJ) {
			return rt
		}

		ret, idx, ok = helper.Next()
	}

	return &object.Null{}
}

func evalObjectCallExpression(call *ast.ObjectCallExpression, env *object.Environment, token token.Token) object.Object {
	obj := Eval(call.Object, env)
	if method, ok := call.Call.(*ast.CallExpression); ok {
		args := evalExpressions(call.Call.(*ast.CallExpression).Arguments, env)
		ret := obj.InvokeMethod(method.Function.String(), args...)
		if ret != nil {
			return ret
		}
	}

	return newError("Failed to invoke method: %s", token, call.Call.(*ast.CallExpression).Function.String())
}
