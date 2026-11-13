//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

// Feet throughout; 5 in = 5/12 ft (1 in = 1/12 ft exactly, NIST SP 811 B.8); 1 yd³ = 27 ft³.
func main() {
	base, _ := geometry.RectangleArea(5, 2.5)
	slab, _ := geometry.PrismVolume(base, 5.0/12)
	hole, _ := geometry.CylinderVolume(1.25, 6)
	fmt.Printf("slab %.4f ft3 %.4f yd3\n", slab, slab/27)
	fmt.Printf("hole %.4f ft3 %.4f yd3\n", hole, hole/27)
}
