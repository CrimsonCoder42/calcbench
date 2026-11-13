//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/business"
)

func main() {
	a, _ := business.TieredMarginalAmount(20000000, []int64{0}, []float64{0.03})
	b, _ := business.TieredMarginalAmount(20000000, []int64{0, 2000000}, []float64{0.03, 0.05})
	fmt.Println(a, b)
}
