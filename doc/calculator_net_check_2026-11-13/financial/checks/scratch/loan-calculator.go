//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.AnnuityPaymentFromPresentValue(10000000, 0.005, 120)
	b, _ := interest.FutureValue(10000000, 0.06, 10)
	c, _ := interest.PresentValue(10000000, 0.06, 10)
	fmt.Println("defaults:", a, b, c)
	j, _ := interest.PeriodicRateFromEffectiveAnnual(0.06, 12) // 6 % compounded annually, paid monthly
	d, _ := interest.AnnuityPaymentFromPresentValue(10000000, j, 120)
	e, _ := interest.FutureValue(10000000, 0.005, 120)
	f, _ := interest.PresentValue(10000000, 0.03, 20)
	fmt.Println("case2:", d, e, f)
}
