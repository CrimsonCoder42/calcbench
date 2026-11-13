package geometry

import "testing"

// Worked examples: CRC 31st ed. §4.14.1 tables, p. 360: V/a³ = 1 and S/a² = 6 for the cube, printed
// to 6 or 7 significant digits (hence 1e-5); and the table's exact forms, V/a³ = 1 and S/a² = 6.
func TestCube_example(t *testing.T) {
	runValCases(t, []valCase{
		{"CubeVolume(1), table", func() (float64, error) { return CubeVolume(1) }, 1, 1e-5},
		{"CubeSurfaceArea(1), table", func() (float64, error) { return CubeSurfaceArea(1) }, 6, 1e-5},
		{"CubeVolume(1), exact form", func() (float64, error) { return CubeVolume(1) }, 1, tolFew},
		{"CubeSurfaceArea(1), exact form", func() (float64, error) { return CubeSurfaceArea(1) }, 6, tolFew},
	})
}

func TestCube_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"volume scales as a³", func() (float64, error) { return CubeVolume(2) }, 8 * (1), tolFew},
		{"area scales as a²", func() (float64, error) { return CubeSurfaceArea(2) }, 4 * (6), tolFew},
		{"tiny edge", func() (float64, error) { return CubeVolume(1e-100) }, (1) * 1e-300, tolFew},
	})
}

func TestCube_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(CubeVolume(nan)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(CubeVolume(0)) }, ErrNonPositive},
		{"volume overflow", func() error { return e1(CubeVolume(1e104)) }, ErrOverflow},
		{"area Inf", func() error { return e1(CubeSurfaceArea(inf)) }, ErrNotFinite},
		{"area negative", func() error { return e1(CubeSurfaceArea(-1)) }, ErrNonPositive},
		{"area overflow", func() error { return e1(CubeSurfaceArea(1e155)) }, ErrOverflow},
	})
}

// Oracles: the surface is 6 regular 4-gons (RegularPolygonAreaFromSide), and the volume is ⅓·r·S with
// the inradius r/a = 0.5 (CRC 31st ed. eq. (4.14.7) and the r/a table).
func TestCube_oracle(t *testing.T) {
	rng := newRNG()
	for range 200 {
		a := 0.01 + rng.Float64()*100
		face, _ := RegularPolygonAreaFromSide(4, a)
		s, err := CubeSurfaceArea(a)
		if err != nil || !relNear(s, 6*face, tolSome) {
			t.Fatalf("CubeSurfaceArea(%v) = (%v, %v), faces %v", a, s, err, 6*face)
		}
		want := s * (0.5) * a / 3
		if got, err := CubeVolume(a); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("CubeVolume(%v) = (%v, %v), ⅓rS %v", a, got, err, want)
		}
	}
}

func TestCube_allocs(t *testing.T) {
	assertNoAllocs(t, "CubeVolume", func() { sink, _ = CubeVolume(2) })
	assertNoAllocs(t, "CubeSurfaceArea", func() { sink, _ = CubeSurfaceArea(2) })
}

func BenchmarkCubeVolume(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = CubeVolume(a)
	}
}

func BenchmarkCubeSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	a := 2.0
	for b.Loop() {
		sink, _ = CubeSurfaceArea(a)
	}
}
