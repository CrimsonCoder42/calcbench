//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	for _, c := range []struct {
		p     int64
		r     float64
		n     int
		extra int64
	}{{20000000, 0.06, 180, 0}, {20000000, 0.06, 180, 20000}, {1000000, 0.12, 5, 0}} {
		rows, err := loan.AmortizationSchedule(c.p, c.r/12, c.n)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var tot, ti int64
		for _, r := range rows {
			tot += r.PaymentCents
			ti += r.InterestCents
		}
		fmt.Printf("p=%d r=%v n=%d pmt=%d total=%d interest=%d row1=%+v row12bal=%d\n", c.p, c.r, c.n, rows[0].PaymentCents, tot, ti, rows[0], rows[min(11, len(rows)-1)].BalanceCents)
		if c.extra > 0 {
			k, saved, err := loan.ExtraPaymentPayoff(c.p, c.r/12, c.n, c.extra)
			fmt.Printf("  extra %d: periods=%d saved=%d interestWithExtra=%d err=%v\n", c.extra, k, saved, ti-saved, err)
		}
	}
}
