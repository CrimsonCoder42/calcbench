//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	sets := [][]float64{
		{10, 12, 23, 23, 16, 23, 21, 16},
		{10000001, 10000003, 10000002},
		{2, 4, 4, 4, 5, 5, 7, 9},
		{1000000000.2, 1000000000.1, 1000000000.3, 1000000000.1, 1000000000.3},
	}
	for _, s := range sets {
		p, _ := stats.PopulationStandardDeviation(s)
		ss, _ := stats.SampleStandardDeviation(s)
		pv, _ := stats.PopulationVariance(s)
		sv, _ := stats.SampleVariance(s)
		m, _ := stats.ArithmeticMean(s)
		fmt.Printf("mean %.15g popSD %.15g sampSD %.15g popVar %.15g sampVar %.15g\n", m, p, ss, pv, sv)
	}
}
