//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/integer"
)

func main() {
	for _, n := range []uint64{120, 97, 1, 963761198400} {
		d, err := integer.Divisors(n)
		p, _ := integer.PrimeFactors(n)
		fmt.Println(n, len(d), d, p, err)
	}
}
