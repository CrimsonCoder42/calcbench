//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
	"github.com/CrimsonCoder42/calcbench/fin/loan"
)

func main() {
	end, _ := loan.BalloonLoanPayment(2000000, 800000, 0.005, 36)
	// payments in advance: PV − FV·v^n spread as an annuity due
	pvRes, _ := interest.PresentValue(800000, 0.005, 36)
	due, _ := interest.AnnuityDuePaymentFromPresentValue(2000000-pvRes, 0.005, 36)
	fmt.Println("end", end, "due", due)
	c2, _ := loan.BalloonLoanPayment(4500000, 2500000, 0.045/12, 51)
	fmt.Println("case2 end", c2)
}
