//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		p, fees int64
		r       float64
		n       int
	}{{1000000, 50000 + 75000, 0.10, 60}, {15000000, 300000 + 50000 + 25000, 0.085, 84}} {
		pmt, _ := interest.AnnuityPaymentFromPresentValue(c.p, c.r/12, c.n)
		apr, _ := interest.AnnualPercentageRate(c.p-c.fees, pmt, c.n, 12)
		fmt.Printf("pmt=%d APR=%.4f%%\n", pmt, apr*100)
	}
}
