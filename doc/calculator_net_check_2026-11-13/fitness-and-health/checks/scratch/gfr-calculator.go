//go:build ignore

package main

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/med/renal"
)

func main() {
	c, err := renal.EGFRCKDEPI2021(0.9, 50, false)
	m, err2 := renal.EGFRMDRD(0.9, 50, false, false)
	fmt.Printf("m 50 y 0.9: CKD-EPI 2021 %.3f %v; MDRD %.3f %v\n", c, err, m, err2)
	c, _ = renal.EGFRCKDEPI2021(1.4, 70, true)
	m, _ = renal.EGFRMDRD(1.4, 70, true, true)
	fmt.Printf("f 70 y 1.4 black: CKD-EPI 2021 %.3f; MDRD %.3f\n", c, m)
	s, err := renal.EGFRSchwartzBedside(110, 0.9, 6)
	fmt.Printf("child 110 cm 0.9: Schwartz %.3f %v\n", s, err)
}
