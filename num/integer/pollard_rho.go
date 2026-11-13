package integer

// rhoBatch is the number of |x − y| products multiplied together before one gcd (Brent 1980, §7).
const rhoBatch = 128

// rhoMaxC is the last increment c tried before findFactor falls back to trial division. Rho with
// f(x) = x² + c has never been seen to fail for every c in 1..255 on a uint64; the bound only guarantees
// that the call ends.
const rhoMaxC = 255

// findFactor returns a divisor d of the composite m with 1 < d < m. m must be odd, composite and have no
// prime factor below 1031. It runs Pollard's rho with Brent's cycle detection for c = 1, 2, …, rhoMaxC and,
// should every c fail, trial division.
func findFactor(m uint64) uint64 {
	for c := uint64(1); c <= rhoMaxC; c++ {
		if d := rhoBrent(m, c); d != m {
			return d
		}
	}
	return trialFactor(m)
}

// rhoBrent runs Pollard's rho on m with f(x) = x² + c mod m and x₀ = 2, detecting the cycle by Brent's
// power-of-two method and taking one gcd per rhoBatch products. It returns a divisor of m greater than 1,
// which is m itself when this c fails. Pollard 1975; Brent 1980, §7 (the batched-gcd algorithm).
// Expected O(p^(1/2)) steps, p = the least prime factor of m (heuristic); O(1) space.
func rhoBrent(m, c uint64) uint64 {
	f := func(v uint64) uint64 { return addMod(mulMod(v, v, m), c, m) }
	y, q, g := uint64(2), uint64(1), uint64(1)
	var x, ys uint64
	for r := 1; g == 1; r *= 2 {
		x = y
		for range r {
			y = f(y)
		}
		for k := 0; k < r && g == 1; k += rhoBatch {
			ys = y
			for range min(rhoBatch, r-k) {
				y = f(y)
				q = mulMod(q, absDiff(x, y), m)
			}
			g = GCD(q, m)
		}
	}
	if g == m {
		// The batch overshot (or q reached 0): replay it one step at a time from its start.
		for {
			ys = f(ys)
			if g = GCD(absDiff(x, ys), m); g > 1 {
				break
			}
		}
	}
	return g
}

// absDiff returns |a − b| for unsigned a and b.
func absDiff(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

// trialFactor returns the least divisor d > 1 of m by odd trial division from 1031 (m odd, with no prime
// factor below 1031). It is the guaranteed-to-end fallback of findFactor: O(√m) divisions, O(1) space.
func trialFactor(m uint64) uint64 {
	for d := uint64(1031); d <= m/d; d += 2 {
		if m%d == 0 {
			return d
		}
	}
	return m
}
