//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/astronomy"
)

func main() {
	for _, d := range [][3]int{{2026, 11, 13}, {1954, 6, 30}, {2000, 2, 29}} {
		w, err := astronomy.GregorianDayOfWeek(int64(d[0]), d[1], d[2])
		fmt.Println(d, w, err)
	}
}
