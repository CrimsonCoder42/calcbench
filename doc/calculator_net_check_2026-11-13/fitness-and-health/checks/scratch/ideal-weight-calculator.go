//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	const lbPerKg = 1 / 0.45359237
	for _, c := range []struct {
		h float64
		m bool
	}{{70, true}, {64, false}} {
		d, _ := anthropometry.IdealBodyWeightDevine(c.h, c.m)
		r, _ := anthropometry.IdealBodyWeightRobinson(c.h, c.m)
		mi, _ := anthropometry.IdealBodyWeightMiller(c.h, c.m)
		h, _ := anthropometry.IdealBodyWeightHamwiLb(c.h, c.m)
		fmt.Printf("%v in male=%v: Robinson %.2f lb, Miller %.2f lb, Devine %.2f lb, Hamwi %.2f lb\n", c.h, c.m, r*lbPerKg, mi*lbPerKg, d*lbPerKg, h)
	}
}
