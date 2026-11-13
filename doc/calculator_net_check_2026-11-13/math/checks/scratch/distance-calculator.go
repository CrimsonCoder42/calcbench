//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
	"github.com/CrimsonCoder42/calcbench/num/linalg"
	"github.com/CrimsonCoder42/calcbench/num/trig"
)

func main() {
	fmt.Println(geometry.PointDistance(geometry.Point{X: -2.5, Y: 7}, geometry.Point{X: 10, Y: -3.25}))
	fmt.Println(linalg.Vec3Distance(linalg.Vec3{1, -2, 3.5}, linalg.Vec3{-4, 6, 0}))
	r := math.Pi / 180
	a, err := trig.CentralAngle(38.8976*r, -77.0366*r, 39.9496*r, -75.1503*r)
	// Reference only: sphere of mean radius 6371.0088 km (Moritz, Geodetic Reference System 1980, J. Geodesy 2000).
	fmt.Println(a, err, "sphere km", a*6371.0088)
}
