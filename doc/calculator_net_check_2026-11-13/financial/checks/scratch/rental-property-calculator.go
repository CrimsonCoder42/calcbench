//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/invest"
)

func main() {
	c1, _ := invest.CapitalizationRate(1610000, 20000000)
	c2, _ := invest.CapitalizationRate(1951440, 35000000)
	fmt.Printf("cap1 %.6f cap2 %.6f\n", c1, c2)
	// case 2 cash flows (no loan): −price, NOI growing 3 %/yr for 20 years, sale at price·1.03^20 less 8 %.
	flows := make([]int64, 21)
	flows[0] = -35000000
	for t := 1; t <= 20; t++ {
		flows[t] = int64(math.Round(1951440 * math.Pow(1.03, float64(t-1))))
	}
	flows[20] += int64(math.Round(35000000 * math.Pow(1.03, 20) * 0.92))
	irr, err := interest.InternalRateOfReturn(flows)
	fmt.Printf("irr2 %.6f err=%v\n", irr, err)
}
