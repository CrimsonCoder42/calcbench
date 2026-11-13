package interest

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Shared test helpers for the package (the first concept's test file holds them).

const bigPrec = 256

// bigOnePlus returns 1 + i exactly (i is a float64, so it is exact in 256 bits).
func bigOnePlus(i float64) *big.Float {
	x := new(big.Float).SetPrec(bigPrec).SetFloat64(i)
	return x.Add(x, new(big.Float).SetPrec(bigPrec).SetInt64(1))
}

// bigPow returns x^n for n ≥ 0 by repeated squaring at 256 bits.
func bigPow(x *big.Float, n int) *big.Float {
	result := new(big.Float).SetPrec(bigPrec).SetInt64(1)
	base := new(big.Float).SetPrec(bigPrec).Set(x)
	for n > 0 {
		if n&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
		n >>= 1
	}
	return result
}

// bigCompound returns (1+i)^n (n may be negative) at 256 bits.
func bigCompound(i float64, n int) *big.Float {
	if n >= 0 {
		return bigPow(bigOnePlus(i), n)
	}
	p := bigPow(bigOnePlus(i), -n)
	return new(big.Float).SetPrec(bigPrec).Quo(new(big.Float).SetPrec(bigPrec).SetInt64(1), p)
}

// bigF returns a 256-bit copy of v.
func bigF(v float64) *big.Float { return new(big.Float).SetPrec(bigPrec).SetFloat64(v) }

// bigI returns a 256-bit copy of v.
func bigI(v int64) *big.Float { return new(big.Float).SetPrec(bigPrec).SetInt64(v) }

// checkCentsOracle applies the accuracy contract (F-5): got must equal the exact value rounded
// half-even, except that a 1-cent difference is accepted when the exact value lies within
// max(1e-6 cent, 1e-14·|x|) of a half cent. The relative term is the float64 evaluation error of
// the closed forms (a few ε times |n·ln(1+i)| ≤ 30 here); without it a value of 10¹² cents could fail
// on a correct implementation. Callers keep |x| ≤ 2⁴³ cents.
func checkCentsOracle(t *testing.T, label string, got int64, exact *big.Float) {
	t.Helper()
	f, _ := exact.Float64()
	floor := math.Floor(f)
	frac := f - floor
	want := int64(math.RoundToEven(f))
	if math.Abs(f) > 1<<43 {
		t.Errorf("%s: test magnitude %.0f cents exceeds the 2^43 contract range", label, f)
	}
	slack := math.Max(1e-6, 1e-14*math.Abs(f))
	if math.Abs(frac-0.5) <= slack {
		if got != want && got != want-1 && got != want+1 {
			t.Errorf("%s: got %d, want %d ±1 (exact %.6f is at a half cent)", label, got, want, f)
		}
		return
	}
	if got != want {
		t.Errorf("%s: got %d, want %d (exact %.6f)", label, got, want, f)
	}
}

// wantErr fails unless err wraps target.
func wantErr(t *testing.T, label string, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("%s: err = %v, want %v", label, err, target)
	}
}

// wantCents fails unless err is nil and got == want.
func wantCents(t *testing.T, label string, got int64, err error, want int64) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error %v", label, err)
		return
	}
	if got != want {
		t.Errorf("%s: got %d cents, want %d", label, got, want)
	}
}

// wantFloat fails unless err is nil and |got − want| ≤ tol.
func wantFloat(t *testing.T, label string, got float64, err error, want, tol float64) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error %v", label, err)
		return
	}
	if math.IsNaN(got) || math.Abs(got-want) > tol {
		t.Errorf("%s: got %.15g, want %.15g ± %g", label, got, want, tol)
	}
}

// wantRel fails unless err is nil and got agrees with want to relative tolerance rel.
func wantRel(t *testing.T, label string, got float64, err error, want, rel float64) {
	t.Helper()
	wantFloat(t, label, got, err, want, rel*math.Max(math.Abs(want), math.SmallestNonzeroFloat64))
}

// wantNoAllocs fails unless f allocates nothing.
func wantNoAllocs(t *testing.T, label string, f func()) {
	t.Helper()
	if a := testing.AllocsPerRun(100, f); a != 0 {
		t.Errorf("%s: %v allocations per run, want 0", label, a)
	}
}

// testRNG is the fixed-seed generator for oracle tests.
func testRNG() *rand.Rand { return rand.New(rand.NewPCG(1, 2)) }

// badRates are the rate inputs every rate-taking function must reject with ErrRate.
var badRates = []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1, -1.5}

var sinkCents int64
var sinkFloat float64

func TestFutureValue_SourceExample(t *testing.T) {
	// Dahlquist & Knight, Principles of Finance (OpenStax 2022) §7.3 Example 1: PV $1,000, I/Y 3, N 4 → $1,125.51.
	got, err := FutureValue(100000, 0.03, 4)
	wantCents(t, "§7.3 Example 1", got, err, 112551)
	// §7.2: $1,000 at 4 % for 10 years → $1,480.24; a $400,000 house at 4 % for 6 years → $506,127.61.
	got, err = FutureValue(100000, 0.04, 10)
	wantCents(t, "§7.2 10 years", got, err, 148024)
	got, err = FutureValue(40000000, 0.04, 6)
	wantCents(t, "§7.2 house", got, err, 50612761)
}

func TestFutureValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		pv      int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate", 12345, 0, 30, 12345},
		{"zero periods", 12345, 0.07, 0, 12345},
		{"one period", 10000, 0.05, 1, 10500},
		{"negative amount", -100000, 0.03, 4, -112551},
		{"zero amount", 0, 0.5, 100, 0},
		{"max exact amount, no growth", maxExactCents, 0.1, 0, maxExactCents},
		{"negative rate", 100000, -0.5, 2, 25000},
	}
	for _, c := range cases {
		got, err := FutureValue(c.pv, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestFutureValue_Errors(t *testing.T) {
	_, err := FutureValue(maxExactCents+1, 0.05, 1)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = FutureValue(-maxExactCents-1, 0.05, 1)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = FutureValue(100, r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = FutureValue(100, math.NaN(), 1)
	wantErr(t, "NaN rate is also ErrNotFinite", err, ErrNotFinite)
	_, err = FutureValue(100, 0.05, -1)
	wantErr(t, "negative periods", err, ErrPeriods)
	_, err = FutureValue(maxExactCents, 0.05, 1)
	wantErr(t, "result above 2^53", err, ErrOverflow)
	_, err = FutureValue(100, 1, 2000)
	wantErr(t, "factor beyond float64", err, ErrOverflow)
}

func TestFutureValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		pv := rng.Int64N(100_000_000) - 50_000_000 // |FV| ≤ 5e7·1.1¹¹⁹ < 2⁴³
		rate := rng.Float64()*0.15 - 0.05
		n := rng.IntN(120)
		got, err := FutureValue(pv, rate, n)
		if err != nil {
			t.Fatalf("FutureValue(%d, %v, %d): %v", pv, rate, n, err)
		}
		exact := new(big.Float).SetPrec(bigPrec).Mul(bigI(pv), bigCompound(rate, n))
		checkCentsOracle(t, "oracle", got, exact)
	}
}

func TestFutureValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "FutureValue", func() { sinkCents, _ = FutureValue(100000, 0.03, 4) })
}

func BenchmarkFutureValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = FutureValue(100000, 0.03, 360)
	}
}
