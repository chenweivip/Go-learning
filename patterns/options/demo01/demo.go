package main

import "fmt"

func add(nums...int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

type MyFuncOptions struct {
	optionStr1 string
	optionStr2 string
	optionInt1 int
	optionInt2 int
}

var defalutMyFuncOptions = MyFuncOptions{
	optionStr1: "defaultStr1",
	optionStr2: "defaultStr2",
	optionInt1: 1,
	optionInt2: 2,
}

type MyFuncOption func(options *MyFuncOptions)

func WithOptionStr1(str1 string) MyFuncOption {
	return func(options *MyFuncOptions) {
		options.optionStr1 = str1
	}
}

func WithOptionInt1(int1 int) MyFuncOption  {
	return func(options *MyFuncOptions) {
		options.optionInt1 = int1
	}
}

func WithOptionStr2AndInt2(str2 string, int2 int) MyFuncOption {
	return func(options *MyFuncOptions) {
		options.optionStr2 = str2
		options.optionInt2 = int2
	}
}

func MyFunc(required string, opts ...MyFuncOption)  {
	options := defalutMyFuncOptions
	for _, o := range opts{
		o(&options)
	}
	fmt.Println(required, options.optionStr1, options.optionInt1, options.optionInt2, options.optionStr2)
}

func main() {
	MyFunc("hello", WithOptionInt1(1), WithOptionStr1("world"), WithOptionStr2AndInt2("age", 2))
}
