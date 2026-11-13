//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func main() {
	fmt.Println(algebra.PercentOf(70, 3500))
	fmt.Println(algebra.PercentFromPart(155, 228))
	fmt.Println(algebra.TotalFromPercent(4000, 10))
	fmt.Println(algebra.PercentDifference(40, 50))
	fmt.Println(algebra.ApplyPercentChange(500, -12.5))
	fmt.Println(algebra.PercentChange(40, 50))
}
