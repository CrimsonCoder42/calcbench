//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	n, err := interest.AnnuityPeriods(25000000, 180000, 0.07/12)
	fmt.Printf("current loan: N=%.4f err=%v\n", n, err)
	pmt, _ := interest.AnnuityPaymentFromPresentValue(25000000, 0.06/12, 240)
	apr, err := interest.AnnualPercentageRate(25000000-500000-150000, pmt, 240, 12)
	fmt.Printf("new loan: pmt=%d APR=%.5f err=%v total=%d\n", pmt, apr, err, pmt*240)
}

func init() {
	b, err := loan.BalloonBalance(25000000, 180000, 0.07/12, 285)
	fmt.Printf("balance after 285 payments of 1,800: %d err=%v\n", b, err)
}
