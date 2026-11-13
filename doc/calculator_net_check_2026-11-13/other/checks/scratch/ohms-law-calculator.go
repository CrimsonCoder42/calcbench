//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/eng/electrical"
)

func main() {
	a, _ := electrical.OhmsLawFromVoltageCurrent(110, 0.2)
	b, _ := electrical.OhmsLawFromResistancePower(47, 5)
	fmt.Printf("%+v\n%.14f %+v\n", a, b.Amperes, b)
	d, _ := electrical.VoltageDropTwoWire(10, 1.2/1000, 500*0.3048)
	fmt.Printf("vd %.5f\n", d)
}
