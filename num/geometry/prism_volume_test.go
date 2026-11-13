package geometry

import "testing"

// Worked example: CRC 31st ed. §4.15 eq. (4.15.1) (V = Ah), hand-derived: base 12, height 5 → 60.
func TestPrismVolume_example(t *testing.T) {
	runValCases(t, []valCase{
		{"PrismVolume(12, 5)", func() (float64, error) { return PrismVolume(12, 5) }, 60, 0},
	})
}

func TestPrismVolume_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"cube as a prism", func() (float64, error) { return PrismVolume(4, 2) }, 8, 0},
		{"no spurious overflow", func() (float64, error) { return PrismVolume(1e300, 1e-300) }, 1, tolFew},
	})
}

func TestPrismVolume_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(PrismVolume(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(PrismVolume(1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(PrismVolume(0, 1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(PrismVolume(1e200, 1e200)) }, ErrOverflow},
	})
}

// Oracle: a prism over a circle is the cylinder.
func TestPrismVolume_vsOtherSolids(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r, h := 0.1+rng.Float64()*10, 0.1+rng.Float64()*10
		base, _ := CircleArea(r)
		cyl, _ := CylinderVolume(r, h)
		if got, err := PrismVolume(base, h); err != nil || !relNear(got, cyl, tolSome) {
			t.Fatalf("PrismVolume(πr², h), r %v, h %v = (%v, %v), cylinder %v", r, h, got, err, cyl)
		}
	}
}

func TestPrismVolume_allocs(t *testing.T) {
	assertNoAllocs(t, "PrismVolume", func() { sink, _ = PrismVolume(12, 5) })
}

func BenchmarkPrismVolume(b *testing.B) {
	b.ReportAllocs()
	base, h := 12.0, 5.0
	for b.Loop() {
		sink, _ = PrismVolume(base, h)
	}
}
