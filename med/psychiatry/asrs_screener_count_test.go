package psychiatry

import (
	"errors"
	"testing"
)

// Chamberlain 2021 / Kittel-Schneider 2019: items 1–3 count at ≥ 2, items 4–6 at ≥ 3; each threshold on both sides.
func TestASRSScreenerCount_example(t *testing.T) {
	cases := []struct {
		in   [6]int
		want int
	}{
		{[6]int{2, 2, 2, 3, 2, 2}, 4},
		{[6]int{1, 1, 1, 1, 1, 1}, 0},
		{[6]int{4, 4, 4, 4, 4, 4}, 6},
		{[6]int{2, 1, 0, 3, 3, 3}, 4},
		{[6]int{0, 0, 0, 2, 2, 2}, 0},
	}
	for _, c := range cases {
		if got, err := ASRSScreenerCount(c.in); err != nil || got != c.want {
			t.Errorf("ASRSScreenerCount(%v) = (%d, %v), want %d", c.in, got, err, c.want)
		}
	}
	for _, b := range [][6]int{{5, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, -1}} {
		if _, err := ASRSScreenerCount(b); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%v: %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger asrs_adhd, source cases 1–6 (source asrs_adhd_test.go: zero, sometimes all, threshold
// inattention, threshold mixed, max, mixed; "score" output).
func TestASRSScreenerCount_parity(t *testing.T) {
	cases := []struct {
		in   [6]int
		want int
	}{
		{[6]int{0, 0, 0, 0, 0, 0}, 0}, // parity: ledger asrs_adhd, source case 1
		{[6]int{2, 2, 2, 2, 2, 2}, 3}, // parity: ledger asrs_adhd, source case 2
		{[6]int{2, 2, 2, 3, 0, 0}, 4}, // parity: ledger asrs_adhd, source case 3
		{[6]int{1, 1, 1, 4, 4, 4}, 3}, // parity: ledger asrs_adhd, source case 4
		{[6]int{4, 4, 4, 4, 4, 4}, 6}, // parity: ledger asrs_adhd, source case 5
		{[6]int{3, 1, 3, 2, 3, 4}, 4}, // parity: ledger asrs_adhd, source case 6
	}
	for i, c := range cases {
		if got, err := ASRSScreenerCount(c.in); err != nil || got != c.want {
			t.Errorf("case %d = (%d, %v), want %d", i+1, got, err, c.want)
		}
	}
}

func TestASRSScreenerCount_noAlloc(t *testing.T) {
	in := [6]int{2, 2, 2, 3, 2, 2}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ASRSScreenerCount(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkASRSScreenerCount(b *testing.B) {
	b.ReportAllocs()
	in := [6]int{2, 2, 2, 3, 2, 2}
	for b.Loop() {
		sinkInt, _ = ASRSScreenerCount(in)
	}
}
