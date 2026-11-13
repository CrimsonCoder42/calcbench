package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.7.1 critical values (3 decimals): u = 0.975 → 1.960, 0.995 → 2.576, 0.999 → 3.090, and
// 0.025 → −1.960 by symmetry. Full precision: Φ⁻¹(0.975) = 1.959963984540054 (the round trip below
// checks every digit through NormalCDF).
func TestNormalQuantile_sourceExample(t *testing.T) {
	for _, c := range []struct{ u, want float64 }{{0.975, 1.960}, {0.995, 2.576}, {0.999, 3.090}, {0.025, -1.960}} {
		if got, err := NormalQuantile(c.u, 0, 1); err != nil || math.Abs(got-c.want) > 0.0005 {
			t.Errorf("NormalQuantile(%v, 0, 1) = (%v, %v), want %v to 3 decimals", c.u, got, err, c.want)
		}
	}
	if got, err := NormalQuantile(0.975, 0, 1); err != nil || !within(got, 1.959963984540054, 1e-15) {
		t.Errorf("NormalQuantile(0.975) = (%v, %v), want 1.959963984540054", got, err)
	}
}

func TestNormalQuantile_cases(t *testing.T) {
	tests := []struct {
		name              string
		u, mean, sd, want float64
	}{
		{"median", 0.5, 0, 1, 0},
		{"scaled median", 0.5, 7, 3, 7},
		{"scaled", 0.975, 10, 2, 10 + 2*1.959963984540054},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NormalQuantile(tt.u, tt.mean, tt.sd); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("NormalQuantile(%v, %v, %v) = (%v, %v), want %v", tt.u, tt.mean, tt.sd, got, err, tt.want)
			}
		})
	}
	// The extreme representable tail is finite: Φ⁻¹(5e-324) ≈ −38.4.
	if got, err := NormalQuantile(5e-324, 0, 1); err != nil || got > -38 || got < -39 {
		t.Errorf("NormalQuantile(5e-324) = (%v, %v), want ≈ −38.4", got, err)
	}
}

func TestNormalQuantile_errors(t *testing.T) {
	for _, u := range []float64{0, 1, -0.1, 1.1} {
		_, err := NormalQuantile(u, 0, 1)
		requireErr(t, "u outside (0, 1)", err, ErrOutOfRange)
	}
	_, err := NormalQuantile(nan, 0, 1)
	requireErr(t, "u NaN", err, ErrNotFinite)
	_, err = NormalQuantile(0.5, 0, 0)
	requireErr(t, "sd = 0", err, ErrOutOfRange)
	_, err = NormalQuantile(0.999, math.MaxFloat64, math.MaxFloat64)
	requireErr(t, "quantile overflows", err, ErrOverflow)
}

// Round trip, all three PPND16 regions: for u ≤ ½, NormalCDF(Φ⁻¹(u)) = u; for u > ½, NormalSurvival of it
// is 1 − u (exact for these u). Tolerance 1e-14·(1 + x²) relative: PPND16 is good to about 1e-16, and the
// relative condition number of Φ in the tail is about x².
func TestNormalQuantile_roundTrip(t *testing.T) {
	rng := newTestRand()
	for trial := range 2000 {
		var u float64
		switch trial % 3 {
		case 0:
			u = 0.075 + 0.85*rng.Float64() // region 1
		case 1:
			u = math.Pow(10, -12*rng.Float64()) * 0.075 // region 2
		default:
			u = math.Pow(10, -300*rng.Float64()) * 1e-11 // region 3
		}
		if trial%2 == 1 && trial%3 != 2 && u < 0.5 {
			// The upper half (regions 1 and 2 only: 1 − u of a region-3 u rounds to 1). 1 − u of the
			// stored u is exact (Sterbenz), and it is what both NormalQuantile and the check use.
			u = 1 - u
		}
		x, err := NormalQuantile(u, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		tol := 1e-14 * (1 + x*x)
		if u <= 0.5 {
			if back, _ := NormalCDF(x, 0, 1); !within(back, u, tol) {
				t.Fatalf("trial %d: NormalCDF(NormalQuantile(%v)) = %v (x = %v)", trial, u, back, x)
			}
		} else if back, _ := NormalSurvival(x, 0, 1); !within(back, 1-u, tol) {
			t.Fatalf("trial %d: NormalSurvival(NormalQuantile(%v)) = %v, want %v (x = %v)", trial, u, back, 1-u, x)
		}
	}
}

func TestNormalQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NormalQuantile", func() { sinkFloat, _ = NormalQuantile(0.975, 0, 1) })
}

func BenchmarkNormalQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NormalQuantile(0.975, 0, 1)
	}
}
