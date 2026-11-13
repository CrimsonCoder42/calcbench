package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.4.2 eq. (4.4.9), hand-derived: the origin is √2/2 from x + y = 1.
func TestPointLineDistance_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PointLineDistance((0,0); (0,1),(1,0))", func() (float64, error) {
			return PointLineDistance(Point{0, 0}, Point{0, 1}, Point{1, 0})
		}, math.Sqrt2 / 2, tolFew},
	})
}

func TestPointLineDistance_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"point on the line", func() (float64, error) { return PointLineDistance(Point{2, 2}, Point{0, 0}, Point{1, 1}) }, 0, 0},
		{"beyond the segment (infinite line)", func() (float64, error) {
			return PointLineDistance(Point{10, 3}, Point{0, 0}, Point{1, 0})
		}, 3, 0},
		{"vertical line", func() (float64, error) { return PointLineDistance(Point{-4, 7}, Point{1, 0}, Point{1, 5}) }, 5, 0},
		{"far from the origin", func() (float64, error) {
			return PointLineDistance(Point{1e12, 1e12 + 2}, Point{1e12, 1e12}, Point{1e12 + 1, 1e12})
		}, 2, 0},
		{"differences overflow, distance finite", func() (float64, error) {
			return PointLineDistance(Point{0, 1}, Point{-1e308, 0}, Point{1e308, 0})
		}, 1, 0},
	})
}

func TestPointLineDistance_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN p", func() error { return e1(PointLineDistance(Point{nan, 0}, Point{0, 0}, Point{1, 0})) }, ErrNotFinite},
		{"Inf a", func() error { return e1(PointLineDistance(Point{}, Point{inf, 0}, Point{1, 0})) }, ErrNotFinite},
		{"Inf b", func() error { return e1(PointLineDistance(Point{}, Point{0, 0}, Point{1, -inf})) }, ErrNotFinite},
		{"a == b", func() error { return e1(PointLineDistance(Point{3, 3}, Point{1, 1}, Point{1, 1})) }, ErrDegenerate},
		{"overflow", func() error {
			return e1(PointLineDistance(Point{0, 1e308}, Point{0, -1e308}, Point{1, -1e308}))
		}, ErrOverflow},
	})
}

// Oracle: |u×w|/|u| computed directly, on points well away from the line (no cancellation).
func TestPointLineDistance_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a := Point{rng.NormFloat64() * 10, rng.NormFloat64() * 10}
		b := Point{a.X + rng.NormFloat64()*10, a.Y + rng.NormFloat64()*10}
		ux, uy := b.X-a.X, b.Y-a.Y
		n := math.Hypot(ux, uy)
		off := 1 + rng.Float64()*10
		along := rng.NormFloat64() * 10
		p := Point{a.X + along*ux/n - off*uy/n, a.Y + along*uy/n + off*ux/n}
		want := math.Abs(ux*(p.Y-a.Y)-uy*(p.X-a.X)) / n
		if got, err := PointLineDistance(p, a, b); err != nil || !relNear(got, want, 1e-12) {
			t.Fatalf("PointLineDistance(%v, %v, %v) = (%v, %v), want %v", p, a, b, got, err, want)
		}
	}
}

func TestPointLineDistance_allocs(t *testing.T) {
	assertNoAllocs(t, "PointLineDistance", func() { sink, _ = PointLineDistance(Point{0, 0}, Point{0, 1}, Point{1, 0}) })
}

func BenchmarkPointLineDistance(b *testing.B) {
	b.ReportAllocs()
	p, a, q := Point{0, 0}, Point{0, 1}, Point{1, 0}
	for b.Loop() {
		sink, _ = PointLineDistance(p, a, q)
	}
}
