//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.FutureValue(2000000, 0.06, 10)
	b, _ := interest.FutureValueAnnuity(200000, 0.06, 10)
	fmt.Println("case1 FV (their sign: −(a − b)) =", -(a - b))
	c, _ := interest.FutureValue(1500000, 0.045/12, 60)
	d, _ := interest.FutureValueAnnuityDue(25000, 0.045/12, 60)
	fmt.Println("case2 FV =", -(c - d))
}
