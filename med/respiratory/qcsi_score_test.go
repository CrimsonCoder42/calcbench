package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Haimovich 2020 Table 2 (PMID 33012378): rate 26 (1), SpO₂ 90 (2), flow 4 (4) = 7; the
// extremes 0 and 2 + 5 + 5 = 12 ("0 to 12 points").
func TestQCSIScore_tableExample(t *testing.T) {
	cases := []struct {
		rr   float64
		spo2 int
		flow float64
		want int
	}{{26, 90, 4, 7}, {18, 96, 0, 0}, {30, 85, 6, 12}}
	for _, c := range cases {
		if got, err := QCSIScore(c.rr, c.spo2, c.flow); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c, got, err, c.want)
		}
	}
}

// Both sides of each band edge.
func TestQCSIScore_edges(t *testing.T) {
	score := func(rr float64, spo2 int, flow float64) int {
		s, err := QCSIScore(rr, spo2, flow)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return s
	}
	for _, c := range []struct {
		rr   float64
		want int
	}{{22, 0}, {23, 1}, {28, 1}, {29, 2}} {
		if got := score(c.rr, 96, 0); got != c.want {
			t.Errorf("rate %v: got %d, want %d", c.rr, got, c.want)
		}
	}
	for _, c := range []struct{ spo2, want int }{{93, 0}, {92, 2}, {89, 2}, {88, 5}, {0, 5}, {100, 0}} {
		if got := score(18, c.spo2, 0); got != c.want {
			t.Errorf("SpO₂ %d: got %d, want %d", c.spo2, got, c.want)
		}
	}
	for _, c := range []struct {
		flow float64
		want int
	}{{2, 0}, {3, 4}, {4, 4}, {5, 5}, {6, 5}} {
		if got := score(18, 96, c.flow); got != c.want {
			t.Errorf("flow %v: got %d, want %d", c.flow, got, c.want)
		}
	}
}

// parity: ledger qcsi_covid, source case 1: band indices RR [0,1,2], SpO₂ [0,2,5], O₂ [0,4,5] agree with
// Table 2; bands (1, 1, 1) → 1 + 2 + 4 = 7.
func TestQCSIScore_parity(t *testing.T) {
	if got, err := QCSIScore(25, 91, 3); err != nil || got != 7 {
		t.Errorf("got (%d, %v), want 7", got, err)
	}
}

func TestQCSIScore_errors(t *testing.T) {
	cases := []struct {
		rr   float64
		spo2 int
		flow float64
		want error
	}{
		{-1, 96, 0, ErrOutOfRange}, {18, -1, 0, ErrOutOfRange}, {18, 101, 0, ErrOutOfRange},
		{18, 96, -0.1, ErrOutOfRange}, {18, 96, 6.1, ErrOutOfRange},
		{math.NaN(), 96, 0, ErrNotFinite}, {math.Inf(1), 96, 0, ErrNotFinite}, {18, 96, math.Inf(-1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := QCSIScore(c.rr, c.spo2, c.flow); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

func TestQCSIScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = QCSIScore(26, 90, 4) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkQCSIScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = QCSIScore(26, 90, 4)
	}
}
