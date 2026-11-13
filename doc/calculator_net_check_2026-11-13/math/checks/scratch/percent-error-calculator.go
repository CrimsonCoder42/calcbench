//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func main() {
	fmt.Println(algebra.PercentError(3.1, 3.0))
	fmt.Println(algebra.PercentError(10, 11))
	fmt.Println(algebra.PercentError(-11, -10))
	fmt.Println(algebra.PercentError(5, 0))
}
