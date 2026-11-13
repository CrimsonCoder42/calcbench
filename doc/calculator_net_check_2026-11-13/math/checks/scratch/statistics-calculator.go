//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	for _, s := range [][]float64{
		{10, 2, 38, 23, 38, 23, 21, 23},
		{10, 2, 38, 23, 38, 23, 21},
		{1.5, -2.25, 7, 7, 3.125, -2.25, 0},
	} {
		m, _ := stats.ArithmeticMean(s)
		med, _ := stats.Median(s)
		modes, cnt, _ := stats.Mode(s)
		r, _ := stats.Range(s)
		ss, _ := stats.SampleStandardDeviation(s)
		ps, _ := stats.PopulationStandardDeviation(s)
		g, gerr := stats.GeometricMean(s)
		fmt.Println("  geometric mean", g, gerr)
		fmt.Printf("n %d mean %.15g median %.15g modes %v x%d range %.15g s %.15g sigma %.15g\n", len(s), m, med, modes, cnt, r, ss, ps)
		for _, meth := range []stats.QuantileMethod{stats.HyndmanFan6, stats.HyndmanFan7, stats.HyndmanFan2} {
			q, _ := stats.Percentiles(s, []float64{0.25, 0.75}, meth)
			fmt.Printf("  HF%d Q1 %.15g Q3 %.15g\n", meth, q[0], q[1])
		}
	}
}
