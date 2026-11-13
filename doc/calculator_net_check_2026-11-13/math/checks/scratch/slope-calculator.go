//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	for _, c := range [][4]float64{{1, 1, 2, 2}, {-3, 4, 5, -2}, {2.5, 7, 2.5, 9}} {
		p, q := geometry.Point{X: c[0], Y: c[1]}, geometry.Point{X: c[2], Y: c[3]}
		m, err := geometry.Slope(p, q)
		a, _ := geometry.SlopeAngleRad(c[3]-c[1], c[2]-c[0])
		fmt.Printf("%v slope %.15g %v angle %.12g deg dist %.12g\n", c, m, err, a*180/math.Pi, math.Hypot(c[2]-c[0], c[3]-c[1]))
	}
}
