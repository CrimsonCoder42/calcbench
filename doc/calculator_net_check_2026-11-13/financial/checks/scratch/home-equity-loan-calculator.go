//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	pmt, _ := interest.AnnuityPaymentFromPresentValue(15000000, 0.08/12, 180)
	apr, err := interest.AnnualPercentageRate(15000000-750000, pmt, 180, 12)
	fmt.Printf("pmt=%d APR=%.5f%% err=%v\n", pmt, apr*100, err)
	ltv, _ := loan.LoanToValueRatio(25000000, 60000000)
	cltv, _ := loan.LoanToValueRatio(25000000+26000000, 60000000)
	fmt.Printf("LTV=%.4f CLTV with 260k second lien=%.4f\n", ltv, cltv)
}
