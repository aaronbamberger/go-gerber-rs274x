package gerber_rs274x

type ExpressionEnvironment struct {
	variables map[int]float64
}

func NewExpressionEnvironment() *ExpressionEnvironment {
	newEnv := new(ExpressionEnvironment)
	
	newEnv.variables = make(map[int]float64, 10) // We'll start with an initial capacity of 10
	
	return newEnv
}

func (env *ExpressionEnvironment) getVariableValue(variable int) float64 {
	if value,found := env.variables[variable]; found {
		return value
	} else {
		return 0.0
	}
}

func (env *ExpressionEnvironment) setVariableValue(variable int, value float64) {
	env.variables[variable] = value
}