package geometry

import (
	"math"
	"testing"
)

// Worked examples, hand-derived: CRC 31st ed. §4.16 eq. (4.16.1) (V = ⅓Ah): a square base of side 6 and
// height 4 → 48; §4.5.3 and §4.5.1 for the regular pyramid: base 36, slant height √(4² + 3²) = 5, four
// faces of ½·6·5 → 96; eq. (4.16.2) (V = ⅓h(A + A′ + √(AA′))): 16, 4, 3 → 28.
func TestPyramid_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PyramidVolume(36, 4)", func() (float64, error) { return PyramidVolume(36, 4) }, 48, tolFew},
		{"RegularPyramidSurfaceArea(4, 6, 4)", func() (float64, error) { return RegularPyramidSurfaceArea(4, 6, 4) }, 96, tolFew},
		{"PyramidFrustumVolume(16, 4, 3)", func() (float64, error) { return PyramidFrustumVolume(16, 4, 3) }, 28, tolFew},
	})
}

func TestPyramid_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"tetrahedron as a triangular pyramid", func() (float64, error) {
			base, _ := RegularPolygonAreaFromSide(3, 1)
			return PyramidVolume(base, math.Sqrt(2.0/3))
		}, math.Sqrt2 / 12, tolSome},
		{"frustum with equal bases is a prism", func() (float64, error) { return PyramidFrustumVolume(5, 5, 2) }, 10, tolFew},
		{"frustum, areas whose product overflows", func() (float64, error) { return PyramidFrustumVolume(1e200, 1e200, 1e-100) }, 1e100, tolSome},
	})
}

func TestPyramid_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(PyramidVolume(nan, 1)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(PyramidVolume(0, 1)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(PyramidVolume(1e200, 1e200)) }, ErrOverflow},
		{"surface too few sides", func() error { return e1(RegularPyramidSurfaceArea(2, 1, 1)) }, ErrTooFewSides},
		{"surface Inf", func() error { return e1(RegularPyramidSurfaceArea(4, inf, 1)) }, ErrNotFinite},
		{"surface zero height", func() error { return e1(RegularPyramidSurfaceArea(4, 1, 0)) }, ErrNonPositive},
		{"surface overflow", func() error { return e1(RegularPyramidSurfaceArea(4, 1e200, 1)) }, ErrOverflow},
		{"frustum NaN", func() error { return e1(PyramidFrustumVolume(1, nan, 1)) }, ErrNotFinite},
		{"frustum zero top", func() error { return e1(PyramidFrustumVolume(1, 0, 1)) }, ErrNonPositive},
		{"frustum overflow", func() error { return e1(PyramidFrustumVolume(1e200, 1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracles: the regular pyramid's area from the CRC's cot forms (base ¼ka² cot(π/k), apothem ½a cot(π/k))
// computed directly; and a cut-off cone (a frustum) as the large pyramid less the small one above it.
func TestPyramid_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		n := 3 + rng.IntN(20)
		s, h := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		cot := 1 / math.Tan(math.Pi/float64(n))
		apothem := s * cot / 2
		want := float64(n)*s*s*cot/4 + float64(n)*s*math.Sqrt(h*h+apothem*apothem)/2
		if got, err := RegularPyramidSurfaceArea(n, s, h); err != nil || !relNear(got, want, 1e-14) {
			t.Fatalf("RegularPyramidSurfaceArea(%d, %v, %v) = (%v, %v), want %v", n, s, h, got, err, want)
		}

		// A pyramid of base A and height H, cut at height H − h' from the base, leaves a top of
		// A·(h'/H)² over a frustum of height H − h'.
		bigA, bigH := 1+rng.Float64()*10, 1+rng.Float64()*10
		top := bigH * (0.1 + 0.8*rng.Float64())
		small := bigA * (top / bigH) * (top / bigH)
		vBig, _ := PyramidVolume(bigA, bigH)
		vSmall, _ := PyramidVolume(small, top)
		if got, err := PyramidFrustumVolume(bigA, small, bigH-top); err != nil || !relNear(got, vBig-vSmall, 1e-13) {
			t.Fatalf("PyramidFrustumVolume(%v, %v, %v) = (%v, %v), want %v", bigA, small, bigH-top, got, err, vBig-vSmall)
		}
	}
}

func TestPyramid_allocs(t *testing.T) {
	assertNoAllocs(t, "PyramidVolume", func() { sink, _ = PyramidVolume(36, 4) })
	assertNoAllocs(t, "RegularPyramidSurfaceArea", func() { sink, _ = RegularPyramidSurfaceArea(4, 6, 4) })
	assertNoAllocs(t, "PyramidFrustumVolume", func() { sink, _ = PyramidFrustumVolume(16, 4, 3) })
}

func BenchmarkPyramidVolume(b *testing.B) {
	b.ReportAllocs()
	base, h := 36.0, 4.0
	for b.Loop() {
		sink, _ = PyramidVolume(base, h)
	}
}

func BenchmarkRegularPyramidSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	n, s, h := 4, 6.0, 4.0
	for b.Loop() {
		sink, _ = RegularPyramidSurfaceArea(n, s, h)
	}
}

func BenchmarkPyramidFrustumVolume(b *testing.B) {
	b.ReportAllocs()
	a1, a2, h := 16.0, 4.0, 3.0
	for b.Loop() {
		sink, _ = PyramidFrustumVolume(a1, a2, h)
	}
}
