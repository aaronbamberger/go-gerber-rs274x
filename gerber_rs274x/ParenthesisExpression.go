package gerber_rs274x

import (
	"fmt"
)

type ParenthesisType int

const (
	LEFT_PARENTHESIS ParenthesisType = iota
	RIGHT_PARENTHESIS
)

type ParenthesisExpression struct {
	parenType ParenthesisType
}

func (expr *ParenthesisExpression) EvaluateExpression(env *ExpressionEnvironment) float64 {
	return 0.0
}

func (expr *ParenthesisExpression) String() string {
	var exprType string
	switch expr.parenType {
		case LEFT_PARENTHESIS:
			exprType = "Left"
		
		case RIGHT_PARENTHESIS:
			exprType = "Right"
	}

	return fmt.Sprintf("{ParenthesisExpr, Type: %s}", exprType)
}