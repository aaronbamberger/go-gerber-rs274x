package gerber_rs274x

import (
	"fmt"
)

type ExpressionStack struct {
	stack []ApertureMacroExpression
}

func NewExpressionStack() *ExpressionStack {
	stack := new(ExpressionStack)
	stack.stack = make([]ApertureMacroExpression, 0, 10) // We'll start with an initial capacity of 10
	
	return stack
}

func (stack *ExpressionStack) Push(expression ApertureMacroExpression) {
	stack.stack = append(stack.stack, expression)
}

func (stack *ExpressionStack) Pop() (expr ApertureMacroExpression, err error) {
	// Make sure the stack has at least one item on it that we can pop
	if len(stack.stack) < 1 {
		return nil,fmt.Errorf("Can't pop from an empty stack")
	}
	
	// Grab the last item in the internal slice (we pop off the end)
	expr = stack.stack[len(stack.stack) - 1]
	
	// Reslice the underlying slice to remove the last item from the stack
	stack.stack = stack.stack[:len(stack.stack) - 1]
	
	return expr,nil
}

func (stack *ExpressionStack) Pop2() (expr1 ApertureMacroExpression, expr2 ApertureMacroExpression, err error) {
	// Make sure the stack has at least two items on it that we can pop
	if len(stack.stack) < 2 {
		return nil,nil,fmt.Errorf("Can't pop 2 from stack with size %d", len(stack.stack))
	}
	
	// Grab the last two items in the internal slice (we pop off the end)
	expr1 = stack.stack[len(stack.stack) - 1]
	expr2 = stack.stack[len(stack.stack) - 2]
	
	// Reslice the underlying slice to remove the last two items from the stack
	stack.stack = stack.stack[:len(stack.stack) - 2]
	
	return expr1,expr2,nil
}

func (stack *ExpressionStack) Peek() (expr ApertureMacroExpression, err error) {
	// Make sure the stack has at least one item on it that we can pop
	if len(stack.stack) < 1 {
		return nil,fmt.Errorf("Can't peek at an empty stack")
	}
	
	// Return the last item in the slice (top of stack is the end of the slice)
	return stack.stack[len(stack.stack) - 1],nil
}

func (stack *ExpressionStack) Len() int {
	return len(stack.stack)
}