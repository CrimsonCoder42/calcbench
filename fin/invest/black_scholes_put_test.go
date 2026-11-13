package invest

import (
	"math"
	"testing"
)

// Hand-checked with eq. (27) of Black and Scholes 1973 in double precision: S $42, K $40, r 10 %,
// σ 20 %, T 0.5 → 80.8599 cents → 81; parity, eq. (25): C − P = S − K·e^(−rT) = 395.08 cents.
func TestBlackScholesPut_handChecked(t *testing.T) {
	got, err := BlackScholesPut(4200, 4000, 0.10, 0.20, 0.5)
	if err != nil || got != 81 {
		t.Errorf("BlackScholesPut = %d, %v; want 81", got, err)
	}
	call, _ := BlackScholesCall(4200, 4000, 0.10, 0.20, 0.5)
	if fwd := 4200 - 4000*math.Exp(-0.05); math.Abs(float64(call-got)-fwd) > 1+1e-6 {
		t.Errorf("parity: C − P = %d, want %v", call-got, fwd)
	}
}

func TestBlackScholesPut_limits(t *testing.T) {
	for _, c := range []struct {
		name          string
		s, k          int64
		r, sigma, yrs float64
		want          int64
	}{
		{"expiry in the money", 3000, 4000, 0.05, 0.3, 0, 1000}, // boundary (23): c − x
		{"expiry out of the money", 5000, 4000, 0.05, 0.3, 0, 0},
		{"zero volatility", 3000, 4000, 0, 0, 1, 1000},
		{"deep out of the money", 1_000_000, 100, 0.01, 0.1, 0.25, 0},
	} {
		got, err := BlackScholesPut(c.s, c.k, c.r, c.sigma, c.yrs)
		if err != nil || got != c.want {
			t.Errorf("%s: %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestBlackScholesPut_errors(t *testing.T) {
	_, err := BlackScholesPut(0, 4000, 0.1, 0.2, 1)
	wantErr(t, "spot 0", err, ErrDomain)
	_, err = BlackScholesPut(4000, -1, 0.1, 0.2, 1)
	wantErr(t, "strike < 0", err, ErrNegative)
	_, err = BlackScholesPut(maxCents+1, 4000, 0.1, 0.2, 1)
	wantErr(t, "spot > 2^53", err, ErrOutOfRange)
	_, err = BlackScholesPut(4000, 4000, math.NaN(), 0.2, 1)
	wantErr(t, "rate NaN", err, ErrRate)
	_, err = BlackScholesPut(4000, 4000, 0.1, math.Inf(1), 1)
	wantErr(t, "sigma Inf", err, ErrNotFinite)
	_, err = BlackScholesPut(4000, 4000, 0.1, 0.2, -0.5)
	wantErr(t, "years < 0", err, ErrNegative)
	_, err = BlackScholesPut(1, maxCents, -0.9, 0.2, 100) // K·e^90 cents
	wantErr(t, "value overflows", err, ErrOverflow)
}

func TestBlackScholesPut_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = BlackScholesPut(4200, 4000, 0.10, 0.20, 0.5) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBlackScholesPut(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = BlackScholesPut(4200, 4000, 0.10, 0.20, 0.5)
	}
}
