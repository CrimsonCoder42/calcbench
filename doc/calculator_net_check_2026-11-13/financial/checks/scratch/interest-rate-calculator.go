//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		pv, pmt int64
		n       int
	}{{3200000, 96000, 36}, {25000000, 160000, 360}} {
		r, err := interest.AnnuityRate(c.pv, c.pmt, 0, c.n)
		fmt.Printf("pv=%d pmt=%d n=%d: rate/month=%.8f annual=%.5f%% err=%v\n", c.pv, c.pmt, c.n, r, r*1200, err)
	}
}
