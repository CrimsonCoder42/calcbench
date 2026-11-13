package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.14.1 tables, p. 360: V/a³ = 2.181695 and S/a² = 8.66025 for the icosahedron, printed
// to 6 or 7 significant digits (hence 1e-5); and the table's exact forms, V/a³ = 5 * (3 + math.Sqrt(5)) / 12 and S/a² = 5 * math.Sqrt(3).
func TestIcosahedron_example(t *testing.T) {
	runValCases(t, []valCase{
		{"IcosahedronVolume(1), table", func() (float64, error) { return IcosahedronVolume(1) }, 2.181695, 1e-5},
		{"IcosahedronSurfaceArea(1), table", func() (float64, error) { return IcosahedronSurfaceArea(1) }, 8.66025, 1e-5},
		{"IcosahedronVolume(1), exact form", func() (float64, error) { return IcosahedronVolume(1) }, 5 * (3 + math.Sqrt(5)) / 12, tolFew},
		{"IcosahedronSurfaceArea(1), exact form", func() (float64, error) { return IcosahedronSurfaceArea(1) }, 5 * math.Sqrt(3), tolFew},
	})
}

func TestIcosahedron_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume scales as a³", func() (float64, error) { return IcosahedronVolume(2) }, 8 * (5 * (3 + math.Sqrt(5)) / 12), tolFew},
		{"area scales as a²", func() (float64, error) { return IcosahedronSurfaceArea(2) }, 4 * (5 * math.Sqrt(3)), tolFew},
		{"tiny edge", func() (float64, error) { return IcosahedronVolume(1e-100) }, (5 * (3 + math.Sqrt(5)) / 12) * 1e-300, tolFew},
	})
}

func TestIcosahedron_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(IcosahedronVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(IcosahedronVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(IcosahedronVolume(1e104)) }, ErrOverflow},
		{"area Inf", func() error { return e1(IcosahedronSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(IcosahedronSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(IcosahedronSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracles: the surface is 20 regular 3-gons (RegularPolygonAreaFromSide), and the volume is ⅓·r·S with
// the inradius r/a = math.Sqrt(42+18*math.Sqrt(5)) / 12 (CRC 31st ed. eq. (4.14.7) and the r/a table).
func TestIcosahedron_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.01 + rng.Float64()*100
		face, _ := RegularPolygonAreaFromSide(3, a)
		s, err := IcosahedronSurfaceArea(a)
		if err != nil || !relNear(s, 20*face, tolSome) {
			t.Fatalf("IcosahedronSurfaceArea(%v) = (%v, %v), faces %v", a, s, err, 20*face)
		}
		want := s * (math.Sqrt(42+18*math.Sqrt(5)) / 12) * a / 3
		if got, err := IcosahedronVolume(a); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("IcosahedronVolume(%v) = (%v, %v), ⅓rS %v", a, got, err, want)
		}
	}
}

func TestIcosahedron_allocs(t *testing.T) {
	assertNoAllocs(t, "IcosahedronVolume", func() { sink, _ = IcosahedronVolume(2) })
	assertNoAllocs(t, "IcosahedronSurfaceArea", func() { sink, _ = IcosahedronSurfaceArea(2) })
}

func BenchmarkIcosahedronVolume(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = IcosahedronVolume(a)
	}
}

func BenchmarkIcosahedronSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = IcosahedronSurfaceArea(a)
	}
}
