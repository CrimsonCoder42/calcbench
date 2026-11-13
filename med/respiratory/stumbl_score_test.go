package respiratory

import (
	"errors"
	"math"
	"testing"
)

// The paper's worked example (Battle 2014, PMID 24887537, Results): 67 years (6) + three rib fractures (9) +
// chronic lung disease (5) + saturation 87 % (4) = 24.
func TestSTUMBLScore_paperExample(t *testing.T) {
	if got, err := STUMBLScore(67, 3, true, false, 87); err != nil || got != 24 {
		t.Errorf("got (%d, %v), want 24", got, err)
	}
}

// Age: ⌊age/10⌋ (both sides of 20 and 30); saturation: both sides of 95, 90, 85, 80.
func TestSTUMBLScore_edges(t *testing.T) {
	score := func(age float64, spo2 int) int {
		s, err := STUMBLScore(age, 0, false, false, spo2)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return s
	}
	for _, c := range []struct {
		age  float64
		want int
	}{{18, 1}, {19.9, 1}, {20, 2}, {29.9, 2}, {30, 3}, {120, 12}} {
		if got := score(c.age, 98); got != c.want {
			t.Errorf("age %v: got %d, want %d", c.age, got, c.want)
		}
	}
	for _, c := range []struct{ spo2, want int }{{100, 0}, {95, 0}, {94, 2}, {90, 2}, {89, 4}, {85, 4}, {84, 6}, {80, 6}, {79, 8}, {0, 38}} {
		if got := score(18, c.spo2) - 1; got != c.want {
			t.Errorf("SpO₂ %d: got %d, want %d", c.spo2, got, c.want)
		}
	}
	if got, _ := STUMBLScore(40, 2, true, true, 96); got != 4+6+5+4 {
		t.Errorf("40 y, 2 fractures, CLD, anticoagulant: got %d, want 19", got)
	}
}

// parity: ledger stumbl_battle_score, source case 1: the source scores ⌊age/10⌋, 3 per fracture, 5, 4 and
// 2 × the saturation band index (≥ 95, 90–94, 85–89, 80–84) — the same as Table 4 read with 94 % in the first
// scored band; 67 y, 3 fractures, CLD, band 85–89 → 24.
func TestSTUMBLScore_parity(t *testing.T) {
	if got, err := STUMBLScore(67, 3, true, false, 86); err != nil || got != 24 {
		t.Errorf("got (%d, %v), want 24", got, err)
	}
}

func TestSTUMBLScore_errors(t *testing.T) {
	cases := []struct {
		age        float64
		ribs, spo2 int
		want       error
	}{
		{17.9, 0, 98, ErrOutOfRange}, {121, 0, 98, ErrOutOfRange}, {40, -1, 98, ErrOutOfRange}, {40, 25, 98, ErrOutOfRange},
		{40, 0, -1, ErrOutOfRange}, {40, 0, 101, ErrOutOfRange},
		{math.NaN(), 0, 98, ErrNotFinite}, {math.Inf(1), 0, 98, ErrNotFinite}, {math.Inf(-1), 0, 98, ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := STUMBLScore(c.age, c.ribs, false, false, c.spo2); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c, err, c.want)
		}
	}
}

func TestSTUMBLScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = STUMBLScore(67, 3, true, false, 87) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSTUMBLScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = STUMBLScore(67, 3, true, false, 87)
	}
}
