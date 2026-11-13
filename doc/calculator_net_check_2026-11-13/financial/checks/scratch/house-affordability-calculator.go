//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	// Their price P fixes the percent-of-price costs; ours then returns the loan.
	// case 1: 120k/yr, 28 %, tax 1.5 % + ins 0.5 % of price, 7.263 %, 30 y, their price 392,644
	ti1 := int64(math.Round(0.02 * 39264400 / 12))
	l1, err := loan.AffordableLoanAmount(1000000, 0.28, ti1, 0.07263/12, 360)
	fmt.Println("case1 T&I", ti1, "loan", l1, err)
	// case 2: 80k/yr, 28 %, tax 1 % + ins 0.5 % of price, PMI 0.5 %/yr of loan, 5 %, 15 y, their price 213,525, loan 192,173
	ti2 := int64(math.Round(0.015*21352500/12 + 0.005*19217300/12))
	l2, err := loan.AffordableLoanAmount(666667, 0.28, ti2, 0.05/12, 180)
	fmt.Println("case2 T&I+PMI", ti2, "loan", l2, err)
}
