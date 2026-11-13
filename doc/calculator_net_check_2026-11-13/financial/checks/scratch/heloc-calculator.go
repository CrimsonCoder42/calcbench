//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		p           int64
		rate        float64
		draw, repay int
	}{{5000000, 0.08, 5, 15}, {12000000, 0.0925, 10, 20}} {
		io, _ := interest.SimpleInterest(c.p, c.rate, 1.0/12) // interest-only draw payment
		rp, _ := interest.AnnuityPaymentFromPresentValue(c.p, c.rate/12, c.repay*12)
		fmt.Printf("p=%d draw pmt=%d repay pmt=%d total≈%d\n", c.p, io, rp, io*int64(c.draw*12)+rp*int64(c.repay*12))
	}
}
