//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func main() {
	m, err := stats.WeightedMean([]float64{90, 80, 88}, []float64{5, 20, 20})
	fmt.Printf("%.10f %v\n", m, err)
	g, err := stats.WeightedMean([]float64{4, 3, 3.3}, []float64{3, 4, 2})
	fmt.Printf("gpa %.10f %v\n", g, err)
}
