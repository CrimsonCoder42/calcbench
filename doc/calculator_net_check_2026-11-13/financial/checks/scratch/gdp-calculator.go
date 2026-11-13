//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/economics"
)

func main() {
	a, _ := economics.GDPExpenditure(140, 36, 39, 21, 27) // tenths of a trillion
	b, _ := economics.GDPExpenditure(125000055, -3000000, 40000000, 15000000, 22000025)
	fmt.Println(a, b)
}
