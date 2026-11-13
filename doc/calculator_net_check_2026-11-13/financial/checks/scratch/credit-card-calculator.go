//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	n, _ := interest.AnnuityPeriods(800000, 20000, 0.18/12)
	b, _ := loan.BalloonBalance(800000, 20000, 0.18/12, 61)
	p, _ := interest.AnnuityPaymentFromPresentValue(800000, 0.18/12, 24)
	fmt.Printf("N=%.3f balance after 61=%d pmt24=%d\n", n, b, p)
}
