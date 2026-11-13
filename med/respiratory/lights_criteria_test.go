package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from the three criteria (Porcel 2026, PMID 41846683; Subedi 2025, PMID 40867084).
func TestLightsCriteria_example(t *testing.T) {
	cases := []struct {
		pp, sp, pl, sl, uln float64
		met                 int
		exudate             bool
	}{
		{4.0, 6.0, 150, 200, 250, 2, true}, // 0.67 > 0.5, 0.75 > 0.6, 150 < 166.7
		{2.0, 7.0, 80, 200, 250, 0, false}, // 0.29, 0.40, 80 < 166.7
		{4.0, 6.0, 300, 400, 250, 3, true}, // 0.67, 0.75, 300 > 166.7
		{2.0, 7.0, 100, 150, 140, 2, true}, // protein no; 0.67 > 0.6; 100 > 93.3
	}
	for _, c := range cases {
		met, ex, err := LightsCriteria(c.pp, c.sp, c.pl, c.sl, c.uln)
		if err != nil || met != c.met || ex != c.exudate {
			t.Errorf("%v: got (%d, %v, %v), want (%d, %v)", c, met, ex, err, c.met, c.exudate)
		}
	}
}

// Each criterion is strict: a ratio of exactly 0.5, 0.6 or 2/3 does not count; just above does.
func TestLightsCriteria_edges(t *testing.T) {
	cases := []struct {
		name                string
		pp, sp, pl, sl, uln float64
		met                 int
	}{
		{"protein ratio 0.5", 3, 6, 10, 200, 300, 0},
		{"protein ratio 0.5001", 3.0006, 6, 10, 200, 300, 1},
		{"LDH ratio 0.6", 1, 6, 120, 200, 300, 0},
		{"LDH ratio 0.605", 1, 6, 121, 200, 300, 1},
		{"LDH 2/3 of ULN", 1, 6, 100, 1000, 150, 0},
		{"LDH above 2/3 of ULN", 1, 6, 100.1, 1000, 150, 1},
	}
	for _, c := range cases {
		if met, _, err := LightsCriteria(c.pp, c.sp, c.pl, c.sl, c.uln); err != nil || met != c.met {
			t.Errorf("%s: got (%d, %v), want %d", c.name, met, err, c.met)
		}
	}
}

// parity: ledger lights_criteria, source case 1: the source applies the same three ratios with "at least one"
// for exudate; protein 3.5/6.0, LDH 100/200, ULN 222 → protein only → 1, exudate.
func TestLightsCriteria_parity(t *testing.T) {
	if met, ex, err := LightsCriteria(3.5, 6.0, 100, 200, 222); err != nil || met != 1 || !ex {
		t.Errorf("got (%d, %v, %v), want (1, true)", met, ex, err)
	}
}

func TestLightsCriteria_errors(t *testing.T) {
	cases := []struct {
		pp, sp, pl, sl, uln float64
		want                error
	}{
		{0, 6, 100, 200, 250, ErrOutOfRange}, {3, 0, 100, 200, 250, ErrOutOfRange}, {3, 6, -1, 200, 250, ErrOutOfRange},
		{3, 6, 100, 0, 250, ErrOutOfRange}, {3, 6, 100, 200, 0, ErrOutOfRange},
		{math.NaN(), 6, 100, 200, 250, ErrNotFinite}, {3, math.Inf(1), 100, 200, 250, ErrNotFinite},
		{3, 6, math.Inf(-1), 200, 250, ErrNotFinite}, {3, 6, 100, math.NaN(), 250, ErrNotFinite}, {3, 6, 100, 200, math.Inf(1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, _, err := LightsCriteria(c.pp, c.sp, c.pl, c.sl, c.uln); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

func TestLightsCriteria_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _, _ = LightsCriteria(4, 6, 150, 200, 250) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkLightsCriteria(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _, _ = LightsCriteria(4, 6, 150, 200, 250)
	}
}
