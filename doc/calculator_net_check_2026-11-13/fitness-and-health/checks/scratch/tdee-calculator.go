//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	r, err := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(140*0.45359237, 65*2.54, 40, true)
	b, _ := anthropometry.BodyMassIndexImperial(140, 65)
	fmt.Printf("REE %.3f ×1.465 (their factor, supplied) = %.3f; BMI %.3f %v\n", r, r*1.465, b, err)
}
