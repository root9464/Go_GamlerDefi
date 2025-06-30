package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func Pointer(a *int) (*int, error) {
	return a, nil
}

func main() {
	coin := decimal.NewFromFloat(0.2)
	fmt.Print(coin.String())

}
