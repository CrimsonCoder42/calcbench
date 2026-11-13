package prob

import (
	"math"
	"math/rand/v2"
	"testing"
)

// e-Handbook §1.3.6.7.4: the quantile at prob equals the printed critical value to 3 decimals (tolerance
// 0.0005 + 1e-9, the table's rounding). Includes 7.815 (ν = 3) and 12.592 (ν = 6), the 0.05 critical
// values of 4- and 7-category goodness-of-fit tests.
func TestChiSquareQuantile_sourceTable(t *testing.T) {
	for _, tc := range chiSquareTable {
		got, err := ChiSquareQuantile(tc.prob, tc.df)
		if err != nil {
			t.Fatalf("ChiSquareQuantile(%v, %v): %v", tc.prob, tc.df, err)
		}
		if math.Abs(got-tc.x) > 0.0005+1e-9 {
			t.Errorf("ChiSquareQuantile(%v, %v) = %v, want %v ± 0.0005", tc.prob, tc.df, got, tc.x)
		}
	}
}

// ν = 2 has the closed form x = −2 ln(1 − u) (from P(1, z) = 1 − e^(−z)). Tolerance 1e-11 relative: near
// 0 the CDF's own error is about 1e-12 (exp of a large logarithm), and the root inherits it.
func TestChiSquareQuantile_closedForm(t *testing.T) {
	for _, u := range []float64{1e-300, 1e-12, 0.001, 0.25, 0.5, 0.75, 0.95, 0.999999, 1 - 0x1p-53} {
		got, err := ChiSquareQuantile(u, 2)
		if want := -2 * math.Log1p(-u); err != nil || !within(got, want, 1e-11) {
			t.Errorf("ChiSquareQuantile(%v, 2) = (%v, %v), want %v", u, got, err, want)
		}
	}
}

// Round trip on random (u, df), fixed seed, 300 cases: F(Q(u)) = u (or 1 − F = 1 − u above ½) within 1e-12
// relative, the CDF's own accuracy times the condition of the root.
func TestChiSquareQuantile_roundTrip(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		df := math.Exp(rng.Float64()*math.Log(1e5/0.5)) * 0.5 // log-uniform in [0.5, 1e5]
		u := rng.Float64()
		if u == 0 {
			continue
		}
		x, err := ChiSquareQuantile(u, df)
		if err != nil {
			t.Fatalf("ChiSquareQuantile(%v, %v): %v", u, df, err)
		}
		if u <= 0.5 {
			p, _ := ChiSquareCDF(x, df)
			if !within(p, u, 1e-12) {
				t.Errorf("F(Q(%v; %v)) = %v", u, df, p)
			}
		} else {
			q, _ := ChiSquareSurvival(x, df)
			if !within(q, 1-u, 1e-11) {
				t.Errorf("1 − F(Q(%v; %v)) = %v, want %v", u, df, q, 1-u)
			}
		}
	}
}

func TestChiSquareQuantile_cases(t *testing.T) {
	// A quantile below the smallest float64 is 0 (ν = 0.5, u = 1e-300: x ≈ 2·(u·Γ(1.25))⁴ ≈ 1e-1200).
	if got, err := ChiSquareQuantile(1e-300, 0.5); err != nil || got != 0 {
		t.Errorf("ChiSquareQuantile(1e-300, 0.5) = (%v, %v), want 0 (underflow)", got, err)
	}
	if got, err := ChiSquareQuantile(0, 5); err != nil || got != 0 {
		t.Errorf("ChiSquareQuantile(0, 5) = (%v, %v), want 0", got, err)
	}
	// Far tails at small and large df stay finite and positive.
	for _, c := range []struct{ u, df float64 }{{1e-30, 0.5}, {1 - 1e-15, 0.5}, {1e-10, 1e6}, {1 - 1e-12, 1e6}} {
		got, err := ChiSquareQuantile(c.u, c.df)
		if err != nil || !(got > 0) || math.IsInf(got, 0) {
			t.Errorf("ChiSquareQuantile(%v, %v) = (%v, %v), want finite > 0", c.u, c.df, got, err)
		}
	}
}

func TestChiSquareQuantile_errors(t *testing.T) {
	_, err := ChiSquareQuantile(nan, 3)
	requireErr(t, "ChiSquareQuantile(NaN, 3)", err, ErrNotFinite)
	_, err = ChiSquareQuantile(0.5, inf)
	requireErr(t, "ChiSquareQuantile(0.5, Inf)", err, ErrNotFinite)
	_, err = ChiSquareQuantile(-0.1, 3)
	requireErr(t, "ChiSquareQuantile(-0.1, 3)", err, ErrProbability)
	_, err = ChiSquareQuantile(1.1, 3)
	requireErr(t, "ChiSquareQuantile(1.1, 3)", err, ErrProbability)
	_, err = ChiSquareQuantile(1, 3)
	requireErr(t, "ChiSquareQuantile(1, 3)", err, ErrOutOfRange)
	_, err = ChiSquareQuantile(0.5, 0)
	requireErr(t, "ChiSquareQuantile(0.5, 0)", err, ErrOutOfRange)
	_, err = ChiSquareQuantile(0.5, 1e7)
	requireErr(t, "ChiSquareQuantile(0.5, 1e7)", err, ErrOutOfRange)
}

func TestChiSquareQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ChiSquareQuantile", func() { sinkFloat, _ = ChiSquareQuantile(0.95, 3) })
}

func BenchmarkChiSquareQuantile(b *testing.B) {
	for _, df := range []float64{3, 100, 1e4, 1e6} {
		b.Run(fmtDF(df), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = ChiSquareQuantile(0.95, df)
			}
		})
	}
}
