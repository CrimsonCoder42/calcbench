//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/geometry"
)

func main() {
	for _, r := range []float64{5, 0.001, 123456.789} {
		c, _ := geometry.CircleCircumference(r)
		a, _ := geometry.CircleArea(r)
		fmt.Printf("r %v d %v C %.15g A %.15g\n", r, 2*r, c, a)
	}
}
