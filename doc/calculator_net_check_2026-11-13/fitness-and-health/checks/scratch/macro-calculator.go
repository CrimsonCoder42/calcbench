//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	kg, cm := 165*0.45359237, 70*2.54
	r, _ := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(kg, cm, 25, false)
	g, err := nutrition.AMDRGrams(r, 25)
	fmt.Printf("REE %.3f %+v %v\n", r, g, err)
	g2, _ := nutrition.AMDRGrams(1740, 25)
	fmt.Printf("at 1740: %+v\n", g2)
}
