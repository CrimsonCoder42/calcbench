package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.1.4.1 (x = r cos θ, y = r sin θ), hand-derived: r 2, θ π/6 → (√3, 1).
func TestPolarToCartesian_example(t *testing.T) {
	got, err := PolarToCartesian(2, math.Pi/6)
	if err != nil || !relNear(got.X, math.Sqrt(3), tolFew) || !relNear(got.Y, 1, tolFew) {
		t.Errorf("PolarToCartesian(2, π/6) = (%v, %v), want (√3, 1)", got, err)
	}
}

func TestPolarToCartesian_cases(t *testing.T) {
	tests := []struct {
		name     string
		r, theta float64
		want     Point
	}{
		{"zero radius", 0, 1.234, Point{0, 0}},
		{"on the x-axis", 3, 0, Point{3, 0}},
		{"largest radius cannot overflow", math.MaxFloat64, 0, Point{math.MaxFloat64, 0}},
	}
	for _, tt := range tests {
		if got, err := PolarToCartesian(tt.r, tt.theta); err != nil || got != tt.want {
			t.Errorf("%s: PolarToCartesian(%v, %v) = (%v, %v), want %v", tt.name, tt.r, tt.theta, got, err, tt.want)
		}
	}
	got, err := PolarToCartesian(1, math.Pi)
	if err != nil || got.X != -1 || math.Abs(got.Y) > 1e-15 {
		t.Errorf("PolarToCartesian(1, π) = (%v, %v), want (−1, ≈0)", got, err)
	}
}

func TestPolarToCartesian_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN r", func() error { return ep(PolarToCartesian(nan, 0)) }, ErrNotFinite},
		{"Inf θ", func() error { return ep(PolarToCartesian(1, inf)) }, ErrNotFinite},
		{"negative r", func() error { return ep(PolarToCartesian(-1, 0)) }, ErrNegative},
	})
}

// Worked example: CRC 31st ed. §4.1.4.1, hand-derived: (−1, 0) → (1, π).
func TestCartesianToPolar_example(t *testing.T) {
	r, theta, err := CartesianToPolar(Point{-1, 0})
	if err != nil || r != 1 || theta != math.Pi {
		t.Errorf("CartesianToPolar((−1, 0)) = (%v, %v, %v), want (1, π)", r, theta, err)
	}
}

func TestCartesianToPolar_cases(t *testing.T) {
	tests := []struct {
		name     string
		p        Point
		r, theta float64
	}{
		{"origin", Point{0, 0}, 0, 0},
		{"-0 ordinate on the negative axis gives +π", Point{-2, math.Copysign(0, -1)}, 2, math.Pi},
		{"3-4-5", Point{3, 4}, 5, math.Atan2(4, 3)},
		{"no overflow of the squares", Point{3e300, 4e300}, 5e300, math.Atan2(4, 3)},
	}
	for _, tt := range tests {
		r, theta, err := CartesianToPolar(tt.p)
		if err != nil || !relNear(r, tt.r, tolFew) || !relNear(theta, tt.theta, tolFew) {
			t.Errorf("%s: CartesianToPolar(%v) = (%v, %v, %v), want (%v, %v)", tt.name, tt.p, r, theta, err, tt.r, tt.theta)
		}
	}
}

func TestCartesianToPolar_errors(t *testing.T) {
	polarErr := func(_, _ float64, err error) error { return err }
	runErrCases(t, []errCase{
		{"NaN", func() error { return polarErr(CartesianToPolar(Point{nan, 0})) }, ErrNotFinite},
		{"Inf", func() error { return polarErr(CartesianToPolar(Point{0, -inf})) }, ErrNotFinite},
		{"overflow", func() error {
			return polarErr(CartesianToPolar(Point{math.MaxFloat64, math.MaxFloat64}))
		}, ErrOverflow},
	})
}

// Oracle: the two conversions are inverse to each other.
func TestPolar_roundTrip(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := rng.Float64() * 100
		theta := (rng.Float64()*2 - 1) * math.Pi
		p, err := PolarToCartesian(r, theta)
		if err != nil {
			t.Fatal(err)
		}
		r2, theta2, err := CartesianToPolar(p)
		if err != nil || !relNear(r2, r, tolSome) || math.Abs(theta2-theta) > 1e-14 {
			t.Fatalf("round trip (%v, %v) → %v → (%v, %v, %v)", r, theta, p, r2, theta2, err)
		}
	}
}

func TestPolar_allocs(t *testing.T) {
	assertNoAllocs(t, "PolarToCartesian", func() { sinkPoint, _ = PolarToCartesian(2, 0.5) })
	assertNoAllocs(t, "CartesianToPolar", func() { sink, _, _ = CartesianToPolar(Point{3, 4}) })
}

func BenchmarkPolarToCartesian(b *testing.B) {
	b.ReportAllocs()
	r, theta := 2.0, 0.5
	for b.Loop() {
		sinkPoint, _ = PolarToCartesian(r, theta)
	}
}

func BenchmarkCartesianToPolar(b *testing.B) {
	b.ReportAllocs()
	p := Point{3, 4}
	for b.Loop() {
		sink, _, _ = CartesianToPolar(p)
	}
}
