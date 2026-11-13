//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	i := 0.068 / 12
	p, _ := interest.AnnuityPaymentFromPresentValue(3000000, i, 120)
	n0, _ := interest.AnnuityPeriods(3000000, 35000, i)
	n1, _ := interest.AnnuityPeriods(3000000, 50000, i)
	fmt.Printf("form1 pmt=%d; form2 N(350)=%.3f N(500)=%.3f\n", p, n0, n1)
	g, _ := interest.FutureValue(4426399, i, 6) // their balance after graduation, 6-month grace, monthly compounding
	q, _ := interest.AnnuityPaymentFromPresentValue(4579044, i, 120)
	fmt.Printf("form3 grace %d; payment on their post-grace balance %d\n", g, q)
}
