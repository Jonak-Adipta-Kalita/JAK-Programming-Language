package evaluator

import (
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/ast"
	"github.com/Jonak-Adipta-Kalita/JAK-Programming-Language/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
