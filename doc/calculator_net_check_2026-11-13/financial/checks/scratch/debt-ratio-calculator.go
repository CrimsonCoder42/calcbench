//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	// case 1: income 60,000/yr = 500,000 c/mo; housing 1,200; other 200 + 250
	b, _ := loan.DebtToIncomeRatio(120000+20000+25000, 500000)
	f, _ := loan.DebtToIncomeRatio(120000, 500000)
	fmt.Printf("case1 back %.4f front %.4f\n", b, f)
	// case 2: income 95,000/yr + 500/mo = 841,667 c/mo; housing 2,100 + 400 tax + 100 ins; other 400 + 300
	b, _ = loan.DebtToIncomeRatio(210000+40000+10000+40000+30000, 841667)
	f, _ = loan.DebtToIncomeRatio(260000, 841667)
	fmt.Printf("case2 back %.4f front %.4f\n", b, f)
}
