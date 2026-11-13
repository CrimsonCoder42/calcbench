package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPRESET PRESETRisk

// HAND-DERIVED from Hilder 2017 Table 4 (PMID 29233160): MAP 75 (3), lactate 4 (2), pH 7.15 (2), platelets 150
// (1), 5 days (1) = 9; the extremes 0 and 15.
func TestPRESETScore_tableExample(t *testing.T) {
	cases := []struct {
		mapMmHg, lac, ph, plt float64
		days, want            int
	}{
		{75, 4, 7.15, 150, 5, 9},
		{110, 1, 7.4, 250, 1, 0},
		{60, 12, 7.0, 50, 10, 15},
	}
	for _, c := range cases {
		if got, err := PRESETScore(c.mapMmHg, c.lac, c.ph, c.plt, c.days); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c, got, err, c.want)
		}
	}
}

// Both sides of every band edge of Table 4, one variable at a time from the all-zero base.
func TestPRESETScore_edges(t *testing.T) {
	score := func(m, l, p, pl float64, d int) int {
		s, err := PRESETScore(m, l, p, pl, d)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return s
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{100.1, 0}, {100, 1}, {91, 1}, {90, 2}, {81, 2}, {80, 3}, {71, 3}, {70, 4}} {
		if got := score(c.v, 1, 7.4, 250, 0); got != c.want {
			t.Errorf("MAP %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{1.5, 0}, {1.51, 1}, {3, 1}, {3.01, 2}, {6, 2}, {6.01, 3}, {10, 3}, {10.01, 4}} {
		if got := score(110, c.v, 7.4, 250, 0); got != c.want {
			t.Errorf("lactate %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{7.301, 0}, {7.3, 1}, {7.201, 1}, {7.2, 2}, {7.101, 2}, {7.1, 3}} {
		if got := score(110, 1, c.v, 250, 0); got != c.want {
			t.Errorf("pH %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct {
		v    float64
		want int
	}{{201, 0}, {200, 1}, {101, 1}, {100, 2}} {
		if got := score(110, 1, 7.4, c.v, 0); got != c.want {
			t.Errorf("platelets %v: got %d, want %d", c.v, got, c.want)
		}
	}
	for _, c := range []struct{ d, want int }{{2, 0}, {3, 1}, {7, 1}, {8, 2}} {
		if got := score(110, 1, 7.4, 250, c.d); got != c.want {
			t.Errorf("days %d: got %d, want %d", c.d, got, c.want)
		}
	}
}

// parity: ledger preset_ecmo, source case 1: the source takes band indices with the same points as Table 4
// (MAP [0..4], lactate [0..4], pH [0..3], platelets [0..2], days [0..2]); bands 3,2,2,1,1 → 9.
func TestPRESETScore_parity(t *testing.T) {
	if got, err := PRESETScore(75, 4, 7.15, 150, 5); err != nil || got != 9 {
		t.Errorf("got (%d, %v), want 9", got, err)
	}
}

func TestPRESETScore_errors(t *testing.T) {
	cases := []struct {
		m, l, p, pl float64
		d           int
		want        error
	}{
		{0, 1, 7.4, 250, 0, ErrOutOfRange}, {110, -1, 7.4, 250, 0, ErrOutOfRange}, {110, 1, 6.4, 250, 0, ErrOutOfRange},
		{110, 1, 8.1, 250, 0, ErrOutOfRange}, {110, 1, 7.4, -1, 0, ErrOutOfRange}, {110, 1, 7.4, 250, -1, ErrOutOfRange},
		{math.NaN(), 1, 7.4, 250, 0, ErrNotFinite}, {110, math.Inf(1), 7.4, 250, 0, ErrNotFinite},
		{110, 1, math.Inf(-1), 250, 0, ErrNotFinite}, {110, 1, 7.4, math.NaN(), 0, ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := PRESETScore(c.m, c.l, c.p, c.pl, c.d); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

// Classes: 0–5 I, 6–9 II, 10–15 III.
func TestPRESETRiskClass_edges(t *testing.T) {
	cases := []struct {
		s    int
		want PRESETRisk
	}{{0, PRESETClassI}, {5, PRESETClassI}, {6, PRESETClassII}, {9, PRESETClassII}, {10, PRESETClassIII}, {15, PRESETClassIII}}
	for _, c := range cases {
		if got, err := PRESETRiskClass(c.s); err != nil || got != c.want {
			t.Errorf("%d: got (%q, %v), want %q", c.s, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 16} {
		if _, err := PRESETRiskClass(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestPRESETScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PRESETScore(75, 4, 7.15, 150, 5) }); n != 0 {
		t.Errorf("PRESETScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPRESET, _ = PRESETRiskClass(9) }); n != 0 {
		t.Errorf("PRESETRiskClass allocs = %v, want 0", n)
	}
}

func BenchmarkPRESETScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = PRESETScore(75, 4, 7.15, 150, 5)
	}
}

func BenchmarkPRESETRiskClass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPRESET, _ = PRESETRiskClass(9)
	}
}
