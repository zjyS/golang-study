package main

import "fmt"

type Number interface {
	int64 | float64
}

func main() {
	c := add[int64](2, 4)
	fmt.Println(c)
}

func add[T Number](a, b T) T {
	return a + b
}
