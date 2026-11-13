//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	fmt.Println(geometry.TriangleAreaHeron(3, 4, 6))
	fmt.Println(geometry.TrapezoidArea(5, 9.5, 4))
	fmt.Println(geometry.SectorArea(6, 75*math.Pi/180))
	fmt.Println(geometry.EllipseArea(7, 2.5))
}
