//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	p, _ := interest.AnnuityPaymentFromPresentValue(1000000, 0.10/12, 60)
	n, _ := interest.AnnuityPeriods(1000000, 20000, 0.10/12)
	fmt.Printf("pmt=%d N(200/mo)=%.3f\n", p, n)
}
