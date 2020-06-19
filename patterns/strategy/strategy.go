package main

import "fmt"

type Operator interface {
	Apply(int, int) int
}

type Operation struct {
	operator Operator
}

func (op *Operation) Apply(a, b int) int  {
	return op.operator.Apply(a, b)
}


type Addition struct {

}

func (add *Addition) Apply(a, b int) int {
	return a + b
}


type Mul struct {

}

func (mul *Mul) Apply(a, b int) int {
	return a - b
}

func CreateOperation(operator Operator) Operation  {
	return Operation{operator: operator}
}

func main()  {
	add := new(Addition)
	operator := CreateOperation(add)
	value := operator.Apply(1, 2)
	fmt.Println(value)
}