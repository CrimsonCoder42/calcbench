//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(2000000, 0.05, 5)
	b, _ := interest.FutureValueAnnuityDue(500000, 0.05, 5)
	real, _ := interest.PresentValue(a+b, 0.03, 5)
	fmt.Printf("case1: FV(init)=%d FVdue(contrib)=%d total=%d buyingPower=%d\n", a, b, a+b, real)
	c, _ := interest.FutureValue(1000000, 0.005, 120)
	d, _ := interest.FutureValueAnnuity(20000, 0.005, 120)
	ear, _ := interest.EffectiveAnnualRate(0.06, 12)
	fmt.Printf("case2: FV(init)=%d FVA=%d total=%d EAR=%.5f\n", c, d, c+d, ear)
}
