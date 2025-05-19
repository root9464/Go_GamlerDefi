package main

import (
	"fmt"
)

func Pointer(a *int) (*int, error) {
	return a, nil
}

func main() {
	var pointer *int
	one := 1
	pointer = &one
	fmt.Printf("%+v\n", pointer)

	two := 2
	pointer, err := Pointer(&two)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", pointer)
}
