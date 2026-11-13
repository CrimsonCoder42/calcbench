//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, c := range []struct {
		pv, c1 int64
		r, g   float64
		n      int
	}{{2000000, 500000, 0.03, 0.03, 10}, {1000000, 200000, 0.07, 0.05, 20}} {
		a, _ := interest.FutureValue(c.pv, c.r, c.n)
		b, err := interest.FutureValueGrowingAnnuity(c.c1, c.r, c.g, c.n)
		fmt.Printf("FV(init)=%d FVGA(end)=%d total=%d err=%v\n", a, b, a+b, err)
	}
}
