//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	p, _ := interest.AnnuityPaymentFromPresentValue(20000000, 0.005, 180)
	n, _ := interest.AnnuityPeriods(20000000, 200000, 0.005)
	fmt.Printf("fixed term pmt=%d; fixed payment 2,000 → N=%.3f months\n", p, n)
}
