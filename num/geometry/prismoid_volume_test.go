package geometry

import (
	"math"
	"testing"
)

// Worked example: EB1911 "Mensuration" §58 (⅙H(S₀ + 4S₁ + S₂)), hand-derived for the square frustum with
// end sides 4 and 2 and height 3: (3/6)(16 + 4·9 + 4) = 28, equal to PyramidFrustumVolume(16, 4, 3).
func TestPrismoidVolume_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PrismoidVolume(16, 9, 4, 3)", func() (float64, error) { return PrismoidVolume(16, 9, 4, 3) }, 28, tolFew},
	})
}

// The solids EB1911 §58 names: cones, spheres, ellipsoids and prismoids, whose cross-section area is a
// quadratic function of the distance, come out exactly.
func TestPrismoidVolume_cases(t *testing.T) {
	sphere, _ := SphereVolume(2)
	cone, _ := ConeVolume(3, 4)
	ellipsoid, _ := EllipsoidVolume(3, 2, 1)
	runValCases(t, []valCase{
		{"sphere: 0, πr², 0 over 2r", func() (float64, error) { return PrismoidVolume(0, 4*math.Pi, 0, 4) }, sphere, tolFew},
		{"cone: πr², π(r/2)², 0 over h", func() (float64, error) { return PrismoidVolume(9*math.Pi, 9*math.Pi/4, 0, 4) }, cone, tolFew},
		{"ellipsoid: 0, πab, 0 over 2c", func() (float64, error) { return PrismoidVolume(0, 6*math.Pi, 0, 2) }, ellipsoid, tolFew},
		{"wedge with a zero end", func() (float64, error) { return PrismoidVolume(6, 3, 0, 2) }, 6, tolFew},
		{"all areas zero", func() (float64, error) { return PrismoidVolume(0, 0, 0, 2) }, 0, 0},
	})
}

func TestPrismoidVolume_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN area", func() error { return e1(PrismoidVolume(nan, 1, 1, 1)) }, ErrNotFinite},
		{"Inf height", func() error { return e1(PrismoidVolume(1, 1, 1, inf)) }, ErrNotFinite},
		{"negative area", func() error { return e1(PrismoidVolume(1, -1, 1, 1)) }, ErrNegative},
		{"zero height", func() error { return e1(PrismoidVolume(1, 1, 1, 0)) }, ErrNonPositive},
		{"overflow", func() error { return e1(PrismoidVolume(1e200, 1e200, 1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: random square and circular frustums, against PyramidFrustumVolume and ConeFrustumVolume.
func TestPrismoidVolume_vsFrustums(t *testing.T) {
	rng := newRNG()
	for range 400 {
		s1, s2, h := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		mid := (s1 + s2) / 2
		want, _ := PyramidFrustumVolume(s1*s1, s2*s2, h)
		if got, err := PrismoidVolume(s1*s1, mid*mid, s2*s2, h); err != nil || !relNear(got, want, 1e-14) {
			t.Fatalf("square frustum %v, %v, %v: (%v, %v), want %v", s1, s2, h, got, err, want)
		}
		cone, _ := ConeFrustumVolume(s1, s2, h)
		if got, err := PrismoidVolume(math.Pi*s1*s1, math.Pi*mid*mid, math.Pi*s2*s2, h); err != nil || !relNear(got, cone, 1e-14) {
			t.Fatalf("circular frustum %v, %v, %v: (%v, %v), want %v", s1, s2, h, got, err, cone)
		}
	}
}

func TestPrismoidVolume_allocs(t *testing.T) {
	assertNoAllocs(t, "PrismoidVolume", func() { sink, _ = PrismoidVolume(16, 9, 4, 3) })
}

func BenchmarkPrismoidVolume(b *testing.B) {
	b.ReportAllocs()
	a1, am, a2, h := 16.0, 9.0, 4.0, 3.0
	for b.Loop() {
		sink, _ = PrismoidVolume(a1, am, a2, h)
	}
}
