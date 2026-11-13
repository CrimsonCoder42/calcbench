package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.14.1 tables, p. 360: V/a³ = 0.117851 and S/a² = 1.73205 for the tetrahedron, printed
// to 6 or 7 significant digits (hence 1e-5); and the table's exact forms, V/a³ = math.Sqrt2 / 12 and S/a² = math.Sqrt(3).
func TestTetrahedron_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TetrahedronVolume(1), table", func() (float64, error) { return TetrahedronVolume(1) }, 0.117851, 1e-5},
		{"TetrahedronSurfaceArea(1), table", func() (float64, error) { return TetrahedronSurfaceArea(1) }, 1.73205, 1e-5},
		{"TetrahedronVolume(1), exact form", func() (float64, error) { return TetrahedronVolume(1) }, math.Sqrt2 / 12, tolFew},
		{"TetrahedronSurfaceArea(1), exact form", func() (float64, error) { return TetrahedronSurfaceArea(1) }, math.Sqrt(3), tolFew},
	})
}

func TestTetrahedron_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume scales as a³", func() (float64, error) { return TetrahedronVolume(2) }, 8 * (math.Sqrt2 / 12), tolFew},
		{"area scales as a²", func() (float64, error) { return TetrahedronSurfaceArea(2) }, 4 * (math.Sqrt(3)), tolFew},
		{"tiny edge", func() (float64, error) { return TetrahedronVolume(1e-100) }, (math.Sqrt2 / 12) * 1e-300, tolFew},
	})
}

func TestTetrahedron_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(TetrahedronVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(TetrahedronVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(TetrahedronVolume(1e104)) }, ErrOverflow},
		{"area Inf", func() error { return e1(TetrahedronSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(TetrahedronSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(TetrahedronSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracles: the surface is 4 regular 3-gons (RegularPolygonAreaFromSide), and the volume is ⅓·r·S with
// the inradius r/a = math.Sqrt(6) / 12 (CRC 31st ed. eq. (4.14.7) and the r/a table).
func TestTetrahedron_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.01 + rng.Float64()*100
		face, _ := RegularPolygonAreaFromSide(3, a)
		s, err := TetrahedronSurfaceArea(a)
		if err != nil || !relNear(s, 4*face, tolSome) {
			t.Fatalf("TetrahedronSurfaceArea(%v) = (%v, %v), faces %v", a, s, err, 4*face)
		}
		want := s * (math.Sqrt(6) / 12) * a / 3
		if got, err := TetrahedronVolume(a); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("TetrahedronVolume(%v) = (%v, %v), ⅓rS %v", a, got, err, want)
		}
	}
}

func TestTetrahedron_allocs(t *testing.T) {
	assertNoAllocs(t, "TetrahedronVolume", func() { sink, _ = TetrahedronVolume(2) })
	assertNoAllocs(t, "TetrahedronSurfaceArea", func() { sink, _ = TetrahedronSurfaceArea(2) })
}

func BenchmarkTetrahedronVolume(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = TetrahedronVolume(a)
	}
}

func BenchmarkTetrahedronSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = TetrahedronSurfaceArea(a)
	}
}
