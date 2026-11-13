package invest

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Black and Scholes 1973 give no numeric example. Hand-checked with eq. (13) in double precision
// (N = ½·erfc(−x/√2)): S $42, K $40, r 10 %, σ 20 %, T 0.5 → 475.9422 cents → 476.
func TestBlackScholesCall_handChecked(t *testing.T) {
	got, err := BlackScholesCall(4200, 4000, 0.10, 0.20, 0.5)
	if err != nil || got != 476 {
		t.Errorf("BlackScholesCall = %d, %v; want 476", got, err)
	}
}

func TestBlackScholesCall_limits(t *testing.T) {
	for _, c := range []struct {
		name          string
		s, k          int64
		r, sigma, yrs float64
		want          int64
	}{
		{"expiry in the money", 5000, 4000, 0.05, 0.3, 0, 1000},  // boundary (8): x − c
		{"expiry out of the money", 3000, 4000, 0.05, 0.3, 0, 0}, // boundary (8): 0
		{"zero volatility", 5000, 4000, 0.05, 0, 1, 1195},        // 5000 − 4000·e^(−0.05) = 1195.08
		{"zero volatility out", 3000, 4000, 0.05, 0, 1, 0},
		{"deep out of the money", 100, 1_000_000, 0.01, 0.1, 0.25, 0},
		{"deep in the money", 1_000_000, 100, 0, 0.2, 1, 999_900},
	} {
		got, err := BlackScholesCall(c.s, c.k, c.r, c.sigma, c.yrs)
		if err != nil || got != c.want {
			t.Errorf("%s: %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// Oracle: bounds max(S − K·e^(−rT), 0) ≤ C ≤ S (to the cent), and put-call parity, eq. (25), within
// 1 cent (both values are rounded), on 300 random inputs.
func TestBlackScholesCall_boundsAndParity(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		s := r.Int64N(1_000_000) + 1
		k := r.Int64N(1_000_000) + 1
		rate := float64(r.IntN(1500)-200) / 10000
		sigma := float64(r.IntN(10000)) / 10000
		yrs := float64(r.IntN(3000)) / 1000
		call, err := BlackScholesCall(s, k, rate, sigma, yrs)
		noErr(t, "call", err)
		put, err := BlackScholesPut(s, k, rate, sigma, yrs)
		noErr(t, "put", err)
		fwd := float64(s) - float64(k)*math.Exp(-rate*yrs)
		if float64(call) < math.Max(fwd, 0)-1 || call > s {
			t.Errorf("S %d K %d r %v σ %v T %v: call %d outside bounds", s, k, rate, sigma, yrs, call)
		}
		if math.Abs(float64(call-put)-fwd) > 1+1e-6 { // two half-cent roundings plus float slack
			t.Errorf("S %d K %d r %v σ %v T %v: C − P = %d, S − K·e^(−rT) = %v", s, k, rate, sigma, yrs, call-put, fwd)
		}
	}
}

// Huge volatility: σ = 1e200 (σ² would overflow; the code never forms it). As σ√T → ∞, d₁ → +∞ and
// d₂ → −∞, so C → S and, by parity, P → K·e^(−rT). Also σ = MaxFloat64 with T = 4 (sd = +Inf).
func TestBlackScholesCall_hugeVolatility(t *testing.T) {
	for _, c := range []struct{ sigma, yrs float64 }{{1e200, 1}, {1e200, 0.5}, {math.MaxFloat64, 4}} {
		call, err := BlackScholesCall(4200, 4000, 0.10, c.sigma, c.yrs)
		if err != nil || call != 4200 {
			t.Errorf("σ %v T %v: call %d, %v; want S = 4200", c.sigma, c.yrs, call, err)
		}
		put, err := BlackScholesPut(4200, 4000, 0.10, c.sigma, c.yrs)
		want := int64(math.RoundToEven(4000 * math.Exp(-0.10*c.yrs)))
		if err != nil || put != want {
			t.Errorf("σ %v T %v: put %d, %v; want K·e^(−rT) = %d", c.sigma, c.yrs, put, err, want)
		}
	}
}

func TestBlackScholesCall_errors(t *testing.T) {
	_, err := BlackScholesCall(0, 4000, 0.1, 0.2, 1)
	wantErr(t, "spot 0", err, ErrDomain)
	_, err = BlackScholesCall(4000, 0, 0.1, 0.2, 1)
	wantErr(t, "strike 0", err, ErrDomain)
	_, err = BlackScholesCall(-1, 4000, 0.1, 0.2, 1)
	wantErr(t, "spot < 0", err, ErrNegative)
	_, err = BlackScholesCall(4000, maxCents+1, 0.1, 0.2, 1)
	wantErr(t, "strike > 2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1} {
		_, err = BlackScholesCall(4000, 4000, bad, 0.2, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1)} {
		_, err = BlackScholesCall(4000, 4000, 0.1, bad, 1)
		wantErr(t, "bad sigma", err, ErrNotFinite)
		_, err = BlackScholesCall(4000, 4000, 0.1, 0.2, bad)
		wantErr(t, "bad years", err, ErrNotFinite)
	}
	_, err = BlackScholesCall(4000, 4000, 0.1, -0.2, 1)
	wantErr(t, "sigma < 0", err, ErrNegative)
	_, err = BlackScholesCall(4000, 4000, 0.1, 0.2, -1)
	wantErr(t, "years < 0", err, ErrNegative)
	_, err = BlackScholesCall(4000, 4000, -0.99, 0.2, 1e6)
	wantErr(t, "discounted strike overflows", err, ErrOverflow)
}

func TestBlackScholesCall_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = BlackScholesCall(4200, 4000, 0.10, 0.20, 0.5) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBlackScholesCall(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = BlackScholesCall(4200, 4000, 0.10, 0.20, 0.5)
	}
}
