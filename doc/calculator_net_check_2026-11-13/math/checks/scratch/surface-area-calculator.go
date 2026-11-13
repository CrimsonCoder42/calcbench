//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	fmt.Println(geometry.SphereSurfaceArea(2))
	fmt.Println(geometry.ConeSurfaceArea(3, 7.5))
	fmt.Println(geometry.CylinderSurfaceArea(1.5, 4))
	fmt.Println(geometry.EllipsoidSurfaceArea(1, 2, 3.5))
}
