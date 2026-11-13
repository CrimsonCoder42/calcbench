//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

// Roof area = plan area × (sloped length per 12 in of run)/12; sloped length = hypot(12, rise).
func main() {
	for _, rise := range []float64{6, 12, 1} {
		h, _ := geometry.RightTriangleHypotenuse(12, rise)
		fmt.Printf("pitch %v/12: %.2f ft2\n", rise, 2000*h/12)
	}
}
