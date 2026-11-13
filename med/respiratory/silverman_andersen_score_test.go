package respiratory

import (
	"errors"
	"testing"
)

// HAND-DERIVED: the sum of five 0–2 grades (Brenne 2023, PMC10313036; Nussbaum 2022, PMID 36138599);
// "Zero points indicate normal respiration and ten points severe respiratory distress".
func TestSilvermanAndersenScore_example(t *testing.T) {
	cases := []struct {
		g    [5]int
		want int
	}{{[5]int{1, 1, 0, 2, 1}, 5}, {[5]int{}, 0}, {[5]int{2, 2, 2, 2, 2}, 10}}
	for _, c := range cases {
		if got, err := SilvermanAndersenScore(c.g); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c.g, got, err, c.want)
		}
	}
}

// parity: ledger silverman_calc, source case 1: the source rounds and sums the five grades; 2,1,1,0,0 → 4.
func TestSilvermanAndersenScore_parity(t *testing.T) {
	if got, err := SilvermanAndersenScore([5]int{2, 1, 1, 0, 0}); err != nil || got != 4 {
		t.Errorf("got (%d, %v), want 4", got, err)
	}
}

func TestSilvermanAndersenScore_errors(t *testing.T) {
	for _, g := range [][5]int{{-1}, {3}, {0, 0, 0, 0, 3}} {
		if _, err := SilvermanAndersenScore(g); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%v: err = %v, want ErrOutOfRange", g, err)
		}
	}
}

func TestSilvermanAndersenScore_noAlloc(t *testing.T) {
	g := [5]int{1, 1, 0, 2, 1}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SilvermanAndersenScore(g) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSilvermanAndersenScore(b *testing.B) {
	g := [5]int{1, 1, 0, 2, 1}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = SilvermanAndersenScore(g)
	}
}
