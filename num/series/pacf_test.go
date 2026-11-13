package series

import (
	"errors"
	"math"
	"testing"
)

// yuleWalkerLastOracle solves the order-k Yule–Walker system R·φ = (r_1…r_k), R_ij = r_|i−j|, by Gaussian
// elimination with partial pivoting and returns φ_kk: the definition of the PACF (BJRL §3.2.6), with no
// recursion, as a slow independent oracle.
func yuleWalkerLastOracle(r []float64, k int) float64 {
	a := make([][]float64, k)
	for i := range a {
		a[i] = make([]float64, k+1)
		for j := 0; j < k; j++ {
			d := i - j
			if d < 0 {
				d = -d
			}
			a[i][j] = r[d]
		}
		a[i][k] = r[i+1]
	}
	for c := 0; c < k; c++ {
		p := c
		for i := c + 1; i < k; i++ {
			if math.Abs(a[i][c]) > math.Abs(a[p][c]) {
				p = i
			}
		}
		a[c], a[p] = a[p], a[c]
		for i := c + 1; i < k; i++ {
			f := a[i][c] / a[c][c]
			for j := c; j <= k; j++ {
				a[i][j] -= f * a[c][j]
			}
		}
	}
	x := make([]float64, k)
	for i := k - 1; i >= 0; i-- {
		s := a[i][k]
		for j := i + 1; j < k; j++ {
			s -= a[i][j] * x[j]
		}
		x[i] = s / a[i][i]
	}
	return x[k-1]
}

// Hand-worked (Durbin 1960 §2 recursion): 1, 2, 3, 4 → r = 1, 0.25, −0.3, −0.45; φ_11 = 0.25;
// φ_22 = (−0.3 − 0.25·0.25)/(1 − 0.25·0.25) = −0.38666…; φ_21 = 0.25 − φ_22·0.25 = 0.34666…;
// φ_33 = (−0.45 − φ_21·(−0.3) − φ_22·0.25)/(1 − φ_21·0.25 − φ_22·(−0.3)) = −0.2493333…/0.7973333…
// Tolerance 1e-14 (a few roundings in the recursion).
func TestPACF_handWorked(t *testing.T) {
	got, err := PACF(nil, []float64{1, 2, 3, 4}, 3)
	phi21 := 0.25 + 0.25*(0.3625/0.9375)
	phi22 := -0.3625 / 0.9375
	num := -0.45 - phi21*(-0.3) - phi22*0.25
	den := 1 - phi21*0.25 - phi22*(-0.3)
	want := []float64{1, 0.25, -0.38666666666666666, num / den}
	if err != nil || len(got) != 4 {
		t.Fatalf("= (%v, %v)", got, err)
	}
	for i := range want {
		if !closeRel(got[i], want[i], 1e-14) {
			t.Errorf("[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

// Oracle: PACF against the Yule–Walker solve on generated series, 50 series × lags 1…8. Tolerance 1e-9
// (the elimination and the recursion round differently; the systems are well conditioned).
func TestPACF_oracleYuleWalker(t *testing.T) {
	rng := newRand()
	for trial := 0; trial < 50; trial++ {
		n := 20 + rng.IntN(200)
		values := make([]float64, n)
		for i := range values {
			values[i] = rng.NormFloat64()
			if i > 0 {
				values[i] += 0.5 * values[i-1]
			}
		}
		const m = 8
		r, _ := ACF(nil, values, m)
		got, err := PACF(nil, values, m)
		if err != nil {
			t.Fatal(err)
		}
		for k := 1; k <= m; k++ {
			if w := yuleWalkerLastOracle(r, k); math.Abs(got[k]-w) > 1e-9 {
				t.Errorf("trial %d lag %d: %v, oracle %v", trial, k, got[k], w)
			}
		}
	}
}

// The PACF of an AR(1) is φ at lag 1 and 0 beyond; SE ≈ 1/√n ≈ 0.007 at n = 20000, tolerance 0.04.
func TestPACF_ar1(t *testing.T) {
	const phi = 0.6
	got, err := PACF(nil, ar1Series(20000, 1, phi), 4)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got[1]-phi) > 0.04 {
		t.Errorf("φ_11 = %v, want ≈ %v", got[1], phi)
	}
	for k := 2; k <= 4; k++ {
		if math.Abs(got[k]) > 0.04 {
			t.Errorf("φ_%d%d = %v, want ≈ 0", k, k, got[k])
		}
	}
}

func TestPACF_cases(t *testing.T) {
	if got, err := PACF(nil, []float64{5, 7}, 0); err != nil || len(got) != 1 || got[0] != 1 {
		t.Errorf("maxLag 0: (%v, %v), want [1]", got, err)
	}
	if got, err := PACF(nil, []float64{0, 1}, 1); err != nil || !closeRel(got[1], -0.5, 1e-15) {
		t.Errorf("two values: (%v, %v), want φ_11 = −0.5", got, err)
	}
}

func TestPACF_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		maxLag int
		want   error
	}{
		{"empty", nil, 0, ErrTooFew},
		{"one", []float64{1}, 0, ErrTooFew},
		{"Inf", []float64{1, math.Inf(1), 2}, 1, ErrNotFinite},
		{"negative maxLag", []float64{1, 2, 3}, -1, ErrOutOfRange},
		{"maxLag = n", []float64{1, 2, 3}, 3, ErrOutOfRange},
		{"constant", []float64{2, 2, 2}, 1, ErrZeroVariance},
	}
	for _, tc := range tests {
		if _, err := PACF(nil, tc.values, tc.maxLag); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkPACF(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.5)
		dst := make([]float64, 21)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = PACF(dst, values, 20)
			}
		})
	}
}
