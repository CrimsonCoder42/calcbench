//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/unit/length"
)

func main() {
	a, _ := length.InchesToCentimeters(12)
	b, _ := length.CentimetersToInches(100)
	fmt.Printf("%.10f %.10f\n", a, b)
}
