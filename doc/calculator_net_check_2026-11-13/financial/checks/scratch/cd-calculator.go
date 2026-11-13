//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(1000000, 0.05, 3)
	b, _ := interest.FutureValue(2500000, 0.045/12, 18)
	e1, _ := interest.EffectiveAnnualRate(0.045, 12)
	c, _ := interest.FutureValueContinuous(1000000, 0.04, 2)
	e2, _ := interest.EffectiveAnnualRateContinuous(0.04)
	fmt.Printf("a=%d b=%d ear=%.5f c=%d earc=%.5f\n", a, b, e1*100, c, e2*100)
}
