//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/astronomy"
)

func main() {
	for _, n := range []int64{27689, -1000, 36525} {
		y, m, d, err := astronomy.GregorianAddDays(2026, 11, 13, n)
		w, _ := astronomy.GregorianDayOfWeek(y, m, d)
		fmt.Println(n, y, m, d, w, err)
	}
}
