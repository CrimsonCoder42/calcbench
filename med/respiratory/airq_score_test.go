package respiratory

import (
	"errors"
	"testing"
)

var (
	sinkInt  int
	sinkAIRQ AIRQControl
)

// parity: ledger airq_asthma, source cases 1–6. The source orders 7 impairment then 3 risk items; the split was
// not seen in the primary's abstract, so only the total is ported. Answers are in the source's order.
func TestAIRQScore_parity(t *testing.T) {
	cases := []struct {
		yes       [10]bool
		wantTotal int
		wantLevel AIRQControl
	}{
		{[10]bool{}, 0, AIRQWellControlled},                                                                  // source case 1
		{[10]bool{0: true}, 1, AIRQWellControlled},                                                           // source case 2
		{[10]bool{0: true, 7: true}, 2, AIRQNotWellControlled},                                               // source case 3
		{[10]bool{0: true, 1: true, 2: true, 3: true}, 4, AIRQNotWellControlled},                             // source case 4
		{[10]bool{0: true, 1: true, 2: true, 7: true, 8: true}, 5, AIRQVeryPoorlyControlled},                 // source case 5
		{[10]bool{true, true, true, true, true, true, true, true, true, true}, 10, AIRQVeryPoorlyControlled}, // source case 6
	}
	for i, c := range cases {
		got := AIRQScore(c.yes)
		if got != c.wantTotal {
			t.Errorf("case %d: total = %d, want %d", i+1, got, c.wantTotal)
		}
		lvl, err := AIRQControlLevel(got)
		if err != nil || lvl != c.wantLevel {
			t.Errorf("case %d: level = (%q, %v), want %q", i+1, lvl, err, c.wantLevel)
		}
	}
}

// Cut points 0-1, 2-4, 5-10 (Murphy 2020 abstract): both sides of each edge, and the ends.
func TestAIRQControlLevel_edges(t *testing.T) {
	cases := []struct {
		total int
		want  AIRQControl
	}{
		{0, AIRQWellControlled}, {1, AIRQWellControlled},
		{2, AIRQNotWellControlled}, {4, AIRQNotWellControlled},
		{5, AIRQVeryPoorlyControlled}, {10, AIRQVeryPoorlyControlled},
	}
	for _, c := range cases {
		if got, err := AIRQControlLevel(c.total); err != nil || got != c.want {
			t.Errorf("%d: got (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 11} {
		if _, err := AIRQControlLevel(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestAIRQScore_noAlloc(t *testing.T) {
	yes := [10]bool{0: true, 5: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = AIRQScore(yes) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkAIRQ, _ = AIRQControlLevel(3) }); n != 0 {
		t.Errorf("level allocs = %v, want 0", n)
	}
}

func BenchmarkAIRQScore(b *testing.B) {
	b.ReportAllocs()
	yes := [10]bool{0: true, 5: true, 9: true}
	for b.Loop() {
		sinkInt = AIRQScore(yes)
	}
}

func BenchmarkAIRQControlLevel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkAIRQ, _ = AIRQControlLevel(3)
	}
}
