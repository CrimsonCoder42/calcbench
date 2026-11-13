//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/intmath"
)

func main() {
	fmt.Println(intmath.PowerBySquaring(2, 5))
	fmt.Println(intmath.PowerBySquaring(-3, 39))
	fmt.Println(intmath.PowerBySquaring(7, 22))
	fmt.Println(intmath.PowerBySquaring(10, 19))
}
