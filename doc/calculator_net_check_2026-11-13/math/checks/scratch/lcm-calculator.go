//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/integer"
)

func fold(xs []uint64) (uint64, uint64, error) {
	l, g := xs[0], xs[0]
	var err error
	for _, x := range xs[1:] {
		if l, err = integer.LCM(l, x); err != nil {
			return 0, 0, err
		}
		g = integer.GCD(g, x)
	}
	return l, g, nil
}

func main() {
	for _, xs := range [][]uint64{{330, 75, 450, 225}, {12, 18}, {17, 19, 23, 29, 31, 37, 41, 43}, {40902, 24140}, {1 << 40, 3 << 30}} {
		fmt.Println(xs, "lcm, gcd =")
		fmt.Println(fold(xs))
	}
}
