//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	j, _ := interest.PeriodicRateFromEffectiveAnnual(0.05, 12)
	for _, r := range []float64{0.05 / 12, j} {
		n, err := interest.AnnuityPeriods(80000000, 500000, r)
		fmt.Printf("rate %.6f: break-even after %.2f months = age %.2f %v\n", r, n, 65+n/12, err)
	}
}
