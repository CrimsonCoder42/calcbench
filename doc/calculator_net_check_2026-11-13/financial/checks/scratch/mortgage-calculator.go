//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	cases := []struct {
		name               string
		loanCents          int64
		ratePerYear        float64
		years              int
		taxCents, insCents int64
	}{
		{"A 300k-20% 6% 30y", 24000000, 0.06, 30, 0, 0},
		{"B 400k-20% 7.263% 30y tax1.2% ins1500", 32000000, 0.07263, 30, 480000, 150000},
		{"C 150k-10% 3.5% 15y", 13500000, 0.035, 15, 0, 0},
		{"D 200k-0% 0% 30y", 20000000, 0, 30, 0, 0},
	}
	for _, c := range cases {
		n := c.years * 12
		pmt, err := interest.AnnuityPaymentFromPresentValue(c.loanCents, c.ratePerYear/12, n)
		if err != nil {
			fmt.Println(c.name, err)
			continue
		}
		rows, _ := loan.AmortizationSchedule(c.loanCents, c.ratePerYear/12, n)
		var tot, ti int64
		for _, r := range rows {
			tot += r.PaymentCents
			ti += r.InterestCents
		}
		piti, _ := loan.PITIPayment(pmt, c.taxCents, c.insCents, 0)
		fmt.Printf("%s: pmt=%d total=%d interest=%d last=%d piti=%d\n", c.name, pmt, tot, ti, rows[len(rows)-1].PaymentCents, piti)
	}
}
