//go:build ignore

package main

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/med/anthropometry"
)

func main() {
	fmt.Printf("hand m: %.15g\n", 86.010*math.Log10(19)-70.041*math.Log10(70)+36.76)
	fmt.Printf("hand f: %.15g\n", 163.205*math.Log10(55)-97.684*math.Log10(65)-78.387)
	// calculator defaults: male 5'10.5", neck 1'7.5", waist 3'1.5", hip 2'10.5"
	m, err := anthropometry.BodyFatPercentNavyMale(37.5, 19.5, 70.5)
	fmt.Printf("default male: %.4f %v\n", m, err)
	f, err := anthropometry.BodyFatPercentNavyFemale(37.5, 34.5, 19.5, 70.5)
	fmt.Printf("default dims female: %.4f %v\n", f, err)
	m2, _ := anthropometry.BodyFatPercentNavyMale(34, 15, 70)
	f2, _ := anthropometry.BodyFatPercentNavyFemale(30, 38, 13, 65)
	fmt.Printf("m 34/15/70: %.4f  f 30/38/13/65: %.4f\n", m2, f2)
}
