//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/invest"
)

func main() {
	f := []int64{-10000000, 500000, 2500000, 3500000, 4000000, 3000000, 1000000}
	p, _ := invest.PaybackPeriod(f)
	d, _ := invest.DiscountedPaybackPeriod(0.10, f)
	irr, _ := interest.InternalRateOfReturn(f)
	fmt.Printf("form2 payback %.4f discounted %.4f IRR %.4f%%\n", p, d, irr*100)
	g := []int64{-10000000}
	for t := 0; t < 5; t++ {
		g = append(g, int64(math.Round(3000000*math.Pow(1.05, float64(t)))))
	}
	p, _ = invest.PaybackPeriod(g)
	d, _ = invest.DiscountedPaybackPeriod(0.10, g)
	irr, _ = interest.InternalRateOfReturn(g)
	fmt.Printf("form1 %v payback %.4f discounted %.4f IRR %.4f%%\n", g, p, d, irr*100)
}
