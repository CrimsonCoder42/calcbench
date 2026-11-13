package geometry

import "testing"

// Worked example: CRC 31st ed. §4.5.2 eq. (4.5.9) with θ = π/2 (and p. 323, rhombus: area = ½pq),
// hand-derived: diagonals 6 and 4 → 12.
func TestKiteArea_example(t *testing.T) {
	runValCases(t, []valCase{
		{"KiteArea(6, 4)", func() (float64, error) { return KiteArea(6, 4) }, 12, 0},
	})
}

func TestKiteArea_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"square of side 1 by its diagonals √2", func() (float64, error) { return KiteArea(1.4142135623730951, 1.4142135623730951) }, 1, tolFew},
		{"no spurious overflow", func() (float64, error) { return KiteArea(1e300, 1e-300) }, 0.5, tolFew},
	})
}

func TestKiteArea_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(KiteArea(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(KiteArea(1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(KiteArea(0, 1)) }, ErrNonPositive},
		{"overflow", func() error { return e1(KiteArea(1e300, 1e300)) }, ErrOverflow},
	})
}

// Oracle: the kite with diagonals d₁ (cut at fraction f) and d₂ as a polygon, by the shoelace formula.
func TestKiteArea_vsPolygon(t *testing.T) {
	rng := newRNG()
	for range 400 {
		d1, d2, f := rng.Float64()*10+0.1, rng.Float64()*10+0.1, 0.05+rng.Float64()*0.9
		kite := []Point{{0, 0}, {f * d1, d2 / 2}, {d1, 0}, {f * d1, -d2 / 2}}
		want, err := PolygonArea(kite)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := KiteArea(d1, d2); err != nil || !relNear(got, want, tolSome) {
			t.Fatalf("KiteArea(%v, %v) = (%v, %v), polygon %v", d1, d2, got, err, want)
		}
	}
}

func TestKiteArea_allocs(t *testing.T) {
	assertNoAllocs(t, "KiteArea", func() { sink, _ = KiteArea(6, 4) })
}

func BenchmarkKiteArea(b *testing.B) {
	b.ReportAllocs()
	d1, d2 := 6.0, 4.0
	for b.Loop() {
		sink, _ = KiteArea(d1, d2)
	}
}
