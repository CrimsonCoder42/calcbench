package prob

import "testing"

// JKB vol. 1 ch. 16, inverted: Q(0.75) = 1, Q(0.25) = −1, Q(½) = x₀; round trip CDF(Q(u)) = u.
func TestCauchyQuantile_sourceExample(t *testing.T) {
	cases := []struct{ u, want float64 }{{0.75, 1}, {0.25, -1}, {0.5, 0}}
	for _, c := range cases {
		if got, err := CauchyQuantile(c.u, 0, 1); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("CauchyQuantile(%v) = (%v, %v), want %v", c.u, got, err, c.want)
		}
	}
	for _, u := range []float64{1e-12, 0.1, 0.6, 1 - 1e-9} {
		x, _ := CauchyQuantile(u, 2, 3)
		if c, _ := CauchyCDF(x, 2, 3); !within(c, u, 1e-9) {
			t.Errorf("u %v: CDF(Q) = %v", u, c)
		}
	}
	for _, u := range []float64{0, 1} {
		_, err := CauchyQuantile(u, 0, 1)
		requireErr(t, "infinite", err, ErrOutOfRange)
	}
	_, err := CauchyQuantile(2, 0, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = CauchyQuantile(0.5, inf, 1)
	requireErr(t, "x₀ Inf", err, ErrNotFinite)
	_, err = CauchyQuantile(1e-300, 0, 1e10)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkCauchyQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CauchyQuantile(0.9, 0, 1)
	}
}
