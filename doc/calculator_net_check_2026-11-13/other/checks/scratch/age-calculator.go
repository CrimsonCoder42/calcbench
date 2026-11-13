//go:build ignore

package main

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/sci/astronomy"
)

func main() {
	j, _ := astronomy.GregorianToJulianDayNumber(-4712, 1, 1)
	k, _ := astronomy.GregorianToJulianDayNumber(1000000, 12, 31)
	fmt.Println("min max", j, k)
	for _, d := range [][3]int{{2000, 1, 1}, {1957, 10, 4}, {1954, 6, 30}} {
		n, _ := astronomy.GregorianToJulianDayNumber(int64(d[0]), d[1], d[2])
		w, _ := astronomy.GregorianDayOfWeek(int64(d[0]), d[1], d[2])
		fmt.Println(d, n, w)
	}
	a, _ := astronomy.DaysBetweenDates(1910, 4, 20, 1986, 2, 9)
	fmt.Println("halley", a)
	for _, d := range [][3]int{{1978, 11, 14}, {1988, 4, 22}} {
		n, _ := astronomy.GregorianDayOfYear(int64(d[0]), d[1], d[2])
		fmt.Println(d, n)
	}
	// round trip vs Go time, every 97 days over the whole range
	bad := 0
	base, _ := astronomy.GregorianToJulianDayNumber(1970, 1, 1)
	for z := int64(38); z <= 366963925; z += 97 {
		y, m, d, err := astronomy.JulianDayNumberToGregorian(z)
		t := time.Unix((z-base)*86400, 0).UTC()
		if err != nil || int(y) != t.Year() || m != int(t.Month()) || d != t.Day() {
			bad++
		}
		if back, _ := astronomy.GregorianToJulianDayNumber(y, m, d); back != z {
			bad++
		}
	}
	fmt.Println("roundtrip bad", bad)
	// compare cases
	fmt.Println(astronomy.DaysBetweenDates(1990, 5, 15, 2026, 11, 13))
	fmt.Println(astronomy.DaysBetweenDates(2000, 2, 29, 2025, 3, 1))
	fmt.Println(astronomy.DaysBetweenDates(1900, 1, 1, 2026, 11, 13))
}
