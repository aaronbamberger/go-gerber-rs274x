package gerber_rs274x

type Modifier interface {
	GetModifierValue() float64
}

type LiteralModifier struct {
	value float64
}

type ExpressionModifier struct {
	expression string
}

type VariableModifier struct {
	variableName string
}

func (literal* LiteralModifier) GetModifierValue() float64 {
	return literal.value
}

func (expression* ExpressionModifier) GetModifierValue() float64 {
	//TODO: Implement
	return 0.0
}

func (variable* VariableModifier) GetModifierValue() float64 {
	//TODO: Implement
	return 0.0
}