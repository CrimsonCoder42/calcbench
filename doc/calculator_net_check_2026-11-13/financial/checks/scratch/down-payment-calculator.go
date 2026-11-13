//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	i := 0.07263 / 12
	// form 1: cash 100,000 covers 20 % down + 3 % closing → price = cash / 0.23
	cash := 10000000.0
	price := int64(cash/0.23 + 0.5)
	l1 := price - price/5
	p1, _ := interest.AnnuityPaymentFromPresentValue(l1, i, 360)
	fmt.Printf("form1 price=%d loan=%d pmt=%d\n", price, l1, p1)
	p2, _ := interest.AnnuityPaymentFromPresentValue(40000000, i, 360)
	fmt.Printf("form2 loan=400000 pmt=%d\n", p2)
	ltv, _ := loan.LoanToValueRatio(41500000, 50000000)
	p3, _ := interest.AnnuityPaymentFromPresentValue(41500000, i, 360)
	fmt.Printf("form3 LTV=%.4f down=%.4f pmt=%d\n", ltv, 1-ltv, p3)
}
