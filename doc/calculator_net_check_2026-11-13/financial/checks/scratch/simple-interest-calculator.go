//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

func main() {
	a, _ := interest.SimpleInterest(2000000, 0.03, 10)
	b, _ := interest.SimpleInterest(123456, 0.0725, 18.0/12)
	fmt.Println(a, 2000000+a, b, 123456+b)
}
