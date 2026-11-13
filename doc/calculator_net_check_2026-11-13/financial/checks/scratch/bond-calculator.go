//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/invest"
)

func main() {
	p, err := invest.BondPrice(100000, 0.05, 0.06, 3, 1)
	fmt.Println("price", p, err)
	y, err := invest.YieldToMaturity(95000, 100000, 0.04, 20, 2)
	fmt.Printf("ytm %.5f%% %v\n", y*100, err)
}
