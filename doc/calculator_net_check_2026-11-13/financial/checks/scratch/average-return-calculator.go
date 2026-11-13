//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/invest"
)

func main() {
	segs := []struct {
		ratePerYear float64
		months      int
	}{{0.10, 14}, {-0.02, 5}, {0.15, 27}}
	var hpr []float64
	bal := int64(100000000) // $1,000,000 in cents, so cent rounding stays below their 2 decimals
	total := 0
	for _, s := range segs {
		m, _ := interest.PeriodicRateFromEffectiveAnnual(s.ratePerYear, 12)
		fv, _ := interest.FutureValue(bal, m, s.months)
		hpr = append(hpr, float64(fv)/float64(bal)-1)
		bal = fv
		total += s.months
	}
	cum, _ := invest.CumulativeReturn(hpr)
	monthly, _ := interest.GrowthRatePerPeriod(100000000, bal, total)
	avg, _ := interest.EffectiveAnnualRate(12*monthly, 12)
	fmt.Printf("balances end %d cumulative %.4f%% average %.4f%%/yr\n", bal, cum*100, avg*100)
}
