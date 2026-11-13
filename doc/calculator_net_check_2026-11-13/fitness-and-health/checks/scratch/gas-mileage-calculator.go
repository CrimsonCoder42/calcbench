//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/automotive"
)

func main() {
	e, err := automotive.FuelEconomy(12360-12000, 15)
	c, _ := automotive.FuelConsumptionPer100(15, 360)
	fmt.Println("US: mpg", e, "gal/100 mi", c, err)
	e, _ = automotive.FuelEconomy(18900-18300, 50)
	c, _ = automotive.FuelConsumptionPer100(50, 600)
	fmt.Println("metric: km/L", e, "L/100 km", c)
}
