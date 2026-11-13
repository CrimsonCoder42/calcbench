//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/fitness"
)

func main() {
	for _, c := range []struct {
		load float64
		reps int
	}{{100, 5}, {225, 8}, {60, 10}} {
		b, err := fitness.OneRepMaxBrzycki(c.load, c.reps)
		fmt.Printf("%+v Brzycki %.3f %v\n", c, b, err)
	}
}
