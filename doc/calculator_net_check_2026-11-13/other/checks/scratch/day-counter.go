//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/astronomy"
)

func main() {
	for _, c := range [][6]int{{1910, 4, 20, 1986, 2, 9}, {2000, 2, 29, 2025, 3, 1}, {2000, 2, 29, 2028, 2, 29}, {1990, 5, 15, 2026, 11, 13}, {2024, 1, 1, 2024, 12, 31}, {1899, 12, 31, 1900, 3, 1}} {
		d, err := astronomy.DaysBetweenDates(int64(c[0]), c[1], c[2], int64(c[3]), c[4], c[5])
		fmt.Println(c, d, err)
	}
}
