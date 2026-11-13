//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
	"github.com/CrimsonCoder42/calcbench/num/trig"
)

func deg(r float64) float64 { return r * 180 / math.Pi }

func show(name string, t trig.Triangle, err error) {
	area, _ := geometry.TriangleAreaHeron(t.SideA, t.SideB, t.SideC)
	fmt.Printf("%s a %.12g b %.12g c %.12g A %.12g B %.12g C %.12g area %.12g %v\n", name, t.SideA, t.SideB, t.SideC, deg(t.AngleARad), deg(t.AngleBRad), deg(t.AngleCRad), area, err)
}

func main() {
	t, err := trig.SolveTriangleSAS(1, 1, 60*math.Pi/180)
	show("SAS 1,1,60", t, err)
	t, err = trig.SolveTriangleSSS(3, 4, 6)
	show("SSS 3,4,6", t, err)
	t, err = trig.SolveTriangleASA(40*math.Pi/180, 75*math.Pi/180, 10)
	show("ASA A40 B75 c10", t, err)
	t, err = trig.SolveTriangleSSS(1000000, 1000000, 1e-3)
	show("needle", t, err)
}
