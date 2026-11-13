package respiratory

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Levine 2021 Table 2 (PMID 33274414): the maximum 5 + 21 + 29 = 55 and the minimum 0.
func TestCHOSENScore_tableExample(t *testing.T) {
	cases := []struct {
		age, spo2 int
		alb       float64
		want      int
	}{
		{40, 99, 4.0, 55},
		{80, 92, 2.5, 0},
		{50, 95, 3.5, 2 + 9 + 15},
		{65, 97, 3.0, 1 + 14 + 5},
	}
	for _, c := range cases {
		if got, err := CHOSENScore(c.age, c.spo2, c.alb); err != nil || got != c.want {
			t.Errorf("CHOSENScore(%d, %d, %v) = (%d, %v), want %d", c.age, c.spo2, c.alb, got, err, c.want)
		}
	}
}

// Both sides of every band edge in Table 2.
func TestCHOSENScore_edges(t *testing.T) {
	score := func(age, spo2 int, alb float64) int {
		s, err := CHOSENScore(age, spo2, alb)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return s
	}
	// age: base SpO₂ 90 (0) and albumin 2.0 (0) isolate the age points.
	for _, c := range []struct{ age, want int }{{18, 5}, {45, 5}, {46, 2}, {59, 2}, {60, 1}, {73, 1}, {74, 0}, {120, 0}} {
		if got := score(c.age, 90, 2.0); got != c.want {
			t.Errorf("age %d: got %d, want %d", c.age, got, c.want)
		}
	}
	for _, c := range []struct{ spo2, want int }{{0, 0}, {93, 0}, {94, 9}, {96, 9}, {97, 14}, {98, 14}, {99, 21}, {100, 21}} {
		if got := score(80, c.spo2, 2.0); got != c.want {
			t.Errorf("SpO₂ %d: got %d, want %d", c.spo2, got, c.want)
		}
	}
	for _, c := range []struct {
		alb  float64
		want int
	}{{2.79, 0}, {2.8, 5}, {3.3, 5}, {3.35, 15}, {3.4, 15}, {3.7, 15}, {3.71, 29}} {
		if got := score(80, 90, c.alb); got != c.want {
			t.Errorf("albumin %v: got %d, want %d", c.alb, got, c.want)
		}
	}
}

// parity: ledger chosen_risk_score, source case 1 (full mode, albumin available): age 50 → 2, SpO₂ 97 → 14,
// albumin 3.5 → 15; total 31. The source's modified mode (no albumin, respiratory rate) is not in the paper's
// main text and is not ported (ledger note).
func TestCHOSENScore_parity(t *testing.T) {
	if got, err := CHOSENScore(50, 97, 3.5); err != nil || got != 31 {
		t.Errorf("got (%d, %v), want 31", got, err)
	}
}

func TestCHOSENScore_errors(t *testing.T) {
	cases := []struct {
		age, spo2 int
		alb       float64
		want      error
	}{
		{17, 95, 3.5, ErrOutOfRange}, {121, 95, 3.5, ErrOutOfRange},
		{50, -1, 3.5, ErrOutOfRange}, {50, 101, 3.5, ErrOutOfRange},
		{50, 95, 0, ErrOutOfRange}, {50, 95, -1, ErrOutOfRange},
		{50, 95, math.NaN(), ErrNotFinite}, {50, 95, math.Inf(1), ErrNotFinite}, {50, 95, math.Inf(-1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := CHOSENScore(c.age, c.spo2, c.alb); !errors.Is(err, c.want) {
			t.Errorf("CHOSENScore(%d, %d, %v): err = %v, want %v", c.age, c.spo2, c.alb, err, c.want)
		}
	}
}

func TestCHOSENScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CHOSENScore(50, 97, 3.5) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCHOSENScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = CHOSENScore(50, 97, 3.5)
	}
}
