//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	v, err := geometry.PrismVolume(200, 2.0/12) // ft² × ft
	fmt.Printf("200 ft² × 2 in: %.4f ft³ = %.4f yd³ (÷27) %v\n", v, v/27, err)
}
