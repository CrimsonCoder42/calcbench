//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/business"
)

func main() {
	s, _ := business.DepreciationSumOfYearsDigits(1100000, 100000, 5)
	d, _ := business.DepreciationDecliningBalance(1100000, 100000, 5, 2)
	fmt.Println("SYD", s, "DDB", d)
}
