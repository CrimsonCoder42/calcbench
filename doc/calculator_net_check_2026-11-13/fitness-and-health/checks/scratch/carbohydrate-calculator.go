//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	r, _ := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(160*0.45359237, 70*2.54, 25, false)
	e := r * 1.2 // their activity factor 1.2, supplied (not ours: NEEDS SOURCE)
	g, err := nutrition.AMDRGrams(e, 25)
	fmt.Printf("REE %.3f energy %.3f carb %.2f-%.2f fat %.2f-%.2f %v\n", r, e, g.CarbohydrateLowG, g.CarbohydrateHighG, g.FatLowG, g.FatHighG, err)
}
