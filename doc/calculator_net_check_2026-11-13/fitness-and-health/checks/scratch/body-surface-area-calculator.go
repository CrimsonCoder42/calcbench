//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	for _, c := range [][2]float64{{155, 69}, {200, 74}} {
		h, err := anthropometry.BodySurfaceAreaHaycock(c[0]*0.45359237, c[1]*2.54)
		fmt.Printf("%v lb %v in: Haycock %.4f %v\n", c[0], c[1], h, err)
	}
}
