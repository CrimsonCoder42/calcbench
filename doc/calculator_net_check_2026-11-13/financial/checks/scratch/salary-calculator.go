//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/personal"
)

func main() {
	w1, _ := personal.WeeklyPayEquivalent(866667, 12)
	w2, _ := personal.WeeklyPayEquivalent(433333, 24)
	h1, _ := personal.HourlyRateFromWeeklyPay(200000, 40)
	fmt.Println("case1 weekly from monthly", w1, "from semimonthly", w2, "hourly", h1)
	w, _ := personal.WeeklyPayEquivalent(650000, 12)
	h, _ := personal.HourlyRateFromWeeklyPay(w, 37.5)
	fmt.Println("case2 weekly", w, "hourly", h, "annual", 650000*12, "biweekly = 2 × weekly", 2*w)
}
