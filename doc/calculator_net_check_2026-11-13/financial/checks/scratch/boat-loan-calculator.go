//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	// 35,000 − 7,000 down + 1,400 tax (4 %) + 2,200 fees financed = 31,600
	p, _ := interest.AnnuityPaymentFromPresentValue(3160000, 0.07/12, 120)
	fmt.Println(p)
}
