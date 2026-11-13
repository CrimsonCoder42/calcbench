package sample

import (
	"fmt"
	"math/rand/v2"
)

// RejectionSample returns the first candidate x = propose(rng) that accept(x, u) accepts, u being a
// fresh uniform draw in [0, 1) per try. To sample a density f from a proposal density g with f ≤ c·g,
// write accept as u < f(x) / (c·g(x)); the expected number of tries is then c. accept may also ignore u
// and test membership (uniform on a region from uniform on a box). After maxTries rejections it returns
// ErrMaxTries: a bad bound or an empty region becomes an error, not an endless loop.
// Errors: ErrNilRand, ErrNilFunc (propose or accept), ErrNonPositive (maxTries < 1), ErrMaxTries.
//
// ID:      ALG-SAMPLE-005
// Time:    O(t · (P + A))   t = tries, expected c (geometric), at most maxTries; P, A = cost of propose and accept
// Space:   O(1) extra
// Source:  von Neumann, "Various techniques used in connection with random digits", in Monte Carlo Method, National Bureau of Standards Applied Mathematics Series 12, 1951, pp. 36–38
// Source:  Devroye, Non-Uniform Random Variate Generation, Springer, 1986, ch. II, §3 (the rejection method)
// Tags:    rejection-sampling, accept-reject, monte-carlo, random-variate
func RejectionSample[T any](rng *rand.Rand, propose func(*rand.Rand) T, accept func(candidate T, u float64) bool, maxTries int) (T, error) {
	var zero T
	switch {
	case rng == nil:
		return zero, fmt.Errorf("%w: rng", ErrNilRand)
	case propose == nil:
		return zero, fmt.Errorf("%w: propose", ErrNilFunc)
	case accept == nil:
		return zero, fmt.Errorf("%w: accept", ErrNilFunc)
	case maxTries < 1:
		return zero, fmt.Errorf("%w: maxTries = %d", ErrNonPositive, maxTries)
	}
	for range maxTries {
		x := propose(rng)
		if accept(x, rng.Float64()) {
			return x, nil
		}
	}
	return zero, fmt.Errorf("%w: %d tries", ErrMaxTries, maxTries)
}
