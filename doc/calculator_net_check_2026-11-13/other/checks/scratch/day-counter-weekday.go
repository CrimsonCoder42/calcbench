//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/sci/astronomy"
)

func main() {
	for _, d := range [][3]int{{1910, 4, 20}, {1986, 2, 9}, {1899, 12, 31}, {1900, 3, 1}, {2024, 1, 1}, {2024, 12, 31}} {
		w, _ := astronomy.GregorianDayOfWeek(int64(d[0]), d[1], d[2])
		fmt.Println(d, w)
	}
}
