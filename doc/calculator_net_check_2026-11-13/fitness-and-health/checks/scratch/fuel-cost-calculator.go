//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/automotive"
)

func main() {
	c, err := automotive.FuelCostCents(500, 25, 300)
	fmt.Println("500 mi, 25 mpg, 300 c/gal:", c, err)
	c, err = automotive.FuelCostCents(380, 100/7.5, 189)
	fmt.Printf("380 km, 7.5 L/100 km (= %.6f km/L), 189 c/L: %d %v (exact 5386.5)\n", 100/7.5, c, err)
	c, err = automotive.FuelCostCents(1234, 31, 359)
	fmt.Println("1234 mi, 31 mpg, 359 c/gal:", c, err, float64(1234)/31*359)
}
