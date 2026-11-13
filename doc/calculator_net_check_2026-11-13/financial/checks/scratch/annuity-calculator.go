//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(2000000, 0.06, 10)
	b, _ := interest.FutureValueAnnuityDue(1000000, 0.06, 10)
	j, _ := interest.PeriodicRateFromEffectiveAnnual(0.07, 12)
	c, _ := interest.FutureValueAnnuity(30000, j, 240)
	d, _ := interest.FutureValueAnnuity(30000, 0.07/12, 240)
	fmt.Println(a, b, a+b, "| monthly eff", c, "nominal", d)
}
