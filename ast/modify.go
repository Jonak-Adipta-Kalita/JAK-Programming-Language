package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Elif != nil {
			for i := range node.Elif {
				node.Elif[i].Condition, _ = Modify(node.Elif[i].Condition, modifier).(Expression)
				node.Elif[i].Consequence, _ = Modify(node.Elif[i].Consequence, modifier).(*BlockStatement)
			}
		}
		if node.Else != nil {
			node.Else, _ = Modify(node.Else, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i], _ = Modify(node.Statements[i], modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *AssignStatement:
		node.Name, _ = Modify(node.Name, modifier).(*Identifier)
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, val := range node.Pairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newVal, _ := Modify(val, modifier).(Expression)
			newPairs[newKey] = newVal
		}
		node.Pairs = newPairs
	case *ForLoopExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
	case *PostfixExpression:
		node.Operator, _ = Modify(node.Operator, modifier).(*StringLiteral)
	case *ImportStatement:
		node.Path, _ = Modify(node.Path, modifier).(*StringLiteral)
	case *CaseExpression:
		node.Default, _ = Modify(node.Default, modifier).(*Boolean)
		node.Expr = Modify(node.Expr, modifier).(Expression)
		node.Block = Modify(node.Block, modifier).(*BlockStatement)
	case *SwitchExpression:
		node.Value = Modify(node.Value, modifier).(Expression)
		for i := range node.Choices {
			node.Choices[i] = Modify(node.Choices[i], modifier).(*CaseExpression)
		}
	case *ForeachStatement:
		node.Index = Modify(node.Index, modifier).(*StringLiteral)
		node.Identifier = Modify(node.Identifier, modifier).(*StringLiteral)
		node.Value = Modify(node.Value, modifier).(Expression)
		node.Body = Modify(node.Body, modifier).(*BlockStatement)
	case *ObjectCallExpression:
		node.Object = Modify(node.Object, modifier).(Expression)
		node.Call = Modify(node.Call, modifier).(Expression)
	}
	return modifier(node)
}
