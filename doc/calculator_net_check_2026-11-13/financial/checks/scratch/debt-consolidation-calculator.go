//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	// consolidation loan: payment 543.44 on 25,000 over 60 months → rate; APR on 23,750 financed
	r, _ := interest.AnnuityRate(2500000, 54344, 0, 60)
	apr, _ := interest.AnnualPercentageRate(2500000-125000, 54344, 60, 12)
	fmt.Printf("loan rate %.4f%%/yr, APR %.4f%%\n", r*1200, apr*100)
	debts := []struct {
		bal, pmt int64
		rate     float64
	}{{1000000, 26000, 0.1799}, {750000, 19000, 0.1999}, {650000, 18000, 0.1899}}
	flows := make([]int64, 1, 80)
	flows[0] = -2400000
	var total int64
	for _, d := range debts {
		i := d.rate / 12
		n, _ := interest.AnnuityPeriods(d.bal, d.pmt, i)
		full := int(math.Floor(n))
		left, _ := loan.BalloonBalance(d.bal, d.pmt, i, full)
		last := int64(math.RoundToEven(float64(left) * (1 + i)))
		for len(flows) < full+2 {
			flows = append(flows, 0)
		}
		for k := 1; k <= full; k++ {
			flows[k] += d.pmt
		}
		if last > 0 {
			flows[full+1] += last
		}
		total += int64(full)*d.pmt + last
		fmt.Printf("debt %d: N=%.3f last=%d\n", d.bal, n, last)
	}
	for flows[len(flows)-1] == 0 {
		flows = flows[:len(flows)-1]
	}
	irr, err := interest.InternalRateOfReturn(flows)
	fmt.Printf("months=%d total=%d combined APR %.4f%% %v\n", len(flows)-1, total, irr*1200, err)
}

func init() {
	n, _ := interest.AnnuityPeriods(2400000, 63000, 0.189218/12)
	fmt.Printf("one pooled debt 24,000 at 18.9218 %%, 630/mo: N=%.3f total≈%.2f\n", n, 630*n)
}
