package series

import (
	"errors"
	"math"
	"testing"
)

// Hand-worked from Ljung & Box 1978 §2, Q = n(n + 2)Σ r_k²/(n − k): 1, 2, 3, 4 (r_1 = 0.25, r_2 = −0.3).
// h = 1: Q = 4·6·0.0625/3 = 0.5, df 1, p = P(χ²_1 > 0.5) = erfc(0.5) = 0.4795001221869535.
// h = 2: Q = 24·(0.0625/3 + 0.09/2) = 1.58, df 2, p = e^(−0.79) = 0.45384479528235583.
// Tolerance 1e-13 (the p-value's incomplete-gamma evaluation).
func TestLjungBox_handWorked(t *testing.T) {
	tests := []struct {
		maxLag, fitted int
		q              float64
		df             int
		p              float64
	}{
		{1, 0, 0.5, 1, 0.4795001221869535},
		{2, 0, 1.58, 2, 0.45384479528235583},
		{2, 1, 1.58, 1, math.Erfc(math.Sqrt(0.79))},
	}
	for _, tc := range tests {
		q, df, p, err := LjungBox([]float64{1, 2, 3, 4}, tc.maxLag, tc.fitted)
		if err != nil || !closeRel(q, tc.q, 1e-14) || df != tc.df || !closeRel(p, tc.p, 1e-13) {
			t.Errorf("h=%d fitted=%d: = (%v, %v, %v, %v), want (%v, %v, %v)", tc.maxLag, tc.fitted, q, df, p, err, tc.q, tc.df, tc.p)
		}
	}
}

// The hand-worked series scaled by the exact power of two 2⁻¹⁰⁶⁰ (all subnormal) shares acfPrepare's
// rescaling path and must give the same Q (scale invariance); before the cap the factor overflowed to +Inf.
func TestLjungBox_subnormal(t *testing.T) {
	values := []float64{0x1p-1060, 0x1p-1059, 3 * 0x1p-1060, 0x1p-1058}
	q, df, p, err := LjungBox(values, 2, 0)
	if err != nil || !closeRel(q, 1.58, 1e-14) || df != 2 || !closeRel(p, 0.45384479528235583, 1e-13) {
		t.Errorf("= (%v, %v, %v, %v), want (1.58, 2, 0.45384479528235583)", q, df, p, err)
	}
}

// An AR(1) with φ = 0.6 is far from white noise: at n = 2000 and h = 10, Q is in the thousands and the
// p-value underflows toward 0 (bound 1e-10, far above what the fixed seed produces).
func TestLjungBox_ar1Rejects(t *testing.T) {
	q, df, p, err := LjungBox(ar1Series(2000, 0, 0.6), 10, 0)
	if err != nil || df != 10 || q < 500 || p > 1e-10 {
		t.Errorf("= (%v, %v, %v, %v), want large Q, p ≈ 0", q, df, p, err)
	}
}

// White noise (φ = 0): Q ≈ χ²_10 (mean 10). With the fixed seed p is not small; the bound 1e-4 fails with
// probability 10⁻⁴ under the null, and the seed is fixed, so the test is deterministic.
func TestLjungBox_whiteNoise(t *testing.T) {
	_, _, p, err := LjungBox(ar1Series(2000, 0, 0), 10, 0)
	if err != nil || p < 1e-4 {
		t.Errorf("p = %v (%v), want not small", p, err)
	}
}

func TestLjungBox_errors(t *testing.T) {
	tests := []struct {
		name           string
		values         []float64
		maxLag, fitted int
		want           error
	}{
		{"empty", nil, 1, 0, ErrTooFew},
		{"NaN", []float64{1, math.NaN(), 3}, 1, 0, ErrNotFinite},
		{"maxLag 0", []float64{1, 2, 3}, 0, 0, ErrOutOfRange},
		{"maxLag = n", []float64{1, 2, 3}, 3, 0, ErrOutOfRange},
		{"negative fitted", []float64{1, 2, 3}, 1, -1, ErrOutOfRange},
		{"fitted = maxLag", []float64{1, 2, 3}, 2, 2, ErrOutOfRange},
		{"constant", []float64{4, 4, 4}, 1, 0, ErrZeroVariance},
	}
	for _, tc := range tests {
		if _, _, _, err := LjungBox(tc.values, tc.maxLag, tc.fitted); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestLjungBox_noAlloc(t *testing.T) {
	values := ar1Series(500, 0, 0.3)
	if a := testing.AllocsPerRun(10, func() { _, _, _, _ = LjungBox(values, 10, 0) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkLjungBox(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.2)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _, _ = LjungBox(values, 20, 0)
			}
		})
	}
}
