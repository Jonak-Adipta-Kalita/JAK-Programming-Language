package evaluator

import (
	"fmt"
	"io"
	"io/ioutil"
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
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, file.GetFileName(), node.Token.Line)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, file.GetFileName(), node.Token.Line)
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
		res := evalAssignStatement(node, env, file.GetFileName(), node.Token.Line)
		if isError(res) {
			fmt.Fprintf(os.Stderr, "%s\n", res.Inspect())
			return NULL
		} else {
			return res
		}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env, file.GetFileName(), node.Token.Line)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		res := applyFunction(function, args, file.GetFileName(), node.Token.Line)
		if isError(res) {
			fmt.Fprintf(os.Stderr, "Error calling `%s` : %s\n", node.Function, res.Inspect())
			return NULL
		} else {
			return res
		}
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
		return evalIndexExpression(left, index, file.GetFileName(), node.Token.Line)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, file.GetFileName(), node.Token.Line)
	case *ast.ForLoopExpression:
		return evalForLoopExpression(node, env)
	case *ast.PostfixExpression:
		res := evalPostfixExpression(env, node.Operator, node, file.GetFileName(), node.Token.Line)
		if isError(res) {
			fmt.Fprintf(os.Stderr, "%s\n", res.Inspect())
			return NULL
		} else {
			return res
		}
	case *ast.ImportStatement:
		evalImportStatement(node, env)
	case *ast.NullLiteral:
		return NULL
	case *ast.SwitchExpression:
		return evalSwitchStatement(node, env)
	case *ast.ForeachStatement:
		return evalForeachExpression(node, env, file.GetFileName(), node.Token.Line)
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
	file string,
	line int,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: "+node.Value, file, line)
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

func evalPrefixExpression(operator string, right object.Object, file string, line int) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, file, line)
	default:
		return newError("unknown operator: %s%s", file, line, operator, right.Type())
	}
}

func evalPostfixExpression(env *object.Environment, operator string, node *ast.PostfixExpression, file string, line int) object.Object {
	switch operator {
	case "++":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError("%s is unknown", file, line, node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v + 1})
			return arg
		default:
			return newError("%s is not an int", file, line, node.Token.Literal)
		}
	case "--":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError("%s is unknown", file, line, node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v - 1})
			return arg
		default:
			return newError("%s is not an int", file, line, node.Token.Literal)
		}
	default:
		return newError("unknown operator: %s", file, line, operator)
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
	file string,
	line int,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, file, line)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, file, line)
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
			file, line, left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			file, line, left.Type(), operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, file string, line int) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", file, line, right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
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
	file string,
	line int,
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
		return newError("unknown operator: %s %s %s", file, line,
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
	file string,
	line int,
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
		return newError("unknown operator: %s %s %s", file, line,
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
	} else if ie.Elif != nil {
		for _, elifExpr := range ie.Elif {
			elifCondition := Eval(elifExpr.Condition, env)
			if isError(elifCondition) {
				return elifCondition
			}
			if isTruthy(elifCondition) {
				return Eval(elifExpr.Consequence, env)
			}
		}
	} else if ie.Else != nil {
		return Eval(ie.Else, env)
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

func newError(format, file string, line int, a ...interface{}) *object.Error {
	if file == "" {
		return &object.Error{Message: fmt.Sprintf(format, a...)}
	}
	args := append([]interface{}{file, line}, a...)
	return &object.Error{Message: fmt.Sprintf("File: %s: Line: %d: "+format, args...)}
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

func applyFunction(fn object.Object, args []object.Object, file string, line int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(file, line, args...)
	default:
		return newError("not a function: %s", file, line, fn.Type())
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

func evalIndexExpression(left, index object.Object, file string, line int) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index, file, line)
	default:
		return newError("index operator not supported: %s", file, line, left.Type())
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
	file string,
	line int,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", file, line, key.Type())
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

func evalHashIndexExpression(hash, index object.Object, file string, line int) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", file, line, index.Type())
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
	contents, err := ioutil.ReadFile(filePath)

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

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func evalAssignStatement(vs *ast.AssignStatement, env *object.Environment, file string, line int) object.Object {
	val := Eval(vs.Value, env)
	if isError(val) {
		return val
	}

	if vs.Token.Type == token.VAR {
		if _, ok := env.Get(vs.Name.Value); ok {
			return newError("Variable `%s` already defined", file, line, vs.Name.Value)
		}
	} else if vs.Token.Type == token.MUTATE {
		if _, ok := env.Get(vs.Name.Value); !ok {
			return newError("Variable `%s` not defined", file, line, vs.Name.Value)
		}
	}

	env.Set(vs.Name.Value, val)
	return NULL
}

func evalSwitchStatement(se *ast.SwitchExpression, env *object.Environment) object.Object {
	obj := Eval(se.Value, env)

	for _, opt := range se.Choices {
		if opt.Default {
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
		if opt.Default {
			out := evalBlockStatement(opt.Block, env)
			return out
		}
	}

	return nil
}

func evalForeachExpression(fle *ast.ForeachStatement, env *object.Environment, file string, line int) object.Object {
	val := Eval(fle.Value, env)

	helper, ok := val.(object.Iterable)
	if !ok {
		return newError("%s object doesn't implement the Iterable interface", file, line, val.Type())
	}

	child := object.NewEnclosedEnvironment(env)

	helper.Reset()

	ret, idx, ok := helper.Next()

	for ok {
		child.Set(fle.Ident, ret)

		idxName := fle.Index
		if idxName != "" {
			child.Set(fle.Index, idx)
		}
		Eval(fle.Body, child)
		ret, idx, ok = helper.Next()
	}

	return &object.Null{}
}
