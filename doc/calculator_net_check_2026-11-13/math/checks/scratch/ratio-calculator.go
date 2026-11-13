//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func main() {
	fmt.Println(algebra.ProportionFourthTerm(3, 4, 600))
	fmt.Println(algebra.ProportionFourthTerm(7, 3, 10))
	fmt.Println(algebra.ProportionFourthTerm(-2, 5, 4))
}
