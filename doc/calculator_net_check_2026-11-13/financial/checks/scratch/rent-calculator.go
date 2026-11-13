//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	income := int64(666667) // 80,000/12 in cents
	for _, rent := range []int64{190000, 136667} {
		front, _ := loan.DebtToIncomeRatio(rent, income)
		back, err := loan.DebtToIncomeRatio(rent+50000, income)
		fmt.Printf("rent %d: housing ratio %.4f, back-end DTI %.4f %v\n", rent, front, back, err)
	}
}
