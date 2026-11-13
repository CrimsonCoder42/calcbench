//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
	"github.com/CrimsonCoder42/calcbench/num/trig"
)

// Total height 120 in, 12 treads of 10 in, 13 risers. Rise = 120/13. Stringer (their convention, read from
// their output): legs = total run (12 × 10 in) and total height − one rise.
func main() {
	for _, treads := range []float64{12, 15} {
		rise := 120 / (treads + 1)
		run := treads * 10
		s, _ := geometry.RightTriangleHypotenuse(run, 120-rise)
		tr, _ := trig.SolveRightTriangleLegs(10, rise)
		fmt.Printf("treads %v rise %.2f in, stringer %.2f in (%.0f ft %.2f in), stringer height %.2f in\n", treads, rise, s, math.Floor(s/12), s-12*math.Floor(s/12), 120-rise)
		fmt.Printf("angles A %.2f B %.2f deg\n", tr.AngleARad*180/math.Pi, tr.AngleBRad*180/math.Pi)
	}
}
