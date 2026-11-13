//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/physics"
	"github.com/CrimsonCoder42/calcbench/unit/mass"
)

// 1 ft³ = 0.028316846592 m³ exactly (NIST SP 811 B.8, (0.3048 m)³); no unit/volume yet, so the factor is here.
func main() {
	d, _ := physics.Density(8900, 1)
	fmt.Printf("%.9f\n", d)
	kg, _ := mass.PoundsToKilograms(100)
	d2, _ := physics.Density(kg, 0.028316846592)
	fmt.Printf("%.9f\n", d2)
	m, _ := physics.MassFromDensity(1000, 0.25)
	fmt.Printf("%.9f\n", m)
}
