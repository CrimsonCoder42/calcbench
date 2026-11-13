//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/business"
)

func main() {
	a, _ := business.SalePriceAfterDiscount(5999, 0.15)
	b, _ := business.SalePriceAfterDiscount(7580, 0.25)
	fmt.Println(a, 5999-a, b, 7580-b)
}
