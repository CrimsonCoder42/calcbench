//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	for _, c := range []struct {
		name string
		p    int64
		r    float64
		n    int
	}{{"cash back", 3900000, 0.05, 60}, {"low rate", 4000000, 0.02, 60}, {"case2 cash back", 2700000, 0.04, 36}, {"case2 low", 3000000, 0.029, 36}} {
		pmt, _ := interest.AnnuityPaymentFromPresentValue(c.p, c.r/12, c.n)
		rows, _ := loan.AmortizationSchedule(c.p, c.r/12, c.n)
		var ti int64
		for _, r := range rows {
			ti += r.InterestCents
		}
		fmt.Printf("%s: pmt=%d interest=%d\n", c.name, pmt, ti)
	}
}
