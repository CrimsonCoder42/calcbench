package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.4.2 (midpoint of P₀P₁), hand-derived.
func TestMidpoint_example(t *testing.T) {
	if got, err := Midpoint(Point{0, 0}, Point{4, -2}); err != nil || got != (Point{2, -1}) {
		t.Errorf("Midpoint((0,0),(4,-2)) = (%v, %v), want (2, -1)", got, err)
	}
}

func TestMidpoint_cases(t *testing.T) {
	tests := []struct {
		name string
		p, q Point
		want Point
	}{
		{"same point", Point{3, 3}, Point{3, 3}, Point{3, 3}},
		{"symmetric about the origin", Point{-5, 7}, Point{5, -7}, Point{0, 0}},
		{"no overflow at the extremes", Point{math.MaxFloat64, math.MaxFloat64}, Point{math.MaxFloat64, -math.MaxFloat64}, Point{math.MaxFloat64, 0}},
	}
	for _, tt := range tests {
		if got, err := Midpoint(tt.p, tt.q); err != nil || got != tt.want {
			t.Errorf("%s: Midpoint(%v, %v) = (%v, %v), want %v", tt.name, tt.p, tt.q, got, err, tt.want)
		}
	}
}

func TestMidpoint_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN p", func() error { return ep(Midpoint(Point{0, nan}, Point{})) }, ErrNotFinite},
		{"Inf q", func() error { return ep(Midpoint(Point{}, Point{inf, 0})) }, ErrNotFinite},
	})
}

// Oracle: x/2 + y/2 and (x + y)/2 round the same exact value once when nothing overflows or is subnormal.
func TestMidpoint_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		// float64(…) forces each coordinate to be rounded: without it the compiler may fuse the
		// multiply into the oracle's addition (an FMA, allowed by the Go spec) and skip that rounding.
		p := Point{float64(rng.NormFloat64() * 1e6), float64(rng.NormFloat64() * 1e6)}
		q := Point{float64(rng.NormFloat64() * 1e6), float64(rng.NormFloat64() * 1e6)}
		want := Point{(p.X + q.X) / 2, (p.Y + q.Y) / 2}
		if got, err := Midpoint(p, q); err != nil || got != want {
			t.Fatalf("Midpoint(%v, %v) = (%v, %v), want %v", p, q, got, err, want)
		}
	}
}

func TestMidpoint_allocs(t *testing.T) {
	assertNoAllocs(t, "Midpoint", func() { sinkPoint, _ = Midpoint(Point{1, 2}, Point{4, 6}) })
}

func BenchmarkMidpoint(b *testing.B) {
	b.ReportAllocs()
	p, q := Point{1, 2}, Point{4, 6}
	for b.Loop() {
		sinkPoint, _ = Midpoint(p, q)
	}
}
