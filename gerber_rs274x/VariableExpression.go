package gerber_rs274x

import (
	"fmt"
)

type VariableExpression struct {
	variableNumber int
}

func (expr *VariableExpression) EvaluateExpression(env *ExpressionEnvironment) float64 {
	return env.getVariableValue(expr.variableNumber)
}

func (expr *VariableExpression) String() string {
	return fmt.Sprintf("{VariableExpr, Variable Number: %d}", expr.variableNumber)
}