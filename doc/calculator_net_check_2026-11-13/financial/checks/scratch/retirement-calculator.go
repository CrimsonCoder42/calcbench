//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	j, _ := interest.PeriodicRateFromEffectiveAnnual(0.06, 12)
	for _, r := range []float64{0.005, j} {
		n, _ := interest.AnnuityPeriods(60000000, 500000, r)
		p10, _ := interest.AnnuityPaymentFromPresentValue(60000000, r, 120)
		d10, _ := interest.AnnuityDuePaymentFromPresentValue(60000000, r, 120)
		d15, _ := interest.AnnuityDuePaymentFromPresentValue(60000000, r, 180)
		fmt.Printf("rate %.8f: N(end)=%.3f pmt10(end)=%d pmt10(due)=%d pmt15(due)=%d\n", r, n, p10, d10, d15)
	}
}
