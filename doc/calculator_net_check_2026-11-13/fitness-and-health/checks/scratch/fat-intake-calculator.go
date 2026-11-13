//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	r, _ := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(160*0.45359237, 70*2.54, 25, false)
	e := r * 1.2
	g, err := nutrition.AMDRGrams(e, 25)
	fmt.Printf("energy %.3f fat %.2f-%.2f %v\n", e, g.FatLowG, g.FatHighG, err)
}
