//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/linalg"
)

func main() {
	a := []float64{2, 1, 0, 3, 1, 3, 2, 0, 0, 2, 4, 1, 5, 0, 1, 2}
	b := []float64{1, 2, 0, -1, 0, 1, 3, 2, 4, 0, 1, 1, -2, 1, 0, 3}
	d, err := linalg.MatDet(a, 4)
	fmt.Println("det", d, err)
	inv, err := linalg.MatInverse(nil, a, 4)
	fmt.Println("inv", inv, err)
	p, err := linalg.MatMul(nil, a, 4, 4, b, 4)
	fmt.Println("AB", p, err)
}
