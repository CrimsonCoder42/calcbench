//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func main() {
	fmt.Println(algebra.QuadraticRealRoots(1, 1, 0.25))
	fmt.Println(algebra.QuadraticRealRoots(2, -7, 3))
	fmt.Println(algebra.QuadraticRoots(1, 2, 5))
	fmt.Println(algebra.QuadraticRealRoots(1, -1e8, 1))
	fmt.Println(algebra.QuadraticRealRoots(1, 0, -2))
}
