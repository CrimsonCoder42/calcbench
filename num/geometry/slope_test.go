package geometry

import "testing"

// Worked example: CRC 31st ed. §4.4.1 item 6, hand-derived: (5 − 1)/(3 − 1) = 2.
func TestSlope_example(t *testing.T) {
	runValCases(t, []valCase{
		{"Slope((1,1),(3,5))", func() (float64, error) { return Slope(Point{1, 1}, Point{3, 5}) }, 2, 0},
	})
}

func TestSlope_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"horizontal", func() (float64, error) { return Slope(Point{-3, 4}, Point{8, 4}) }, 0, 0},
		{"negative", func() (float64, error) { return Slope(Point{0, 0}, Point{2, -6}) }, -3, 0},
		{"differences overflow, slope finite", func() (float64, error) {
			return Slope(Point{-1e308, -1e308}, Point{1e308, 1e308})
		}, 1, 0},
	})
}

func TestSlope_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(Slope(Point{nan, 0}, Point{1, 1})) }, ErrNotFinite},
		{"Inf", func() error { return e1(Slope(Point{0, 0}, Point{1, inf})) }, ErrNotFinite},
		{"same point", func() error { return e1(Slope(Point{2, 3}, Point{2, 3})) }, ErrDegenerate},
		{"vertical", func() error { return e1(Slope(Point{2, 3}, Point{2, 9})) }, ErrVertical},
		{"overflow", func() error { return e1(Slope(Point{0, -1e308}, Point{1e-300, 1e308})) }, ErrOverflow},
	})
}

// Property: the slope does not depend on the order of the points (q − p and p − q are exact negatives).
func TestSlope_symmetric(t *testing.T) {
	rng := newRNG()
	for range 400 {
		p := Point{rng.NormFloat64() * 100, rng.NormFloat64() * 100}
		q := Point{rng.NormFloat64() * 100, rng.NormFloat64() * 100}
		a, errA := Slope(p, q)
		b, errB := Slope(q, p)
		if errA != nil || errB != nil || a != b {
			t.Fatalf("Slope(%v, %v) = (%v, %v), reversed (%v, %v)", p, q, a, errA, b, errB)
		}
	}
}

func TestSlope_allocs(t *testing.T) {
	assertNoAllocs(t, "Slope", func() { sink, _ = Slope(Point{1, 1}, Point{3, 5}) })
}

func BenchmarkSlope(b *testing.B) {
	b.ReportAllocs()
	p, q := Point{1, 1}, Point{3, 5}
	for b.Loop() {
		sink, _ = Slope(p, q)
	}
}
