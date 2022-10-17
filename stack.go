package main

type StackConstraint interface {
	rune | float64
}

type Stack[T StackConstraint] struct {
	data []T
}

func (o *Stack[T]) push(elem T) {
	o.data = append(o.data, elem)
}

func (o *Stack[T]) peek() T {
	return o.data[len(o.data)-1]
}

func (o *Stack[T]) top() T {
	n := len(o.data) - 1
	elem := o.data[n]
	o.data = o.data[:n]
	return elem
}

func (o *Stack[T]) isEmpty() bool {
	return len(o.data) == 0
}
