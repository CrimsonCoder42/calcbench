//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/combinatorics"
)

func main() {
	for _, c := range [][2]int{{6, 2}, {52, 5}, {66, 33}, {20, 20}, {70, 35}} {
		p, perr := combinatorics.Permutations(c[0], c[1])
		b, berr := combinatorics.Binomial(c[0], c[1])
		fmt.Println(c, "P", p, perr, "C", b, berr)
	}
}
