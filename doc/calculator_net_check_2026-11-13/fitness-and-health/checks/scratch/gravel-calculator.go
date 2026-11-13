//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	a, _ := geometry.CircleArea(5) // ft²
	v, err := geometry.PrismVolume(a, 0.25)
	m3 := v * 0.028316846592 // ft³ → m³, NIST SP 811 B.8 (exact)
	kg := m3 * 1600
	fmt.Printf("area %.3f ft², volume %.4f ft³ = %.5f m³, mass %.2f kg = %.2f lb %v\n", a, v, m3, kg, kg/0.45359237, err)
}
