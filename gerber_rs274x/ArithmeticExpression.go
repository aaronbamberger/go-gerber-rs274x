package gerber_rs274x

import (
	"fmt"
)

type ArithmeticOperator int

const (
	OPERATOR_ADD ArithmeticOperator = iota
	OPERATOR_SUBTRACT
	OPERATOR_MULTIPLY
	OPERATOR_DIVIDE
)

type ArithmeticExpression struct {
	operator ArithmeticOperator
	lhs ApertureMacroExpression
	rhs ApertureMacroExpression
}

func (expr *ArithmeticExpression) EvaluateExpression(env *ExpressionEnvironment) float64 {
	switch expr.operator {
		case OPERATOR_ADD:
			return expr.lhs.EvaluateExpression(env) + expr.rhs.EvaluateExpression(env)
		
		case OPERATOR_SUBTRACT:
			return expr.lhs.EvaluateExpression(env) - expr.rhs.EvaluateExpression(env)
		
		case OPERATOR_MULTIPLY:
			return expr.lhs.EvaluateExpression(env) * expr.rhs.EvaluateExpression(env)
		
		case OPERATOR_DIVIDE:
			return expr.lhs.EvaluateExpression(env) / expr.rhs.EvaluateExpression(env)
		
		default:
			return 0.0
	}
}

func (expr *ArithmeticExpression) String() string {
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
	
	return fmt.Sprintf("{ArithmeticExpr, Operator: %s, LHS: %v, RHS: %v}\n", operator, expr.lhs, expr.rhs)
}