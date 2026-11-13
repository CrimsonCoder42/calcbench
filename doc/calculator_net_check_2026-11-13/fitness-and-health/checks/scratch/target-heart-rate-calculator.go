//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/cardio"
)

func main() {
	t, err := cardio.MaxHeartRateTanaka(30)
	fmt.Printf("Tanaka 30: %.3f %v\n", t, err)
	n, err := cardio.MaxHeartRateNes(50)
	fmt.Printf("Nes 50: %.3f %v\n", n, err)
}
