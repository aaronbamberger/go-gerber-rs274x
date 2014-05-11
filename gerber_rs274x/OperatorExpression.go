package gerber_rs274x

import (
	"fmt"
)

type OperatorExpression struct {
	operator ArithmeticOperator
}

func (expr *OperatorExpression) EvaluateExpression(env *ExpressionEnvironment) float64 {
	return 0.0
}

func (expr *OperatorExpression) String() string {
	var operator string
	switch expr.operator {
		case OPERATOR_ADD:
			operator = "Add"
		
		case OPERATOR_SUBTRACT:
			operator = "Subtract"
		
		case OPERATOR_MULTIPLY:
			operator = "Multiply"
		
		case OPERATOR_DIVIDE:
			operator = "Divide"
	}
	
	return fmt.Sprintf("{OperatorExpr, Operator: %s}", operator)
}