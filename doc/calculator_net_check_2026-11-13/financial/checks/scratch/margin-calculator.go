//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/business"
)

func main() {
	m, _ := business.GrossProfitMargin(16000, 12000)
	k, _ := business.MarkupRate(12000, 16000)
	m2, _ := business.GrossProfitMargin(4999, 3725)
	k2, _ := business.MarkupRate(3725, 4999)
	fmt.Printf("case1 margin %.4f markup %.4f; case2 margin %.4f markup %.4f\n", m, k, m2, k2)
}
