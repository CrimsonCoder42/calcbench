//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/invest"
)

func main() {
	h1, _ := invest.HoldingPeriodYield(100000, 200000, 0)
	a1, _ := invest.AnnualizedHoldingPeriodYield(h1, 913) // 2.5 y ≈ 913 days (whole days only)
	g1, _ := interest.GrowthRatePerPeriod(200, 400, 5)    // per half-year over 5 half-years
	e1, _ := interest.EffectiveAnnualRate(2*g1, 2)
	fmt.Printf("case1 ROI %.4f annualized(913 d) %.4f%% via half-years %.4f%%\n", h1, a1*100, e1*100)
	h2, _ := invest.HoldingPeriodYield(2500000, 2100000, 0)
	g2, _ := interest.GrowthRatePerPeriod(2500000, 2100000, 3)
	a2, _ := invest.AnnualizedHoldingPeriodYield(h2, 1095)
	fmt.Printf("case2 ROI %.4f CAGR %.4f%% annualized(1095 d) %.4f%%\n", h2, g2*100, a2*100)
}
