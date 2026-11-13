//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(100000, 0.06, 10)
	b, _ := interest.FutureValueAnnuity(10000, 0.06, 10)
	c, _ := interest.FutureValueAnnuityDue(25000, 0.005, 40)
	fmt.Println(a, b, a+b, c)
}
