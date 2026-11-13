//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	const lbPerKg = 1 / 0.45359237
	for _, in := range []float64{70, 62, 76} {
		m := in * 0.0254
		lo, _ := anthropometry.WeightAtBMI(18.5, m)
		hi, _ := anthropometry.WeightAtBMI(25, m)
		loI, _ := anthropometry.WeightAtBMIImperial(18.5, in)
		hiI, _ := anthropometry.WeightAtBMIImperial(25, in)
		fmt.Printf("%v in: metric %.3f-%.3f lb (%.3f-%.3f kg); imperial x703 %.3f-%.3f lb\n", in, lo*lbPerKg, hi*lbPerKg, lo, hi, loI, hiI)
	}
}
