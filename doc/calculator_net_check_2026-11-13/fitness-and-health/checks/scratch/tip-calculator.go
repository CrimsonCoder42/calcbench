//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/shopping"
)

func main() {
	for _, c := range [][2]int64{{5500, 1500}, {1010, 1500}, {4733, 1800}} {
		t, err := shopping.TipCents(c[0], c[1])
		fmt.Printf("bill %d c, %d bp: tip %d c, total %d c %v\n", c[0], c[1], t, c[0]+t, err)
	}
}
