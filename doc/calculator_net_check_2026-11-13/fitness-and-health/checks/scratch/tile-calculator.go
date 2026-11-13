//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	floor, _ := geometry.RectangleArea(20*12, 15*12) // in²
	cell, _ := geometry.RectangleArea(12.25, 12.25)  // tile + 0.25 in gap, in²
	n := floor / cell
	fmt.Printf("floor %.0f in², cell %.4f in², tiles %.3f → %v; +5%% %.1f → %v; +10%% %.1f → %v\n", floor, cell, n, math.Ceil(n), math.Ceil(n)*1.05, math.Ceil(math.Ceil(n)*1.05), math.Ceil(n)*1.10, math.Ceil(math.Ceil(n)*1.10))
}
