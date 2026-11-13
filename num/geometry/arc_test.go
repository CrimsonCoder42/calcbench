package geometry

import (
	"math"
	"testing"
)

// Worked examples: CRC 31st ed. §4.6.6 p. 334 (arc rθ, chord 2R sin ½θ, sector ½R²θ), hand-derived for
// r 2, θ π/3: arc 2π/3, chord 2 (an equilateral triangle), sector 2π/3.
func TestArc_example(t *testing.T) {
	runValCases(t, []valCase{
		{"ArcLength(2, π/3)", func() (float64, error) { return ArcLength(2, math.Pi/3) }, 2 * math.Pi / 3, tolFew},
		{"ChordLength(2, π/3)", func() (float64, error) { return ChordLength(2, math.Pi/3) }, 2, tolFew},
		{"SectorArea(2, π/3)", func() (float64, error) { return SectorArea(2, math.Pi/3) }, 2 * math.Pi / 3, tolFew},
	})
}

func TestArc_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"full-circle arc is the circumference", func() (float64, error) { return ArcLength(3, 2*math.Pi) }, 6 * math.Pi, tolFew},
		{"chord of a half circle is the diameter", func() (float64, error) { return ChordLength(3, math.Pi) }, 6, tolFew},
		{"chord of a full circle is 0", func() (float64, error) { return ChordLength(3, 2*math.Pi) }, 0, 1e-15},
		{"chord of a tiny angle keeps its digits", func() (float64, error) { return ChordLength(1, 1e-10) }, 1e-10, tolFew},
		{"full-circle sector is the disc", func() (float64, error) { return SectorArea(3, 2*math.Pi) }, 9 * math.Pi, tolFew},
	})
}

func TestArc_errors(t *testing.T) {
	above := math.Nextafter(2*math.Pi, 7)
	runErrCases(t, []errCase{
		{"arc NaN", func() error { return e1(ArcLength(nan, 1)) }, ErrNotFinite},
		{"arc Inf angle", func() error { return e1(ArcLength(1, inf)) }, ErrNotFinite},
		{"arc zero radius", func() error { return e1(ArcLength(0, 1)) }, ErrNonPositive},
		{"arc zero angle", func() error { return e1(ArcLength(1, 0)) }, ErrOutOfRange},
		{"arc beyond a turn", func() error { return e1(ArcLength(1, above)) }, ErrOutOfRange},
		{"arc overflow", func() error { return e1(ArcLength(math.MaxFloat64, 6)) }, ErrOverflow},
		{"chord NaN", func() error { return e1(ChordLength(1, nan)) }, ErrNotFinite},
		{"chord negative radius", func() error { return e1(ChordLength(-1, 1)) }, ErrNonPositive},
		{"chord negative angle", func() error { return e1(ChordLength(1, -1)) }, ErrOutOfRange},
		{"chord overflow", func() error { return e1(ChordLength(math.MaxFloat64, math.Pi)) }, ErrOverflow},
		{"sector Inf", func() error { return e1(SectorArea(inf, 1)) }, ErrNotFinite},
		{"sector zero radius", func() error { return e1(SectorArea(0, 1)) }, ErrNonPositive},
		{"sector beyond a turn", func() error { return e1(SectorArea(1, 7)) }, ErrOutOfRange},
		{"sector overflow", func() error { return e1(SectorArea(1e160, 1)) }, ErrOverflow},
	})
}

// Oracles: the chord against the law of cosines √(2r²(1 − cos θ)) away from θ = 0 (where that form
// cancels), and the sector as the segment plus its isosceles triangle ½r² sin θ.
func TestArc_oracle(t *testing.T) {
	rng := newRNG()
	for range 400 {
		r := 0.1 + rng.Float64()*100
		theta := 0.1 + rng.Float64()*(2*math.Pi-0.2)
		want := math.Sqrt(2 * r * r * (1 - math.Cos(theta)))
		if got, err := ChordLength(r, theta); err != nil || !relNear(got, want, 1e-13) {
			t.Fatalf("ChordLength(%v, %v) = (%v, %v), want %v", r, theta, got, err, want)
		}
		sector, err := SectorArea(r, theta)
		segment, _ := SegmentArea(r, theta)
		triangle := 0.5 * r * r * math.Sin(theta)
		if err != nil || math.Abs(sector-(segment+triangle)) > 1e-13*r*r {
			t.Fatalf("SectorArea(%v, %v) = (%v, %v), segment + triangle %v", r, theta, sector, err, segment+triangle)
		}
	}
}

func TestArc_allocs(t *testing.T) {
	assertNoAllocs(t, "ArcLength", func() { sink, _ = ArcLength(2, 1) })
	assertNoAllocs(t, "ChordLength", func() { sink, _ = ChordLength(2, 1) })
	assertNoAllocs(t, "SectorArea", func() { sink, _ = SectorArea(2, 1) })
}

func BenchmarkArcLength(b *testing.B) {
	b.ReportAllocs()
	r, theta := 2.0, 1.0
	for b.Loop() {
		sink, _ = ArcLength(r, theta)
	}
}

func BenchmarkChordLength(b *testing.B) {
	b.ReportAllocs()
	r, theta := 2.0, 1.0
	for b.Loop() {
		sink, _ = ChordLength(r, theta)
	}
}

func BenchmarkSectorArea(b *testing.B) {
	b.ReportAllocs()
	r, theta := 2.0, 1.0
	for b.Loop() {
		sink, _ = SectorArea(r, theta)
	}
}
