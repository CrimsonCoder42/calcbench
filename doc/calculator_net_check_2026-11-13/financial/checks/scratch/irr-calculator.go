//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	for _, f := range [][]int64{{-5000000, -1000000, 3000000, 5000000}, {-10000000, 3000000, 4000000, 5000000, 2000000}} {
		r, err := interest.InternalRateOfReturn(f)
		fmt.Printf("%v: IRR %.4f%% %v\n", f, r*100, err)
	}
}
