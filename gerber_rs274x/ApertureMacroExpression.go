package gerber_rs274x

import (
	"fmt"
	"unicode"
	"strconv"
)

type ExpressionParseState int

const (
	PARSING_NORMAL ExpressionParseState = iota
	PARSING_LITERAL
	PARSING_VARIABLE
)

type ApertureMacroExpression interface {
	EvaluateExpression(env *ExpressionEnvironment) float64
}

func parseExpression(infixExpression string) (ApertureMacroExpression, error) {
	//NOTE: This function uses the expression parsing algorithms found here: http://www.smccd.net/accounts/hasson/C++2Notes/ArithmeticParsing.html

	// First, we parse the string into list of expressions in postfix order
	postfixExpressionList := make([]ApertureMacroExpression, 0, 10) // Start with an initial capacity of 10
	stack := NewExpressionStack()
	state := PARSING_NORMAL
	numberAccumulator := ""
	
	for _,char := range infixExpression {
		switch state {
			case PARSING_NORMAL:
				if newState,err := nextParseState(char, &numberAccumulator, &postfixExpressionList, stack); err != nil {
					return nil,err
				} else {
					state = newState
				}
			
			case PARSING_LITERAL:
				if unicode.IsDigit(char) || char == '.' {
					numberAccumulator += string(char)
				} else {
					if literalVal,err := strconv.ParseFloat(numberAccumulator, 64); err != nil {
						return nil,fmt.Errorf("Error parsing literal: %s\n", err.Error())
					} else {
						postfixExpressionList = append(postfixExpressionList, &LiteralExpression{literalVal})
					}
					
					// Reset the string accumulator
					numberAccumulator = ""
					
					if newState,err := nextParseState(char, &numberAccumulator, &postfixExpressionList, stack); err != nil {
						return nil,err
					} else {
						state = newState
					}
				}
			
			case PARSING_VARIABLE:
				if unicode.IsDigit(char) {
					numberAccumulator += string(char)
				} else {
					if varNum,err := strconv.ParseInt(numberAccumulator, 10, 32); err != nil {
						return nil,fmt.Errorf("Error parsing variable number: %s\n", err.Error())
					} else {
						postfixExpressionList = append(postfixExpressionList, &VariableExpression{int(varNum)})
					}
					
					// Reset the string accumulator
					numberAccumulator = ""
					
					if newState,err := nextParseState(char, &numberAccumulator, &postfixExpressionList, stack); err != nil {
						return nil,err
					} else {
						state = newState
					}
				}
		}
	}
	
	// If the last thing we parsed from the string was a literal or a variable, then we need to finalize the parsing
	if state == PARSING_LITERAL {
		if literalVal,err := strconv.ParseFloat(numberAccumulator, 64); err != nil {
			return nil,fmt.Errorf("Error parsing literal: %s\n", err.Error())
		} else {
			postfixExpressionList = append(postfixExpressionList, &LiteralExpression{literalVal})
		}
	} else if state == PARSING_VARIABLE {
		if varNum,err := strconv.ParseInt(numberAccumulator, 10, 32); err != nil {
			return nil,fmt.Errorf("Error parsing variable number: %s\n", err.Error())
		} else {
			postfixExpressionList = append(postfixExpressionList, &VariableExpression{int(varNum)})
		}
	}
	
	// Empty the stack onto the end of the list
	for stack.Len() > 0 {
		// We've already checked that the stack isn't empty, so we can pop the stack and ignore the error
		expr,_ := stack.Pop()
		postfixExpressionList = append(postfixExpressionList, expr)
	}
	
	// We then parse the postfix expression list into an expression tree that we can return
	// We've now emptied the stack, so we can reuse it
	for _,expr := range postfixExpressionList {
		switch exprValue := expr.(type) {
			case *VariableExpression,*LiteralExpression:
				stack.Push(expr)
			
			case *OperatorExpression:
				if expr1,expr2,err := stack.Pop2(); err != nil {
					return nil,err
				} else {
					stack.Push(&ArithmeticExpression{exprValue.operator, expr2, expr1})
				}
			
			default:
				return nil,fmt.Errorf("Unexpected type encountered during parse: %T\n", expr)
		}
	}
	
	// Once we're done, there should only be 1 expression left on the stack, which is the root of the expression tree
	if stack.Len() != 1 {
		return nil,fmt.Errorf("Incorrect number of elements left on postfix stack after parsing: %d\n", stack.Len()) 
	}
	
	// If we're here, we know there is only 1 expression left on the stack, so we pop it and return it
	expr,_ := stack.Pop()
	return expr,nil
}

func nextParseState(char rune, numberAccumulator *string, expressionList *[]ApertureMacroExpression, stack *ExpressionStack) (ExpressionParseState, error) {
	if char == '$' {
		return PARSING_VARIABLE,nil
		
	} else if unicode.IsDigit(char) || char == '.' {
		*numberAccumulator = *numberAccumulator + string(char)
		return PARSING_LITERAL,nil
		
	} else if char == '+' || char == '-' || char == 'x' || char == '/' {
		switch char {
			case '+', '-':
				done := false
				for (stack.Len() > 0) && !done {
					if topOfStack,err := stack.Peek(); err != nil {
						return PARSING_NORMAL,err
					} else {
						switch exprValue := topOfStack.(type) {
							case *OperatorExpression:
								if exprValue.operator == OPERATOR_MULTIPLY || exprValue.operator == OPERATOR_DIVIDE {
									// We've already checked that there are items on the stack, so it's safe to pop and ignore the error
									expr,_ := stack.Pop()
									*expressionList = append(*expressionList, expr)
								} else {
									done = true
								}
							
							default:
								done = true
						}
					}
				}
				fallthrough // We always want to push this new operator on the stack, so we fallthrough here
			
			case 'x', '/':
				switch char {
					case '+':
						stack.Push(&OperatorExpression{OPERATOR_ADD})
						
					case '-':
						stack.Push(&OperatorExpression{OPERATOR_SUBTRACT})
					
					case 'x':
						stack.Push(&OperatorExpression{OPERATOR_MULTIPLY})
					
					case '/':
						stack.Push(&OperatorExpression{OPERATOR_DIVIDE})
				}
		}
		
		return PARSING_NORMAL,nil
		
	} else if char == '(' || char == ')' {
		switch char {
			case '(':
				stack.Push(&ParenthesisExpression{LEFT_PARENTHESIS})
				
			case ')':
				done := false
				for !done {
					if topOfStack,err := stack.Peek(); err != nil {
						return PARSING_NORMAL,err
					} else {
						switch exprValue := topOfStack.(type) {
							case *ParenthesisExpression:
								if exprValue.parenType == LEFT_PARENTHESIS {
									// Pop the top of the stack, and throw it away
									if _,err := stack.Pop(); err != nil {
										return PARSING_NORMAL,err
									}
									done = true
								} else {
									return PARSING_NORMAL,fmt.Errorf("Encountered unmatched parentheses")
								}
							
							default:
								// If we're here, we already know there's something on the stack to be popped,
								// so we can pop and ignore the error
								expr,_ := stack.Pop()
								*expressionList = append(*expressionList, expr);
						}
					}
				}
		}
		
		return PARSING_NORMAL,nil
		
	} else {
		return PARSING_NORMAL,fmt.Errorf("Unexpected character encountered: %c", char)
		
	}
}