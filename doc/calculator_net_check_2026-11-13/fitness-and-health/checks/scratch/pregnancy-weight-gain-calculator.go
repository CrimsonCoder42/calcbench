//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
	"github.com/CrimsonCoder42/calcbench/med/obstetric"
)

func main() {
	b1, _ := anthropometry.BodyMassIndexImperial(120, 66)
	r, err := obstetric.GestationalWeightGainIOM(b1)
	fmt.Printf("BMI %.2f %+v %v -> term %.1f-%.1f lb\n", b1, r, err, 120+r.TotalLowLb, 120+r.TotalHighLb)
	b2, _ := anthropometry.BodyMassIndexImperial(190, 64)
	c, lo, hi, err := obstetric.GestationalWeightGainTwinsIOM(b2)
	fmt.Printf("BMI %.2f twins %s %v-%v %v -> term %.1f-%.1f lb\n", b2, c, lo, hi, err, 190+lo, 190+hi)
}
