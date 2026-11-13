package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.17 p. 363 (area 4π²Rr, volume 2π²Rr²), hand-derived for R 3, r 1:
// volume 6π², area 12π².
func TestTorus_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TorusVolume(3, 1)", func() (float64, error) { return TorusVolume(3, 1) }, 6 * math.Pi * math.Pi, tolFew},
		{"TorusSurfaceArea(3, 1)", func() (float64, error) { return TorusSurfaceArea(3, 1) }, 12 * math.Pi * math.Pi, tolFew},
	})
}

func TestTorus_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"horn torus (r = R) is allowed", func() (float64, error) { return TorusVolume(2, 2) }, 16 * math.Pi * math.Pi, tolFew},
		{"thin ring", func() (float64, error) { return TorusSurfaceArea(1e10, 1e-10) }, 4 * math.Pi * math.Pi, tolFew},
	})
}

func TestTorus_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(TorusVolume(nan, 1)) }, ErrNotFinite},
		{"volume zero", func() error { return e1(TorusVolume(3, 0)) }, ErrNonPositive},
		{"volume spindle (r > R)", func() error { return e1(TorusVolume(1, 2)) }, ErrOutOfRange},
		{"volume overflow", func() error { return e1(TorusVolume(1e200, 1e200)) }, ErrOverflow},
		{"area Inf", func() error { return e1(TorusSurfaceArea(inf, 1)) }, ErrNotFinite},
		{"area negative", func() error { return e1(TorusSurfaceArea(3, -1)) }, ErrNonPositive},
		{"area spindle", func() error { return e1(TorusSurfaceArea(1, 1.5)) }, ErrOutOfRange},
		{"area overflow", func() error { return e1(TorusSurfaceArea(1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: the Pappus–Guldinus theorem the CRC states beside the formulas: the disc's area (or the circle's
// length) times the length of the path of its centre.
func TestTorus_pappus(t *testing.T) {
	rng := newRNG()
	for range 400 {
		big := 0.1 + rng.Float64()*10
		small := big * rng.Float64()
		if small == 0 {
			continue
		}
		path, _ := CircleCircumference(big)
		disc, _ := CircleArea(small)
		rim, _ := CircleCircumference(small)
		if got, err := TorusVolume(big, small); err != nil || !relNear(got, disc*path, tolSome) {
			t.Fatalf("TorusVolume(%v, %v) = (%v, %v), want %v", big, small, got, err, disc*path)
		}
		if got, err := TorusSurfaceArea(big, small); err != nil || !relNear(got, rim*path, tolSome) {
			t.Fatalf("TorusSurfaceArea(%v, %v) = (%v, %v), want %v", big, small, got, err, rim*path)
		}
	}
}

func TestTorus_allocs(t *testing.T) {
	assertNoAllocs(t, "TorusVolume", func() { sink, _ = TorusVolume(3, 1) })
	assertNoAllocs(t, "TorusSurfaceArea", func() { sink, _ = TorusSurfaceArea(3, 1) })
}

func BenchmarkTorusVolume(b *testing.B) {
	b.ReportAllocs()
	major, minor := 3.0, 1.0
	for b.Loop() {
		sink, _ = TorusVolume(major, minor)
	}
}

func BenchmarkTorusSurfaceArea(b *testing.B) {
	b.ReportAllocs()
	major, minor := 3.0, 1.0
	for b.Loop() {
		sink, _ = TorusSurfaceArea(major, minor)
	}
}
