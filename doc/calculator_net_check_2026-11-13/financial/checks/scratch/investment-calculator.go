//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(2000000, 0.06, 10)
	j, _ := interest.PeriodicRateFromEffectiveAnnual(0.06, 12) // FIN-INT-044
	b, _ := interest.FutureValueAnnuity(100000, j, 120)
	fmt.Printf("case1: %d + %d = %d (j=%.8f)\n", a, b, a+b, j)
	c, _ := interest.FutureValue(500000, 0.08, 15)
	d, _ := interest.FutureValueAnnuityDue(300000, 0.08, 15)
	fmt.Printf("case2: %d + %d = %d\n", c, d, c+d)
}
