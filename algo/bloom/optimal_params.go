package bloom

import (
	"fmt"
	"math"
)

// OptimalParams returns the bit count m and hash count k of the smallest Bloom filter predicted to hold
// expectedItems distinct items at falsePositiveRate p: m = ⌈−n·ln p / (ln 2)²⌉ and k = round((m/n)·ln 2),
// at least 1 and at most MaxHashCount. With the real optimum k the predicted rate is exactly p; rounding k to an integer leaves it
// within 1% of p (FalsePositiveRate(m, n, k) ≤ 1.01·p for every n and p tested).
// Errors: ErrItems when expectedItems < 1; ErrProbability when p is NaN or not in (0, 1); ErrTooLarge
// when m > MaxBits.
//
// ID:      ALG-BLOOM-003
// Time:    O(1)
// Space:   O(1)
// Source:  Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (k = ln 2 · m/n minimises the rate) and §2.2 (rate ≤ ε requires m ≥ n·log₂(1/ε)/ln 2)
// Tags:    bloom-filter, sizing, optimal-parameters
func OptimalParams(expectedItems uint64, falsePositiveRate float64) (bitCount uint64, hashCount int, err error) {
	if expectedItems < 1 {
		return 0, 0, fmt.Errorf("%w: expectedItems = %d", ErrItems, expectedItems)
	}
	p := falsePositiveRate
	if !(p > 0 && p < 1) { // also rejects NaN
		return 0, 0, fmt.Errorf("%w: falsePositiveRate = %v", ErrProbability, p)
	}
	n := float64(expectedItems)
	m := math.Ceil(-n * math.Log(p) / (math.Ln2 * math.Ln2)) // > 0; +Inf only past MaxBits anyway
	if m > MaxBits {
		return 0, 0, fmt.Errorf("%w: %v items at rate %v need %v bits", ErrTooLarge, expectedItems, p, m)
	}
	k := min(MaxHashCount, max(1, int(math.Round(m/n*math.Ln2)))) // ≤ m too: (m/n)·ln 2 < m
	return uint64(m), k, nil
}
