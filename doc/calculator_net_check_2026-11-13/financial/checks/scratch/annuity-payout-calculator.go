//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		pv    int64
		r     float64
		m, yr int
	}{{50000000, 0.06, 12, 10}, {25000000, 0.04, 4, 20}} {
		j, _ := interest.PeriodicRateFromEffectiveAnnual(c.r, c.m)
		e, _ := interest.AnnuityPaymentFromPresentValue(c.pv, j, c.m*c.yr)
		d, _ := interest.AnnuityDuePaymentFromPresentValue(c.pv, j, c.m*c.yr)
		n, _ := interest.AnnuityPaymentFromPresentValue(c.pv, c.r/float64(c.m), c.m*c.yr)
		nd, _ := interest.AnnuityDuePaymentFromPresentValue(c.pv, c.r/float64(c.m), c.m*c.yr)
		fmt.Printf("pv=%d eff end=%d due=%d | nominal end=%d due=%d\n", c.pv, e, d, n, nd)
	}
}
