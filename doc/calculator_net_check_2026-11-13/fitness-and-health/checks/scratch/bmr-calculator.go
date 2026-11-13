//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/nutrition"
)

func main() {
	for _, x := range []struct {
		lb, in, age float64
		f           bool
	}{{160, 70, 25, false}, {130, 64, 50, true}} {
		kg, cm := x.lb*0.45359237, x.in*2.54
		s, err := nutrition.RestingEnergyExpenditureMifflinStJeorSimplified(kg, cm, x.age, x.f)
		fmt.Printf("%+v simplified=%.3f kJ(×4.184)=%.1f kJ(×4.1868)=%.1f %v\n", x, s, s*4.184, s*4.1868, err)
	}
}
