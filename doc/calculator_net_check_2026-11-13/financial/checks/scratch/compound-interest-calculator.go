//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	ear, _ := interest.EffectiveAnnualRate(0.06, 12)
	fmt.Printf("6%% monthly -> EAR %.5f%%\n", ear*100)
	ear2, _ := interest.EffectiveAnnualRate(0.09, 4)
	p, _ := interest.PeriodicRateFromEffectiveAnnual(ear2, 26)
	fmt.Printf("9%% quarterly -> biweekly nominal %.5f%%, per period %.5f%%\n", p*26*100, p*100)
}
