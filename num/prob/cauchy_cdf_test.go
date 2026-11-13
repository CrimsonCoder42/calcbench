package prob

import (
	"math"
	"testing"
)

// JKB vol. 1 ch. 16: F(x) = ½ + atan(z)/π; F(1) = 0.75, F(−1) = 0.25; the far lower tail ≈ 1/(π|z|).
func TestCauchyCDF_sourceExample(t *testing.T) {
	cases := []struct{ x, want float64 }{{1, 0.75}, {-1, 0.25}, {0, 0.5}, {-1e12, 1 / (math.Pi * 1e12)}}
	for _, c := range cases {
		if got, err := CauchyCDF(c.x, 0, 1); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("CauchyCDF(%v) = (%v, %v), want %v", c.x, got, err, c.want)
		}
	}
	_, err := CauchyCDF(1, 0, -1)
	requireErr(t, "γ < 0", err, ErrOutOfRange)
	_, err = CauchyCDF(inf, 0, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func BenchmarkCauchyCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CauchyCDF(1.5, 0, 1)
	}
}
