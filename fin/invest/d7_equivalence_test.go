package invest

import (
	"math"
	"testing"
)

// D7 (doc/numpy_build_2026-11-13/BRIEF.md): normalCDF now calls NUM-PROB-048 prob.NormalCDF(x, 0, 1).
// oldNormalCDF is the implementation it replaced, copied verbatim (the original is in
// _delete/fin/invest/black_scholes_call.go). prob.NormalCDF computes ½·erfc(−z/√2) with z = (x−0)/1 = x
// exactly, so the new path must equal the old one bit for bit; ±Inf (which prob.NormalCDF rejects) is
// handled by normalCDF's limits, which are the values ½·erfc gives there.
func oldNormalCDF(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

func TestD7_normalCDFMatchesOld(t *testing.T) {
	xs := []float64{math.Inf(-1), -1e308, -40, -38.5, -10, -8.3, -5, -1.96, -1, -0.5, -1e-300, 0, math.Copysign(0, -1),
		1e-300, 0.1, 0.5, 1, 1.96, 3, 5, 8.3, 10, 38.5, 40, 1e308, math.Inf(1)}
	for i := -400; i <= 400; i++ {
		xs = append(xs, float64(i)/37)
	}
	for _, x := range xs {
		if got, want := normalCDF(x), oldNormalCDF(x); got != want {
			t.Errorf("normalCDF(%v) = %v, old %v", x, got, want)
		}
	}
}

func TestD7_blackScholesMatchesOld(t *testing.T) {
	type in struct {
		s, k            int64
		r, sigma, years float64
	}
	cases := []in{{4200, 4000, 0.10, 0.20, 0.5}, {4000, 4000, 0, 0.2, 1}, {1, 1 << 53, 0.05, 0.3, 2}, {1 << 53, 1, 0.05, 0.3, 2},
		{4200, 4000, 0.10, 1e300, 0.5}, {4200, 4000, 0.10, math.MaxFloat64, 1}, {4200, 4000, -0.5, 0.2, 30}, {4200, 4000, 0.1, 0, 1},
		{4200, 4000, 0.1, 0.2, 0}, {100000, 50, 0.02, 0.01, 0.01}, {50, 100000, 0.02, 0.01, 0.01}}
	for _, s := range []int64{1000, 4000, 4200, 9000} {
		for _, sigma := range []float64{0.01, 0.2, 0.8, 3} {
			for _, yrs := range []float64{0.01, 0.5, 5} {
				cases = append(cases, in{s, 4000, 0.05, sigma, yrs})
			}
		}
	}
	for _, c := range cases {
		sp, kDisc, d1, d2, degenerate, err := blackScholesTerms(c.s, c.k, c.r, c.sigma, c.years)
		var oldCall, oldPut int64
		var oldCallErr, oldPutErr error
		if err == nil {
			if degenerate {
				oldCall, oldCallErr = roundCents(math.Max(sp-kDisc, 0))
				oldPut, oldPutErr = roundCents(math.Max(kDisc-sp, 0))
			} else {
				oldCall, oldCallErr = roundCents(sp*oldNormalCDF(d1) - kDisc*oldNormalCDF(d2))
				oldPut, oldPutErr = roundCents(-sp*oldNormalCDF(-d1) + kDisc*oldNormalCDF(-d2))
			}
		} else {
			oldCallErr, oldPutErr = err, err
		}
		call, callErr := BlackScholesCall(c.s, c.k, c.r, c.sigma, c.years)
		put, putErr := BlackScholesPut(c.s, c.k, c.r, c.sigma, c.years)
		if call != oldCall || (callErr == nil) != (oldCallErr == nil) {
			t.Errorf("%+v: call %d (%v), old %d (%v)", c, call, callErr, oldCall, oldCallErr)
		}
		if put != oldPut || (putErr == nil) != (oldPutErr == nil) {
			t.Errorf("%+v: put %d (%v), old %d (%v)", c, put, putErr, oldPut, oldPutErr)
		}
	}
}
