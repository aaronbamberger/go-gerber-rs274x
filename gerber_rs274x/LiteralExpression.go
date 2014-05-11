package gerber_rs274x

import (
	"fmt"
)

type LiteralExpression struct {
	value float64
}

func (expr *LiteralExpression) EvaluateExpression(env *ExpressionEnvironment) float64 {
	return expr.value
}

func (expr *LiteralExpression) String() string {
	return fmt.Sprintf("{LiteralExpr, Value: %f}", expr.value)
}