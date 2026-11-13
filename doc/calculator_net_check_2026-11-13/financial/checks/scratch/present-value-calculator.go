//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.PresentValue(100000, 0.06, 10)
	b, _ := interest.PresentValueAnnuityDue(10000, 0.06, 10)
	c, _ := interest.PresentValueAnnuity(10000, 0.06, 10)
	fmt.Println(a, b, c)
}
