package psychiatry

import (
	"errors"
	"testing"
)

var (
	sinkInt    int
	sinkString string
)

// Butler 2007: 17 items, 0–4 each, summed. HAND-DERIVED: all 0 → 0; all 4 → 68; nine 1s → 9 (the cut-off value).
func TestCOMMScore_example(t *testing.T) {
	var zero, max, nine [17]int
	for i := range max {
		max[i] = 4
	}
	for i := 0; i < 9; i++ {
		nine[i] = 1
	}
	for _, tc := range []struct {
		in   [17]int
		want int
	}{{zero, 0}, {max, 68}, {nine, 9}} {
		if got, err := COMMScore(tc.in); err != nil || got != tc.want {
			t.Errorf("COMMScore(%v) = (%d, %v), want %d", tc.in, got, err, tc.want)
		}
	}
}

// parity: ledger comm_opioid_misuse, source cases 1–6 (the total; the source's five domain sub-sums are not
// ported, ledger note).
func TestCOMMScore_parity(t *testing.T) {
	tests := []struct {
		in   [17]int
		want int
	}{
		{[17]int{}, 0},                                                   // source case 1
		{[17]int{1, 1, 1, 1, 0, 1, 1, 1, 1}, 8},                          // source case 2
		{[17]int{1, 1, 1, 1, 1, 1, 1, 1, 1}, 9},                          // source case 3
		{[17]int{2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 1}, 24},             // source case 4
		{[17]int{2, 2, 2, 2, 1, 2, 2, 2, 2, 2, 2, 2, 1, 1}, 25},          // source case 5
		{[17]int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}, 68}, // source case 6
	}
	for i, tc := range tests {
		if got, err := COMMScore(tc.in); err != nil || got != tc.want {
			t.Errorf("source case %d: COMMScore = (%d, %v), want %d", i+1, got, err, tc.want)
		}
	}
}

func TestCOMMScore_errors(t *testing.T) {
	for i := range 17 {
		for _, v := range []int{-1, 5} {
			var in [17]int
			in[i] = v
			if _, err := COMMScore(in); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("items[%d] = %d: error = %v, want ErrOutOfRange", i, v, err)
			}
		}
	}
}

func TestCOMMScore_noAlloc(t *testing.T) {
	in := [17]int{1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = COMMScore(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCOMMScore(b *testing.B) {
	b.ReportAllocs()
	in := [17]int{1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 0, 1, 2}
	for b.Loop() {
		sinkInt, _ = COMMScore(in)
	}
}
