//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/trig"
)

func show(name string, t trig.Triangle, err error) {
	d := 180 / math.Pi
	fmt.Printf("%s a %.12g b %.12g c %.12g A %.12g B %.12g C %.12g %v\n", name, t.SideA, t.SideB, t.SideC, t.AngleARad*d, t.AngleBRad*d, t.AngleCRad*d, err)
}

func main() {
	t, err := trig.SolveRightTriangleHypotenuseLeg(5, 3)
	show("c5 a3", t, err)
	t, err = trig.SolveRightTriangleLegAngle(7, 25*math.Pi/180)
	show("a7 A25", t, err)
	t, err = trig.SolveRightTriangleHypotenuseAngle(10, 0.3)
	show("c10 A0.3rad", t, err)
}
