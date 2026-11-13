//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	type c struct {
		lb, in, age float64
		f           bool
	}
	for _, x := range []c{{165, 70, 25, false}, {140, 65, 40, true}, {220, 74, 60, false}} {
		kg := x.lb * 0.45359237
		cm := x.in * 2.54
		s, _ := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(kg, cm, x.age, x.f)
		fmt.Printf("simplified=%.3f x1.2=%.3f\n", s, s*1.2)
		r, err := nutrition.RestingEnergyExpenditureMifflinStJeor(kg, cm, x.age, x.f)
		fmt.Printf("%+v kg=%.4f cm=%.2f REE=%.3f x1.465=%.3f x1.2=%.3f %v\n", x, kg, cm, r, r*1.465, r*1.2, err)
	}
}
