package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.14.1 tables, p. 360: V/a³ = 0.471405 and S/a² = 3.46410 for the octahedron, printed
// to 6 or 7 significant digits (hence 1e-5); and the table's exact forms, V/a³ = math.Sqrt2 / 3 and S/a² = 2 * math.Sqrt(3).
func TestOctahedron_example(t *testing.T) {
	runValCases(t, []valCase{
		{"OctahedronVolume(1), table", func() (float64, error) { return OctahedronVolume(1) }, 0.471405, 1e-5},
		{"OctahedronSurfaceArea(1), table", func() (float64, error) { return OctahedronSurfaceArea(1) }, 3.46410, 1e-5},
		{"OctahedronVolume(1), exact form", func() (float64, error) { return OctahedronVolume(1) }, math.Sqrt2 / 3, tolFew},
		{"OctahedronSurfaceArea(1), exact form", func() (float64, error) { return OctahedronSurfaceArea(1) }, 2 * math.Sqrt(3), tolFew},
	})
}

func TestOctahedron_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume scales as a³", func() (float64, error) { return OctahedronVolume(2) }, 8 * (math.Sqrt2 / 3), tolFew},
		{"area scales as a²", func() (float64, error) { return OctahedronSurfaceArea(2) }, 4 * (2 * math.Sqrt(3)), tolFew},
		{"tiny edge", func() (float64, error) { return OctahedronVolume(1e-100) }, (math.Sqrt2 / 3) * 1e-300, tolFew},
	})
}

func TestOctahedron_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(OctahedronVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(OctahedronVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(OctahedronVolume(1e104)) }, ErrOverflow},
		{"area Inf", func() error { return e1(OctahedronSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(OctahedronSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(OctahedronSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracles: the surface is 8 regular 3-gons (RegularPolygonAreaFromSide), and the volume is ⅓·r·S with
// the inradius r/a = math.Sqrt(6) / 6 (CRC 31st ed. eq. (4.14.7) and the r/a table).
func TestOctahedron_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.01 + rng.Float64()*100
		face, _ := RegularPolygonAreaFromSide(3, a)
		s, err := OctahedronSurfaceArea(a)
		if err != nil || !relNear(s, 8*face, tolSome) {
			t.Fatalf("OctahedronSurfaceArea(%v) = (%v, %v), faces %v", a, s, err, 8*face)
		}
		want := s * (math.Sqrt(6) / 6) * a / 3
		if got, err := OctahedronVolume(a); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("OctahedronVolume(%v) = (%v, %v), ⅓rS %v", a, got, err, want)
		}
	}
}

func TestOctahedron_allocs(t *testing.T) {
	assertNoAllocs(t, "OctahedronVolume", func() { sink, _ = OctahedronVolume(2) })
	assertNoAllocs(t, "OctahedronSurfaceArea", func() { sink, _ = OctahedronSurfaceArea(2) })
}

func BenchmarkOctahedronVolume(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = OctahedronVolume(a)
	}
}

func BenchmarkOctahedronSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = OctahedronSurfaceArea(a)
	}
}
