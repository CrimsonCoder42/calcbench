//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	kg := 160 * 0.45359237
	lo, hi, err := nutrition.ProteinRangeISSN(kg)
	fmt.Printf("kg %.3f ISSN %.2f-%.2f %v; 0.8 g/kg %.2f\n", kg, lo, hi, err, 0.8*kg)
	r, _ := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(kg, 70*2.54, 25, false)
	g, _ := nutrition.AMDRGrams(r*1.375, 25)
	fmt.Printf("energy %.1f AMDR protein %.2f-%.2f\n", r*1.375, g.ProteinLowG, g.ProteinHighG)
}
