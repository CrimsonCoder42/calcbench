package bloom

import (
	"fmt"
	"math"
)

// FalsePositiveRate returns the predicted probability that a Bloom filter of bitCount bits and
// hashCount hash functions, holding items distinct items, reports an item it never held:
// (1 − e^(−k·n/m))^k. 1 − e^(−x) is computed as −expm1(−x), accurate when k·n/m is small.
// items = 0 gives 0.
// Errors: ErrBitCount when bitCount < 1; ErrHashCount when hashCount < 1.
//
// ID:      ALG-BLOOM-002
// Time:    O(1)
// Space:   O(1)
// Source:  Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (f = (1 − e^(−kn/m))^k)
// Source:  Kirsch and Mitzenmacher, Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2
// Tags:    bloom-filter, false-positive-rate, probability
func FalsePositiveRate(bitCount, items uint64, hashCount int) (float64, error) {
	if bitCount < 1 {
		return 0, fmt.Errorf("%w: bitCount = %d", ErrBitCount, bitCount)
	}
	if hashCount < 1 {
		return 0, fmt.Errorf("%w: hashCount = %d", ErrHashCount, hashCount)
	}
	k := float64(hashCount)
	x := k * float64(items) / float64(bitCount) // finite: every factor is finite and m ≥ 1
	return math.Pow(-math.Expm1(-x), k), nil
}
