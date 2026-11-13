//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
	"github.com/CrimsonCoder42/calcbench/unit/area"
)

func main() {
	a, _ := geometry.RectangleArea(12.5, 8)
	f, _ := area.SquareMetersToSquareFeet(a)
	fmt.Printf("%.9f\n", f)
}
