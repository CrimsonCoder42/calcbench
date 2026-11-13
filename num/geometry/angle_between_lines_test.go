package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.4.3 eq. (4.4.10), hand-derived: y = 0 and the line of slope 1/√3 meet
// at π/6.
func TestAngleBetweenLinesRad_example(t *testing.T) {
	runValCases(t, []valCase{
		{"AngleBetweenLinesRad(y = 0, slope 1/√3)", func() (float64, error) {
			return AngleBetweenLinesRad(Point{0, 0}, Point{1, 0}, Point{0, 0}, Point{math.Sqrt(3), 1})
		}, math.Pi / 6, tolFew},
	})
}

func TestAngleBetweenLinesRad_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"parallel", func() (float64, error) {
			return AngleBetweenLinesRad(Point{0, 0}, Point{1, 1}, Point{0, 5}, Point{3, 8})
		}, 0, 0},
		{"perpendicular, one vertical", func() (float64, error) {
			return AngleBetweenLinesRad(Point{0, 0}, Point{1, 0}, Point{2, -1}, Point{2, 7})
		}, math.Pi / 2, tolFew},
		{"obtuse pair gives the acute angle", func() (float64, error) {
			return AngleBetweenLinesRad(Point{0, 0}, Point{1, 0}, Point{0, 0}, Point{-1, 1})
		}, math.Pi / 4, tolFew},
		{"direction reversal changes nothing", func() (float64, error) {
			return AngleBetweenLinesRad(Point{1, 0}, Point{0, 0}, Point{-1, 1}, Point{0, 0})
		}, math.Pi / 4, tolFew},
		{"differences overflow", func() (float64, error) {
			return AngleBetweenLinesRad(Point{-1e308, 0}, Point{1e308, 0}, Point{0, -1e308}, Point{0, 1e308})
		}, math.Pi / 2, tolFew},
	})
}

func TestAngleBetweenLinesRad_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN a", func() error { return e1(AngleBetweenLinesRad(Point{nan, 0}, Point{1, 0}, Point{}, Point{0, 1})) }, ErrNotFinite},
		{"Inf d", func() error { return e1(AngleBetweenLinesRad(Point{}, Point{1, 0}, Point{}, Point{0, inf})) }, ErrNotFinite},
		{"a == b", func() error { return e1(AngleBetweenLinesRad(Point{1, 1}, Point{1, 1}, Point{}, Point{0, 1})) }, ErrDegenerate},
		{"c == d", func() error { return e1(AngleBetweenLinesRad(Point{}, Point{1, 0}, Point{2, 2}, Point{2, 2})) }, ErrDegenerate},
	})
}

// Oracle: lines with directions φ₁ and φ₂ meet at |φ₁ − φ₂| reduced to [0, π/2].
func TestAngleBetweenLinesRad_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		phi1, phi2 := rng.Float64()*2*math.Pi, rng.Float64()*2*math.Pi
		r1, r2 := 0.5+rng.Float64()*10, 0.5+rng.Float64()*10
		o := Point{rng.NormFloat64(), rng.NormFloat64()}
		s1, c1 := math.Sincos(phi1)
		s2, c2 := math.Sincos(phi2)
		want := math.Mod(math.Abs(phi1-phi2), math.Pi)
		if want > math.Pi/2 {
			want = math.Pi - want
		}
		got, err := AngleBetweenLinesRad(o, Point{o.X + r1*c1, o.Y + r1*s1}, o, Point{o.X + r2*c2, o.Y + r2*s2})
		if err != nil || math.Abs(got-want) > 1e-13 {
			t.Fatalf("φ₁ %v, φ₂ %v: got (%v, %v), want %v", phi1, phi2, got, err, want)
		}
	}
}

func TestAngleBetweenLinesRad_allocs(t *testing.T) {
	assertNoAllocs(t, "AngleBetweenLinesRad", func() {
		sink, _ = AngleBetweenLinesRad(Point{0, 0}, Point{1, 0}, Point{0, 0}, Point{1, 2})
	})
}

func BenchmarkAngleBetweenLinesRad(b *testing.B) {
	b.ReportAllocs()
	p, q, r, s := Point{0, 0}, Point{1, 0}, Point{0, 0}, Point{1, 2}
	for b.Loop() {
		sink, _ = AngleBetweenLinesRad(p, q, r, s)
	}
}
