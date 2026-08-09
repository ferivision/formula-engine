package evaluator

import (
	"strconv"

	"github.com/ferivision/formula-engine/internal/parser"
	"github.com/ferivision/formula-engine/internal/registry"
)

// Evaluate walks an AST produced by parser.Parse and computes its
// value. Field references and function calls are handled in later
// tickets (FASE-5.2, Phase 6).
func Evaluate(node *parser.Node) (registry.Value, error) {
	switch node.Type {
	case parser.NodeNumber:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, newRuntimeError("invalid number literal: " + node.Value)
		}
		return v, nil
	case parser.NodeString:
		return node.Value, nil
	case parser.NodeBinary:
		return evalBinary(node)
	case parser.NodeIdentifier:
		return nil, newRuntimeError("field references are not yet supported")
	case parser.NodeFunctionCall:
		return nil, newRuntimeError("function calls are not yet supported")
	default:
		return nil, newRuntimeError("unknown node type")
	}
}

func evalBinary(node *parser.Node) (registry.Value, error) {
	leftVal, err := Evaluate(node.Left)
	if err != nil {
		return nil, err
	}
	rightVal, err := Evaluate(node.Right)
	if err != nil {
		return nil, err
	}
	left, err := toFloat64(leftVal)
	if err != nil {
		return nil, err
	}
	right, err := toFloat64(rightVal)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case parser.TokenPlus:
		return left + right, nil
	case parser.TokenMinus:
		return left - right, nil
	case parser.TokenStar:
		return left * right, nil
	case parser.TokenSlash:
		if right == 0 {
			return nil, newRuntimeError("division by zero")
		}
		return left / right, nil
	default:
		return nil, newRuntimeError("unknown operator")
	}
}

func toFloat64(v registry.Value) (float64, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, newTypeError("expected a number")
	}
	return f, nil
}
