//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		p, fee int64
		r      float64
		n      int
	}{{2000000, 100000, 0.10, 60}, {800000, 30000, 0.135, 36}} {
		pmt, _ := interest.AnnuityPaymentFromPresentValue(c.p, c.r/12, c.n)
		apr, _ := interest.AnnualPercentageRate(c.p-c.fee, pmt, c.n, 12)
		fmt.Printf("pmt=%d APR=%.4f%%\n", pmt, apr*100)
	}
}
