package geometry

import (
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// pentagon is the catalog's irregular pentagon: a 4 × 3 rectangle with a triangular roof of apex (2, 5).
var pentagon = []Point{{0, 0}, {4, 0}, {4, 3}, {2, 5}, {0, 3}}

// Worked example: CRC 31st ed. §4.5 eq. (4.5.1), hand-derived: rectangle 12 + roof triangle ½·4·2 = 16.
func TestPolygonArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PolygonArea(pentagon)", func() (float64, error) { return PolygonArea(pentagon) }, 16, tolFew},
	})
}

func TestPolygonArea_cases(t *testing.T) {
	reversed := slices.Clone(pentagon)
	slices.Reverse(reversed)
	far := []Point{{1e15, 1e15}, {1e15 + 1, 1e15}, {1e15 + 1, 1e15 + 1}, {1e15, 1e15 + 1}}
	runValCases(t, []valCase{
		{"triangle (CRC eq. (4.5.2))", func() (float64, error) { return PolygonArea([]Point{{0, 0}, {4, 0}, {0, 3}}) }, 6, tolFew},
		{"clockwise order", func() (float64, error) { return PolygonArea(reversed) }, 16, tolFew},
		{"unit square far from the origin", func() (float64, error) { return PolygonArea(far) }, 1, 0},
		{"non-convex (an L shape)", func() (float64, error) {
			return PolygonArea([]Point{{0, 0}, {3, 0}, {3, 1}, {1, 1}, {1, 3}, {0, 3}})
		}, 5, tolFew},
		{"huge coordinates, no overflow", func() (float64, error) {
			return PolygonArea([]Point{{0, 0}, {1e154, 0}, {1e154, 1e154}, {0, 1e154}})
		}, 1e308, tolFew},
	})
	if !slices.Equal(reversed, []Point{{0, 3}, {2, 5}, {4, 3}, {4, 0}, {0, 0}}) {
		t.Errorf("PolygonArea modified its input: %v", reversed)
	}
}

func TestPolygonArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"nil", func() error { return e1(PolygonArea(nil)) }, ErrDegenerate},
		{"two vertices", func() error { return e1(PolygonArea([]Point{{0, 0}, {1, 1}})) }, ErrDegenerate},
		{"collinear", func() error { return e1(PolygonArea([]Point{{0, 0}, {1, 1}, {2, 2}, {3, 3}})) }, ErrDegenerate},
		{"self-intersecting bow tie: the signed areas cancel", func() error {
			return e1(PolygonArea([]Point{{0, 0}, {2, 2}, {2, 0}, {0, 2}}))
		}, ErrDegenerate},
		{"NaN vertex", func() error { return e1(PolygonArea([]Point{{0, 0}, {1, nan}, {0, 1}})) }, ErrNotFinite},
		{"Inf vertex", func() error { return e1(PolygonArea([]Point{{0, 0}, {1, 0}, {inf, 1}})) }, ErrNotFinite},
		{"overflow", func() error {
			return e1(PolygonArea([]Point{{-1e200, -1e200}, {1e200, -1e200}, {1e200, 1e200}, {-1e200, 1e200}}))
		}, ErrOverflow},
	})
}

// starPolygon returns n vertices at increasing angles and random radii around center: a simple polygon.
func starPolygon(rng *rand.Rand, n int, center Point, scale float64) []Point {
	angles := make([]float64, n)
	for i := range angles {
		angles[i] = rng.Float64() * 2 * math.Pi
	}
	slices.Sort(angles)
	vs := make([]Point, n)
	for i, a := range angles {
		r := scale * (0.2 + rng.Float64())
		vs[i] = Point{center.X + r*math.Cos(a), center.Y + r*math.Sin(a)}
	}
	return vs
}

// shoelaceExact returns ½|Σ (xᵢyᵢ₊₁ − xᵢ₊₁yᵢ)| of the float64 vertices in exact big.Float arithmetic
// (4096 bits hold every product and sum here exactly), rounded once.
func shoelaceExact(vs []Point) float64 {
	const prec = 4096
	sum := new(big.Float).SetPrec(prec)
	term := new(big.Float).SetPrec(prec)
	for i, v := range vs {
		w := vs[(i+1)%len(vs)]
		term.Mul(new(big.Float).SetFloat64(v.X), new(big.Float).SetFloat64(w.Y))
		sum.Add(sum, term)
		term.Mul(new(big.Float).SetFloat64(w.X), new(big.Float).SetFloat64(v.Y))
		sum.Sub(sum, term)
	}
	f, _ := sum.Abs(sum).Float64()
	return f / 2
}

// Oracle: the exact shoelace sum on random simple polygons, near and far from the origin. Translation
// rounds each coordinate difference once and the compensated sum adds about one ulp, so 1e-13 leaves a
// wide margin for the cancellation between positive and negative cross terms.
func TestPolygonArea_oracle(t *testing.T) {
	rng := newRNG()
	for trial := range 300 {
		center := Point{rng.NormFloat64() * 10, rng.NormFloat64() * 10}
		if trial%3 == 0 {
			center = Point{1e6 + rng.Float64(), -1e6 + rng.Float64()}
		}
		vs := starPolygon(rng, 3+rng.IntN(200), center, 1+rng.Float64()*100)
		want := shoelaceExact(vs)
		if got, err := PolygonArea(vs); err != nil || !relNear(got, want, 1e-13) {
			t.Fatalf("trial %d (n = %d): PolygonArea = (%v, %v), exact %v", trial, len(vs), got, err, want)
		}
	}
}

// Oracle: the regular n-gon of circumradius 1 against RegularPolygonAreaFromCircumradius.
func TestPolygonArea_regular(t *testing.T) {
	for _, n := range []int{3, 4, 7, 100, 10000} {
		want, _ := RegularPolygonAreaFromCircumradius(n, 1)
		if got, err := PolygonArea(regularPolygon(n)); err != nil || !relNear(got, want, 1e-13) {
			t.Errorf("n %d: PolygonArea = (%v, %v), want %v", n, got, err, want)
		}
	}
}

// regularPolygon returns the vertices of the regular n-gon of circumradius 1 centred at (3, −2).
func regularPolygon(n int) []Point {
	vs := make([]Point, n)
	for i := range vs {
		s, c := math.Sincos(2 * math.Pi * float64(i) / float64(n))
		vs[i] = Point{3 + c, -2 + s}
	}
	return vs
}

func TestPolygonArea_allocs(t *testing.T) {
	vs := regularPolygon(64)
	assertNoAllocs(t, "PolygonArea", func() { sink, _ = PolygonArea(vs) })
}

var polygonSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

func BenchmarkPolygonArea(b *testing.B) {
	for _, size := range polygonSizes {
		vs := regularPolygon(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sink, _ = PolygonArea(vs)
			}
		})
	}
}
