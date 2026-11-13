//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/integer"
)

func main() {
	p := func(s string, n, d int64, err error) { fmt.Printf("%s = %d/%d %v\n", s, n, d, err) }
	n, d, e := integer.FractionAdd(2, 7, 3, 8)
	p("2/7+3/8", n, d, e)
	n, d, e = integer.FractionDiv(5, 12, 10, 9)
	p("5/12÷10/9", n, d, e)
	n, d, e = integer.FractionSub(1, 6, 5, 6)
	p("1/6-5/6", n, d, e)
	n, d, e = integer.FractionAdd(-11, 4, 26, 7)
	p("-2 3/4 + 3 5/7", n, d, e)
	n, d, e = integer.FractionReduce(21, 98)
	p("21/98", n, d, e)
}
