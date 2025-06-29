package main

import (
	"fmt"
)

func Pointer(a *int) (*int, error) {
	return a, nil
}

func main() {
	var coin uint64 = 0
	var float_coin = 0.2

	coin = uint64(float_coin)
	fmt.Print(coin)
}
