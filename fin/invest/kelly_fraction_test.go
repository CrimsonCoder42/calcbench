package invest

import (
	"math"
	"testing"
)

func TestKellyFraction_cases(t *testing.T) {
	for _, c := range []struct{ p, b, want float64 }{
		{0.6, 1, 0.2},     // even money: Kelly's p − q
		{0.5, 2, 0.25},    // 0.5 − 0.5/2
		{0.5, 1, 0},       // fair coin, even money: no edge
		{0.4, 1, -0.2},    // disadvantage reported as computed
		{1, 3, 1},         // certain win
		{0, 1, -1},        // certain loss
		{0.25, 4, 0.0625}, // 0.25 − 0.75/4
	} {
		got, err := KellyFraction(c.p, c.b)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("KellyFraction(%v, %v) = %v, %v; want %v", c.p, c.b, got, err, c.want)
		}
	}
}

func TestKellyFraction_errors(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1)} {
		_, err := KellyFraction(bad, 1)
		wantErr(t, "bad p", err, ErrNotFinite)
		_, err = KellyFraction(0.5, bad)
		wantErr(t, "bad b", err, ErrNotFinite)
	}
	for _, bad := range []float64{-0.1, 1.1} {
		_, err := KellyFraction(bad, 1)
		wantErr(t, "p outside [0, 1]", err, ErrDomain)
	}
	for _, bad := range []float64{0, -1} {
		_, err := KellyFraction(0.5, bad)
		wantErr(t, "b ≤ 0", err, ErrDomain)
	}
	_, err := KellyFraction(0.5, 1e-320)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestKellyFraction_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = KellyFraction(0.6, 1) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkKellyFraction(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = KellyFraction(0.6, 1)
	}
}
