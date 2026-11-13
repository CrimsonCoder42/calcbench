package geometry

import (
	"math"
	"testing"
)

// Worked example: CRC 31st ed. §4.5 (the sides AᵢAᵢ₊₁) with §4.4.2 eq. (4.4.5), hand-derived for the
// pentagon: 4 + 3 + 2√2 + 2√2 + 3 = 10 + 4√2.
func TestPolygonPerimeter_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PolygonPerimeter(pentagon)", func() (float64, error) { return PolygonPerimeter(pentagon) }, 10 + 4*math.Sqrt2, tolFew},
	})
}

func TestPolygonPerimeter_cases(t *testing.T) {
	far := []Point{{1e15, 1e15}, {1e15 + 1, 1e15}, {1e15 + 1, 1e15 + 1}, {1e15, 1e15 + 1}}
	runValCases(t, []valCase{
		{"two vertices: there and back", func() (float64, error) { return PolygonPerimeter([]Point{{0, 0}, {3, 4}}) }, 10, 0},
		{"repeated vertex adds nothing", func() (float64, error) {
			return PolygonPerimeter([]Point{{0, 0}, {3, 0}, {3, 0}, {3, 4}})
		}, 12, 0},
		{"unit square far from the origin", func() (float64, error) { return PolygonPerimeter(far) }, 4, 0},
		{"huge coordinates, perimeter near the float64 limit", func() (float64, error) {
			return PolygonPerimeter([]Point{{-3e307, 0}, {3e307, 0}, {3e307, 1e307}, {-3e307, 1e307}})
		}, 1.4e308, tolFew},
	})
}

func TestPolygonPerimeter_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"nil", func() error { return e1(PolygonPerimeter(nil)) }, ErrDegenerate},
		{"one vertex", func() error { return e1(PolygonPerimeter([]Point{{1, 1}})) }, ErrDegenerate},
		{"all coincide", func() error { return e1(PolygonPerimeter([]Point{{2, 2}, {2, 2}, {2, 2}})) }, ErrDegenerate},
		{"NaN vertex", func() error { return e1(PolygonPerimeter([]Point{{0, 0}, {nan, 1}})) }, ErrNotFinite},
		{"Inf vertex", func() error { return e1(PolygonPerimeter([]Point{{0, -inf}, {0, 1}})) }, ErrNotFinite},
		{"overflow", func() error { return e1(PolygonPerimeter([]Point{{-1e308, 0}, {1e308, 0}})) }, ErrOverflow},
	})
}

// Oracle: the plain sum of math.Hypot over the sides, on random simple polygons.
func TestPolygonPerimeter_oracle(t *testing.T) {
	rng := newRNG()
	for trial := range 300 {
		vs := starPolygon(rng, 2+rng.IntN(200), Point{rng.NormFloat64() * 10, rng.NormFloat64() * 10}, 1+rng.Float64()*100)
		want := 0.0
		for i, v := range vs {
			w := vs[(i+1)%len(vs)]
			want += math.Hypot(w.X-v.X, w.Y-v.Y)
		}
		if got, err := PolygonPerimeter(vs); err != nil || !relNear(got, want, 1e-13) {
			t.Fatalf("trial %d (n = %d): PolygonPerimeter = (%v, %v), want %v", trial, len(vs), got, err, want)
		}
	}
}

// Oracle: the regular n-gon of circumradius 1 has perimeter n·RegularPolygonSideInscribed(n, 1).
func TestPolygonPerimeter_regular(t *testing.T) {
	for _, n := range []int{3, 4, 7, 100, 10000} {
		side, _ := RegularPolygonSideInscribed(n, 1)
		if got, err := PolygonPerimeter(regularPolygon(n)); err != nil || !relNear(got, float64(n)*side, 1e-12) {
			t.Errorf("n %d: PolygonPerimeter = (%v, %v), want %v", n, got, err, float64(n)*side)
		}
	}
}

func TestPolygonPerimeter_allocs(t *testing.T) {
	vs := regularPolygon(64)
	assertNoAllocs(t, "PolygonPerimeter", func() { sink, _ = PolygonPerimeter(vs) })
}

func BenchmarkPolygonPerimeter(b *testing.B) {
	for _, size := range polygonSizes {
		vs := regularPolygon(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sink, _ = PolygonPerimeter(vs)
			}
		})
	}
}
