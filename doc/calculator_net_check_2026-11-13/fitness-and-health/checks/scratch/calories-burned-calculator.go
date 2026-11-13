//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/fitness"
)

func main() {
	for _, c := range []struct{ met, kg, h float64 }{{2.0, 150 * 0.45359237, 0.75}, {2.0, 100, 1}, {6.0, 70, 0.5}} {
		e, err := fitness.EnergyExpenditureMET(c.met, c.kg, c.h)
		fmt.Printf("%+v -> %.3f kcal (×1.05 for 3.5 mL × 5 kcal/L: %.3f) %v\n", c, e, e*1.05, err)
	}
}
