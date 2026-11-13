package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.14.1 tables, p. 360: V/a³ = 7.663119 and S/a² = 20.64573 for the dodecahedron, printed
// to 6 or 7 significant digits (hence 1e-5); and the table's exact forms, V/a³ = (15 + 7*math.Sqrt(5)) / 4 and S/a² = 3 * math.Sqrt(25+10*math.Sqrt(5)).
func TestDodecahedron_example(t *testing.T) {
	runValCases(t, []valCase{
		{"DodecahedronVolume(1), table", func() (float64, error) { return DodecahedronVolume(1) }, 7.663119, 1e-5},
		{"DodecahedronSurfaceArea(1), table", func() (float64, error) { return DodecahedronSurfaceArea(1) }, 20.64573, 1e-5},
		{"DodecahedronVolume(1), exact form", func() (float64, error) { return DodecahedronVolume(1) }, (15 + 7*math.Sqrt(5)) / 4, tolFew},
		{"DodecahedronSurfaceArea(1), exact form", func() (float64, error) { return DodecahedronSurfaceArea(1) }, 3 * math.Sqrt(25+10*math.Sqrt(5)), tolFew},
	})
}

func TestDodecahedron_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume scales as a³", func() (float64, error) { return DodecahedronVolume(2) }, 8 * ((15 + 7*math.Sqrt(5)) / 4), tolFew},
		{"area scales as a²", func() (float64, error) { return DodecahedronSurfaceArea(2) }, 4 * (3 * math.Sqrt(25+10*math.Sqrt(5))), tolFew},
		{"tiny edge", func() (float64, error) { return DodecahedronVolume(1e-100) }, ((15 + 7*math.Sqrt(5)) / 4) * 1e-300, tolFew},
	})
}

func TestDodecahedron_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(DodecahedronVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(DodecahedronVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(DodecahedronVolume(1e104)) }, ErrOverflow},
		{"area Inf", func() error { return e1(DodecahedronSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(DodecahedronSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(DodecahedronSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracles: the surface is 12 regular 5-gons (RegularPolygonAreaFromSide), and the volume is ⅓·r·S with
// the inradius r/a = math.Sqrt(250+110*math.Sqrt(5)) / 20 (CRC 31st ed. eq. (4.14.7) and the r/a table).
func TestDodecahedron_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.01 + rng.Float64()*100
		face, _ := RegularPolygonAreaFromSide(5, a)
		s, err := DodecahedronSurfaceArea(a)
		if err != nil || !relNear(s, 12*face, tolSome) {
			t.Fatalf("DodecahedronSurfaceArea(%v) = (%v, %v), faces %v", a, s, err, 12*face)
		}
		want := s * (math.Sqrt(250+110*math.Sqrt(5)) / 20) * a / 3
		if got, err := DodecahedronVolume(a); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("DodecahedronVolume(%v) = (%v, %v), ⅓rS %v", a, got, err, want)
		}
	}
}

func TestDodecahedron_allocs(t *testing.T) {
	assertNoAllocs(t, "DodecahedronVolume", func() { sink, _ = DodecahedronVolume(2) })
	assertNoAllocs(t, "DodecahedronSurfaceArea", func() { sink, _ = DodecahedronSurfaceArea(2) })
}

func BenchmarkDodecahedronVolume(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = DodecahedronVolume(a)
	}
}

func BenchmarkDodecahedronSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = DodecahedronSurfaceArea(a)
	}
}
