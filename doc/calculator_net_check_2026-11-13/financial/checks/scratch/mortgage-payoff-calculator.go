//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	i := 0.06 / 12
	pmt, _ := interest.AnnuityPaymentFromPresentValue(40000000, i, 360)
	b60, _ := loan.RemainingBalance(40000000, i, 360, 60)
	n0, _ := interest.AnnuityPeriods(b60, pmt, i)
	n1, _ := interest.AnnuityPeriods(b60, pmt+50000, i)
	k, saved, err := loan.ExtraPaymentPayoff(b60, i, 300, 50000)
	fmt.Printf("form1: pmt=%d bal60=%d Norig=%.4f Nextra=%.4f | ExtraPaymentPayoff periods=%d saved=%d err=%v\n", pmt, b60, n0, n1, k, saved, err)
	for _, p := range []int64{150000, 200000} {
		n, err := interest.AnnuityPeriods(23000000, p, i)
		fmt.Printf("form2: pmt=%d N=%.4f interest≈%.2f err=%v\n", p, n, (float64(p)*n-23000000)/100, err)
	}
}
