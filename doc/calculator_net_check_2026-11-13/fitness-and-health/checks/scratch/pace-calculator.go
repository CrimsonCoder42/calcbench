//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/life/fitness"
)

func main() {
	p, err := fitness.Pace(3025, 5)
	pk, _ := fitness.Pace(3025, 5*1.609344)
	fmt.Printf("50:25 over 5 mi: %.3f s/mi, %.3f s/km %v\n", p, pk, err)
	p2, err := fitness.Pace(2700, 10)
	fmt.Printf("45:00 over 10 km: %.3f s/km %v\n", p2, err)
	p3, err := fitness.Pace(3*3600+30*60+15, 42.195)
	fmt.Printf("3:30:15 over 42.195 km: %.3f s/km %v\n", p3, err)
}
