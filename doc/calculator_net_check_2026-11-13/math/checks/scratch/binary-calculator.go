//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/intmath"
)

func main() {
	fmt.Println(intmath.AddBinaryStrings("10101010", "11001100"))
	fmt.Println(intmath.AddBinaryStrings("1111111111111111111111111111111111111111111111111111111111111111111111", "1"))
	fmt.Println(intmath.AddBinaryStrings("0", "0"))
}
