package prob

import (
	"math"
	"testing"
)

// chiSquareTable holds e-Handbook §1.3.6.7.4 ("Critical Values of the Chi-Square Distribution", read
// 2026-11-12): the critical value x with P(X < x) = prob for ν degrees of freedom, printed to 3 decimals.
var chiSquareTable = []struct{ df, prob, x float64 }{
	{1, 0.90, 2.706}, {1, 0.95, 3.841}, {1, 0.975, 5.024}, {1, 0.99, 6.635}, {1, 0.999, 10.828},
	{2, 0.90, 4.605}, {2, 0.95, 5.991}, {2, 0.975, 7.378}, {2, 0.99, 9.210}, {2, 0.999, 13.816},
	{3, 0.90, 6.251}, {3, 0.95, 7.815}, {3, 0.975, 9.348}, {3, 0.99, 11.345}, {3, 0.999, 16.266},
	{6, 0.90, 10.645}, {6, 0.95, 12.592}, {6, 0.975, 14.449}, {6, 0.99, 16.812}, {6, 0.999, 22.458},
	{10, 0.90, 15.987}, {10, 0.95, 18.307}, {10, 0.975, 20.483}, {10, 0.99, 23.209}, {10, 0.999, 29.588},
	{30, 0.90, 40.256}, {30, 0.95, 43.773}, {30, 0.975, 46.979}, {30, 0.99, 50.892}, {30, 0.999, 59.703},
	{100, 0.90, 118.498}, {100, 0.95, 124.342}, {100, 0.975, 129.561}, {100, 0.99, 135.807}, {100, 0.999, 149.449},
	// Lower-tail table (same page): P(X < x) = 0.10, 0.05 for ν = 3 and 10.
	{3, 0.10, 0.584}, {3, 0.05, 0.352}, {10, 0.10, 4.865}, {10, 0.05, 3.940},
}

// The table rounds x to 3 decimals, so F at the printed x differs from prob by at most f(x)·0.0005 ≤ 0.0005
// (the density is below 1 for ν ≥ 2 and away from 0 for ν = 1); the check is |F(x) − prob| ≤ 0.0005 + 1e-9.
func TestChiSquareCDF_sourceTable(t *testing.T) {
	for _, tc := range chiSquareTable {
		got, err := ChiSquareCDF(tc.x, tc.df)
		if err != nil {
			t.Fatalf("ChiSquareCDF(%v, %v): %v", tc.x, tc.df, err)
		}
		if math.Abs(got-tc.prob) > 0.0005+1e-9 {
			t.Errorf("ChiSquareCDF(%v, %v) = %v, want %v ± 0.0005", tc.x, tc.df, got, tc.prob)
		}
	}
}

// Closed forms (DLMF 8.4: P(1, z) = 1 − e^(−z), P(2, z) = 1 − e^(−z)(1 + z), P(½, z) = erf √z), so ν = 2, 4
// and 1 are oracles. Tolerance 1e-13 relative: the closed forms themselves round in exp and erf.
func TestChiSquareCDF_closedForms(t *testing.T) {
	for _, x := range []float64{1e-8, 0.01, 0.3, 1, 2, 3.5, 7, 15, 40} {
		cases := []struct {
			df   float64
			want float64
		}{
			{2, -math.Expm1(-x / 2)},
			{4, 1 - math.Exp(-x/2)*(1+x/2)},
			{1, math.Erf(math.Sqrt(x / 2))},
		}
		for _, c := range cases {
			if c.df == 4 && x < 0.3 {
				continue // 1 − e^(−z)(1 + z) cancels to z²/2 below z ≈ 0.15: the oracle itself loses digits
			}
			got, err := ChiSquareCDF(x, c.df)
			if err != nil || !within(got, c.want, 1e-13) {
				t.Errorf("ChiSquareCDF(%v, %v) = (%v, %v), want %v", x, c.df, got, err, c.want)
			}
		}
	}
}

func TestChiSquareCDF_cases(t *testing.T) {
	tests := []struct {
		name      string
		x, df     float64
		want, tol float64
	}{
		{"zero", 0, 3, 0, 0},
		{"negative", -5, 3, 0, 0},
		{"median of df 2 is 2 ln 2", 2 * math.Ln2, 2, 0.5, 1e-15},
		{"far upper tail", 1e4, 3, 1, 0},
		{"tiny lower tail df 2", 1e-300, 2, 5e-301, 1e-12}, // exp(ln z) amplifies the rounding of ln z (692 ulps)
	}
	for _, tc := range tests {
		got, err := ChiSquareCDF(tc.x, tc.df)
		if err != nil || !within(got, tc.want, tc.tol) {
			t.Errorf("%s: ChiSquareCDF(%v, %v) = (%v, %v), want %v", tc.name, tc.x, tc.df, got, err, tc.want)
		}
	}
}

// Even df = 2k: Q(k, z) = e^(−z)·Σ_{j<k} zʲ/j! (DLMF §8.4), so ChiSquareSurvival(x, 2k) = PoissonCDF(k − 1, x/2),
// an independent oracle (NUM-PROB-042 sums the Poisson terms). Large k and x near the mean exercise the
// saddle-point prefactor and both evaluation branches. Tolerance 1e-11 relative on the smaller tail (the
// continued fraction's rounding grows with its length, about √df terms).
func TestChiSquareCDF_poissonOracle(t *testing.T) {
	for _, k := range []int{1, 3, 10, 50, 1000, 100000} {
		df := float64(2 * k)
		for _, r := range []float64{0.3, 0.9, 0.99, 1, 1.01, 1.1, 2} {
			x := df * r
			pois, err := PoissonCDF(k-1, x/2)
			if err != nil {
				t.Fatal(err)
			}
			q, err := ChiSquareSurvival(x, df)
			if err != nil {
				t.Fatal(err)
			}
			p, err := ChiSquareCDF(x, df)
			if err != nil {
				t.Fatal(err)
			}
			if pois < 0.5 {
				if !within(q, pois, 1e-11) {
					t.Errorf("df %v, x %v: Q = %v, Poisson oracle %v", df, x, q, pois)
				}
			} else if !within(p, 1-pois, 1e-11) && math.Abs(p-(1-pois)) > 1e-15 {
				t.Errorf("df %v, x %v: F = %v, Poisson oracle %v", df, x, p, 1-pois)
			}
		}
	}
}

func TestChiSquareCDF_errors(t *testing.T) {
	_, err := ChiSquareCDF(nan, 3)
	requireErr(t, "ChiSquareCDF(NaN, 3)", err, ErrNotFinite)
	_, err = ChiSquareCDF(inf, 3)
	requireErr(t, "ChiSquareCDF(Inf, 3)", err, ErrNotFinite)
	_, err = ChiSquareCDF(1, nan)
	requireErr(t, "ChiSquareCDF(1, NaN)", err, ErrNotFinite)
	_, err = ChiSquareCDF(1, -inf)
	requireErr(t, "ChiSquareCDF(1, -Inf)", err, ErrNotFinite)
	for _, df := range []float64{0, -1, 1e6 + 1} {
		_, err = ChiSquareCDF(1, df)
		requireErr(t, "ChiSquareCDF(1, df)", err, ErrOutOfRange)
	}
	if _, err := ChiSquareCDF(1e6, 1e6); err != nil {
		t.Errorf("ChiSquareCDF at df = 1e6 (the Range edge): %v", err)
	}
}

// Monotone and complementary over a grid: F non-decreasing in x (to 1e-14, the switch between series and
// continued fraction at x = df + 2), F + Q = 1 within 4 ulps of 1.
func TestChiSquareCDF_monotoneAndComplement(t *testing.T) {
	for _, df := range []float64{0.1, 0.5, 1, 2.5, 7, 31, 250, 5e4} {
		prev := 0.0
		for i := 1; i <= 400; i++ {
			x := df * float64(i) / 100
			p, err := ChiSquareCDF(x, df)
			if err != nil {
				t.Fatal(err)
			}
			q, err := ChiSquareSurvival(x, df)
			if err != nil {
				t.Fatal(err)
			}
			if p < prev-1e-14 {
				t.Fatalf("df %v: F(%v) = %v < previous %v", df, x, p, prev)
			}
			if math.Abs(p+q-1) > 4*0x1p-52 {
				t.Fatalf("df %v, x %v: F + Q = %v", df, x, p+q)
			}
			prev = p
		}
	}
}

func TestChiSquareCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ChiSquareCDF", func() { sinkFloat, _ = ChiSquareCDF(7.815, 3) })
	requireNoAlloc(t, "ChiSquareCDF large df", func() { sinkFloat, _ = ChiSquareCDF(1e5, 1e5) })
}

func BenchmarkChiSquareCDF(b *testing.B) {
	for _, df := range []float64{3, 100, 1e4, 1e6} {
		b.Run(fmtDF(df), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = ChiSquareCDF(df, df)
			}
		})
	}
}

// fmtDF names a benchmark case by its degrees of freedom.
func fmtDF(df float64) string {
	switch df {
	case 3:
		return "df=3"
	case 100:
		return "df=1e2"
	case 1e4:
		return "df=1e4"
	}
	return "df=1e6"
}
