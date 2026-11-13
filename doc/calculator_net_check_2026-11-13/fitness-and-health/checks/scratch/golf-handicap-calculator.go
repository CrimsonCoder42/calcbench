//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/life/sports"
)

// diff is the Rule 5.1 Score Differential as quoted in A-OTHER's note (113 / slope × (score − rating), PCC 0),
// NOT read first hand; used here only to prepare comparison inputs.
func diff(rating, slope, score float64) float64 {
	return math.Round(113/slope*(score-rating)*10) / 10
}

func main() {
	d1, d3, d4, d5, d6 := diff(70.3, 126, 85), diff(66.8, 120, 82), diff(65.9, 118, 85), diff(74.3, 137, 88), diff(71.1, 126, 90)
	fmt.Println("differentials", d1, d3, d4, d5, d6)
	a, err := sports.HandicapIndexWHS([]float64{d1, d3, d4, d5, d6})
	fmt.Println("5 eighteen-hole rounds:", a, err)
	b, err := sports.HandicapIndexWHS([]float64{d1, d3, d5})
	fmt.Println("rounds 1, 3, 5:", b, err)
}
