package sample

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/CrimsonCoder42/calcbench/algo/prefix"
	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// integer is the built-in integer types, the weights a WeightedPicker accepts (the same set as
// algo/prefix accepts).
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Errors of the weighted picker.
var (
	// ErrNegativeWeight: a weight is below zero.
	ErrNegativeWeight = errors.New("sample: weight is negative")
	// ErrZeroTotal: there are no weights, or they are all zero.
	ErrZeroTotal = errors.New("sample: weights sum to zero")
	// ErrOverflow: the total of the weights does not fit the weight type.
	ErrOverflow = errors.New("sample: total weight does not fit the integer type")
	// ErrNilPicker: Pick was called on a nil *WeightedPicker.
	ErrNilPicker = errors.New("sample: picker is nil")
)

// WeightedPicker picks indices with probability proportional to integer weights. Build it with
// NewWeightedPicker; it is read-only afterwards and safe for concurrent Pick calls with distinct
// *rand.Rand values.
type WeightedPicker[W integer] struct {
	prefix []W // prefix[i] = w₀ + … + w_{i−1}; len = n+1, prefix[n] = total > 0
}

// NewWeightedPicker returns a picker for weights: Pick returns index i with probability
// weights[i] / Σ weights, exactly (integer arithmetic, no rounding). Zero weights are never picked.
// weights is read, never modified; the picker keeps its own prefix sums.
// Errors: ErrNegativeWeight (the first negative index is named), ErrZeroTotal (empty or all zero),
// ErrOverflow (the total does not fit W).
//
// ID:      ALG-SAMPLE-003
// Time:    O(n)   n = len(weights)
// Space:   O(n) for the picker's n+1 prefix sums
// Source:  Devroye, Non-Uniform Random Variate Generation, Springer, 1986, ch. III, §2 (the inversion method for discrete distributions: search in the table of cumulative probabilities)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 (prefix sum arrays), via PrefixSums (ALG-PREFIX-001)
// Tags:    weighted-random, inversion, prefix-sum, discrete-distribution
func NewWeightedPicker[W integer](weights []W) (*WeightedPicker[W], error) {
	var zero W
	for i, w := range weights {
		if w < zero {
			return nil, fmt.Errorf("%w: weights[%d]", ErrNegativeWeight, i)
		}
	}
	sums, err := prefix.PrefixSums(nil, weights)
	if err != nil {
		return nil, fmt.Errorf("%w: weights: %w", ErrOverflow, err)
	}
	if sums[len(sums)-1] == zero {
		return nil, fmt.Errorf("%w: %d weights", ErrZeroTotal, len(weights))
	}
	return &WeightedPicker[W]{prefix: sums}, nil
}

// Pick returns an index i chosen with probability weights[i] / Σ weights: it draws r uniform in
// [0, total) and returns the i with prefix[i] ≤ r < prefix[i+1].
// Errors: ErrNilPicker (p is nil or a zero WeightedPicker), ErrNilRand.
//
// ID:      ALG-SAMPLE-004
// Time:    O(log n)   n = number of weights; one random draw and a binary search
// Space:   O(1) extra, no allocation
// Source:  Devroye, Non-Uniform Random Variate Generation, Springer, 1986, ch. III, §2 (inversion by binary search of the cumulative table), via UpperBound (ALG-SEARCH-002)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search)
// Tags:    weighted-random, inversion, binary-search, discrete-distribution
func (p *WeightedPicker[W]) Pick(rng *rand.Rand) (int, error) {
	if p == nil || len(p.prefix) < 2 {
		return 0, fmt.Errorf("%w: p", ErrNilPicker)
	}
	if rng == nil {
		return 0, fmt.Errorf("%w: rng", ErrNilRand)
	}
	total := p.prefix[len(p.prefix)-1]
	r := W(rng.Uint64N(uint64(total)))
	return search.UpperBound(p.prefix, r) - 1, nil
}
