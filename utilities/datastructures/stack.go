package datastructures

import "fmt"

type Stack[K any] struct {
	values []*K
}

func (stack *Stack[K]) Push(value *K) {
	stack.values = append(stack.values, value)
}

func (stack *Stack[K]) Pop() *K {
	var result *K
	if len(stack.values) > 0 {
		result = stack.values[len(stack.values)-1]
		stack.values = stack.values[:len(stack.values)-1]
	}
	return result
}
func (stack *Stack[K]) Top() *K {
	if len(stack.values) > 0 {
		return stack.values[len(stack.values)-1]
	}
	return nil
}

func (stack *Stack[K]) Size() int {
	return len(stack.values)
}

func (stack *Stack[K]) Clear() {
	stack.values = nil
}

func (stack *Stack[K]) String() string {
	result := ""
	for _, val := range stack.values {
		if len(result) > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v", *val)
	}
	return result
}
