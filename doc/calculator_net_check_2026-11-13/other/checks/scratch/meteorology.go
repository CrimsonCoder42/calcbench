//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/meteorology"
	"github.com/CrimsonCoder42/calcbench/unit/temperature"
)

func main() {
	for _, c := range [][2]float64{{10, 15}, {-20, 40}, {50, 3}} {
		w, err := meteorology.WindChillF(c[0], c[1])
		fmt.Printf("wc %v °F %v mph → %.4f %v\n", c[0], c[1], w, err)
	}
	for _, c := range [][2]float64{{85, 70}, {70, 50}, {100, 10}, {84, 90}, {105, 40}} {
		h, err := meteorology.HeatIndexF(c[0], c[1])
		fmt.Printf("hi %v °F %v %% → %.4f %v\n", c[0], c[1], h, err)
	}
	for _, c := range [][2]float64{{70, 65}, {90, 30}, {32, 80}} {
		tc, _ := temperature.FahrenheitToCelsius(c[0])
		d, err := meteorology.DewPointC(tc, c[1])
		df, _ := temperature.CelsiusToFahrenheit(d)
		fmt.Printf("dp %v °F %v %% → %.4f °C %.4f °F %v\n", c[0], c[1], d, df, err)
	}
	t, _ := temperature.FahrenheitToCelsius(85)
	td, _ := temperature.FahrenheitToCelsius(70)
	rh, _ := meteorology.RelativeHumidityFromDewPoint(t, td)
	hi, _ := meteorology.HeatIndexF(85, rh)
	fmt.Printf("rh %.4f hi(85, Td 70) %.4f\n", rh, hi)
}
