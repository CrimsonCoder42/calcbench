//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/prob"
)

func main() {
	for _, c := range [][2]float64{{0.5, 0.4}, {0.37, 0.81}} {
		a, b := c[0], c[1]
		and, _ := prob.IndependentJointProbability(a, b)
		or, _ := prob.UnionProbability(a, b, and)
		na, _ := prob.ComplementProbability(a)
		nb, _ := prob.ComplementProbability(b)
		neither, _ := prob.IndependentJointProbability(na, nb)
		fmt.Printf("A %v B %v: and %.17g or %.17g xor %.17g neither %.17g notA %.17g\n", a, b, and, or, or-and, neither, na)
	}
	fmt.Println(prob.NormalIntervalProbability(-1, 1, 0, 1))
	fmt.Println(prob.NormalIntervalProbability(-1.96, 1.96, 0, 1))
	fmt.Println(prob.NormalIntervalProbability(110, 130, 100, 15))
}
