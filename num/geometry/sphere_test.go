package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.18.1 item 3, p. 367 (area 4πr², volume (4/3)πr³), hand-derived for
// r 3: both 36π.
func TestSphere_example(t *testing.T) {
	runValCases(t, []valCase{
		{"SphereVolume(3)", func() (float64, error) { return SphereVolume(3) }, 36 * math.Pi, tolFew},
		{"SphereSurfaceArea(3)", func() (float64, error) { return SphereSurfaceArea(3) }, 36 * math.Pi, tolFew},
	})
}

func TestSphere_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"unit sphere volume", func() (float64, error) { return SphereVolume(1) }, 4 * math.Pi / 3, tolFew},
		{"tiny radius", func() (float64, error) { return SphereSurfaceArea(1e-150) }, 4 * math.Pi * 1e-300, tolFew},
	})
}

func TestSphere_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(SphereVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(SphereVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(SphereVolume(1e103)) }, ErrOverflow},
		{"area Inf", func() error { return e1(SphereSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(SphereSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(SphereSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracle: the sphere is the ellipsoid with three equal semi-axes, and the full cap (h = 2r).
func TestSphere_vsEllipsoidAndCap(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.01 + rng.Float64()*100
		v, _ := SphereVolume(r)
		s, _ := SphereSurfaceArea(r)
		ev, _ := EllipsoidVolume(r, r, r)
		es, _ := EllipsoidSurfaceArea(r, r, r)
		cv, _ := SphericalCapVolume(r, 2*r)
		cs, _ := SphericalCapArea(r, 2*r)
		if !relNear(v, ev, tolSome) || !relNear(s, es, tolSome) || !relNear(v, cv, tolSome) || !relNear(s, cs, tolSome) {
			t.Fatalf("r %v: volume %v, ellipsoid %v, cap %v; area %v, ellipsoid %v, cap %v", r, v, ev, cv, s, es, cs)
		}
	}
}

func TestSphere_allocs(t *testing.T) {
	assertNoAllocs(t, "SphereVolume", func() { sink, _ = SphereVolume(3) })
	assertNoAllocs(t, "SphereSurfaceArea", func() { sink, _ = SphereSurfaceArea(3) })
}

func BenchmarkSphereVolume(b *testing.B) {
	b.ReportAllocs()
	r := 3.0
	for b.Loop() {
		sink, _ = SphereVolume(r)
	}
}

func BenchmarkSphereSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	r := 3.0
	for b.Loop() {
		sink, _ = SphereSurfaceArea(r)
	}
}
