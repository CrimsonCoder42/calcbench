//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	for _, c := range [][2]float64{{160, 70}, {250, 64}, {100, 72}} {
		b, err := anthropometry.BodyMassIndexImperial(c[0], c[1])
		cat, _ := anthropometry.BMICategoryCDCAdult(b)
		fmt.Printf("imperial %v lb %v in -> %.4f %s %v\n", c[0], c[1], b, cat, err)
	}
	for _, c := range [][2]float64{{65, 1.80}, {120, 1.70}} {
		b, err := anthropometry.BodyMassIndex(c[0], c[1])
		cat, _ := anthropometry.BMICategoryCDCAdult(b)
		fmt.Printf("metric %v kg %v m -> %.4f %s %v\n", c[0], c[1], b, cat, err)
	}
}
