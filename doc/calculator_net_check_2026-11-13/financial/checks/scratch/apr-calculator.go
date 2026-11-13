//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		name               string
		principal, charges int64 // note amount (incl. financed fees), prepaid finance charges
		rate               float64
		n                  int
	}{
		{"100k 10y 6% upfront 2,500", 10000000, 250000, 0.06, 120},
		{"25k+750 financed, 5y 8.5%, upfront 300", 2575000, 75000 + 30000, 0.085, 60},
		{"280k 30y 6.2%, fees 3,500 + 0.5 pt", 28000000, 350000 + 140000, 0.062, 360},
	} {
		pmt, _ := interest.AnnuityPaymentFromPresentValue(c.principal, c.rate/12, c.n)
		apr, err := interest.AnnualPercentageRate(c.principal-c.charges, pmt, c.n, 12)
		fmt.Printf("%s: pmt=%d APR=%.5f%% err=%v\n", c.name, pmt, apr*100, err)
	}
}
