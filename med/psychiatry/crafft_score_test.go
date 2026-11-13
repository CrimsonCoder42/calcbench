package psychiatry

import "testing"

// Knight 1999: six yes/no items, the count. HAND-DERIVED: none → 0; all → 6.
func TestCRAFFTScore_example(t *testing.T) {
	if got := CRAFFTScore([6]bool{}); got != 0 {
		t.Errorf("none = %d, want 0", got)
	}
	if got := CRAFFTScore([6]bool{true, true, true, true, true, true}); got != 6 {
		t.Errorf("all = %d, want 6", got)
	}
}

// parity: ledger crafft, source cases 1–6. Items in order car, relax, alone, forget, friends, trouble. The
// source first asks three prior-use questions (CRAFFT 2.1) and counts only "car" when all are "no"; the 1999
// primary has no such gate, so case 2 asserts the primary's count (6), not the source's 1 (ledger note). The
// prior-use answers of the other cases do not change the count.
func TestCRAFFTScore_parity(t *testing.T) {
	const T, F = true, false
	tests := []struct {
		in   [6]bool
		want int
	}{
		{[6]bool{F, F, F, F, F, F}, 0}, // source case 1 "all no"
		{[6]bool{T, T, T, T, T, T}, 6}, // source case 2 "no part A, only CAR counts": source 1, primary 6
		{[6]bool{F, F, F, F, F, F}, 0}, // source case 3 "part A only"
		{[6]bool{F, T, F, F, F, F}, 1}, // source case 4 "medium"
		{[6]bool{T, T, F, F, F, F}, 2}, // source case 5 "high edge"
		{[6]bool{T, T, T, T, T, T}, 6}, // source case 6 "max"
	}
	for i, tc := range tests {
		if got := CRAFFTScore(tc.in); got != tc.want {
			t.Errorf("source case %d: CRAFFTScore = %d, want %d", i+1, got, tc.want)
		}
	}
}

func TestCRAFFTScore_noAlloc(t *testing.T) {
	in := [6]bool{true, true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = CRAFFTScore(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCRAFFTScore(b *testing.B) {
	b.ReportAllocs()
	in := [6]bool{true, true}
	for b.Loop() {
		sinkInt = CRAFFTScore(in)
	}
}
