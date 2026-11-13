//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/prob"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	fmt.Println(stats.ZScore(5, 3, 2))
	fmt.Println(stats.ZScore(-12.5, 4, 7.25))
	for _, z := range []float64{2, -3.5, 0.674} {
		lo, _ := prob.NormalCDF(z, 0, 1)
		hi, _ := prob.NormalSurvival(z, 0, 1)
		fmt.Printf("z %v P(<z) %.10g P(>z) %.10g\n", z, lo, hi)
	}
	fmt.Println(prob.NormalIntervalProbability(-1, 0, 0, 1))
}
