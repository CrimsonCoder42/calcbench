//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.AnnuityPaymentFromPresentValue(4000000, 0.05/12, 60)
	// case 2: 32,000 − 1,000 incentive − 3,000 down − 5,000 trade + 2,000 owed + 1,890 tax (7 % of 27,000) + 900 fees
	b, _ := interest.AnnuityPaymentFromPresentValue(2779000, 0.069/12, 72)
	fmt.Println(a, b)
}
