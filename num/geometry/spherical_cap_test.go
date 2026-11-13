package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.18.1.1 p. 367 (area 2πrh, volume ⅓πh²(3r − h)), hand-derived for
// r 5, h 2: volume ⅓π·4·13 = 52π/3; area 20π.
func TestSphericalCap_example(t *testing.T) {
	runValCases(t, []valCase{
		{"SphericalCapVolume(5, 2)", func() (float64, error) { return SphericalCapVolume(5, 2) }, 52 * math.Pi / 3, tolFew},
		{"SphericalCapArea(5, 2)", func() (float64, error) { return SphericalCapArea(5, 2) }, 20 * math.Pi, tolFew},
	})
}

func TestSphericalCap_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"hemisphere volume", func() (float64, error) { return SphericalCapVolume(3, 3) }, 18 * math.Pi, tolFew},
		{"hemisphere area", func() (float64, error) { return SphericalCapArea(3, 3) }, 18 * math.Pi, tolFew},
		{"3r would overflow, h² is subnormal", func() (float64, error) { return SphericalCapVolume(1e308, 1e-160) }, math.Pi * 1e-320 * 1e308, tolSome},
	})
}

func TestSphericalCap_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"volume NaN", func() error { return e1(SphericalCapVolume(nan, 1)) }, ErrNotFinite},
		{"volume zero height", func() error { return e1(SphericalCapVolume(1, 0)) }, ErrNonPositive},
		{"volume height beyond the diameter", func() error { return e1(SphericalCapVolume(1, 2.5)) }, ErrOutOfRange},
		{"volume overflow", func() error { return e1(SphericalCapVolume(1e103, 1e103)) }, ErrOverflow},
		{"area Inf", func() error { return e1(SphericalCapArea(inf, 1)) }, ErrNotFinite},
		{"area negative radius", func() error { return e1(SphericalCapArea(-1, 1)) }, ErrNonPositive},
		{"area height beyond the diameter", func() error { return e1(SphericalCapArea(1, 3)) }, ErrOutOfRange},
		{"area overflow", func() error { return e1(SphericalCapArea(1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: the CRC's second forms on the same page, with the cap's base radius a² = h(2r − h):
// volume (1/6)πh(3a² + h²) and area πp², p² = a² + h² (p the slant distance from the pole to the rim).
func TestSphericalCap_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.1 + rng.Float64()*10
		h := 2 * r * (0.01 + 0.98*rng.Float64())
		a2 := h * (2*r - h)
		wantV := math.Pi * h * (3*a2 + h*h) / 6
		wantS := math.Pi * (a2 + h*h)
		if got, err := SphericalCapVolume(r, h); err != nil || !relNear(got, wantV, 1e-14) {
			t.Fatalf("SphericalCapVolume(%v, %v) = (%v, %v), want %v", r, h, got, err, wantV)
		}
		if got, err := SphericalCapArea(r, h); err != nil || !relNear(got, wantS, 1e-14) {
			t.Fatalf("SphericalCapArea(%v, %v) = (%v, %v), want %v", r, h, got, err, wantS)
		}
	}
}

func TestSphericalCap_allocs(t *testing.T) {
	assertNoAllocs(t, "SphericalCapVolume", func() { sink, _ = SphericalCapVolume(5, 2) })
	assertNoAllocs(t, "SphericalCapArea", func() { sink, _ = SphericalCapArea(5, 2) })
}

func BenchmarkSphericalCapVolume(b *testing.B) {
	b.ReportAllocs()
	r, h := 5.0, 2.0
	for b.Loop() {
		sink, _ = SphericalCapVolume(r, h)
	}
}

func BenchmarkSphericalCapArea(b *testing.B) {
	b.ReportAllocs()
	r, h := 5.0, 2.0
	for b.Loop() {
		sink, _ = SphericalCapArea(r, h)
	}
}
