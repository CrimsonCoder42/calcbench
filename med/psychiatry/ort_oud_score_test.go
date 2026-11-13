package psychiatry

import "testing"

// Cheatle 2019: nine unweighted items, yes = 1. HAND-DERIVED: none → 0; all → 9.
func TestORTOUDScore_example(t *testing.T) {
	if got := ORTOUDScore([9]bool{}); got != 0 {
		t.Errorf("none = %d, want 0", got)
	}
	if got := ORTOUDScore([9]bool{true, true, true, true, true, true, true, true, true}); got != 9 {
		t.Errorf("all = %d, want 9", got)
	}
}

// parity: ledger ort_oud, source cases 1–6 (the total; the source's subtotals are not ported, ledger note).
// Items in source order: age_16_45, hx_depression, hx_psych, personal_alcohol, personal_illegal, personal_rx,
// family_alcohol, family_illegal, family_rx.
func TestORTOUDScore_parity(t *testing.T) {
	const T, F = true, false
	tests := []struct {
		in   [9]bool
		want int
	}{
		{[9]bool{F, F, F, F, F, F, F, F, F}, 0}, // source case 1 "zero"
		{[9]bool{T, T, F, F, F, F, F, F, F}, 2}, // source case 2 "low top"
		{[9]bool{T, T, T, F, F, F, F, F, F}, 3}, // source case 3 "high edge (age + psych)"
		{[9]bool{F, F, F, T, T, T, T, F, F}, 4}, // source case 4 "high, substance"
		{[9]bool{F, T, F, F, F, T, F, T, F}, 3}, // source case 5 "high, mixed"
		{[9]bool{T, T, T, T, T, T, T, T, T}, 9}, // source case 6 "max"
	}
	for i, tc := range tests {
		if got := ORTOUDScore(tc.in); got != tc.want {
			t.Errorf("source case %d: ORTOUDScore = %d, want %d", i+1, got, tc.want)
		}
	}
}

func TestORTOUDScore_noAlloc(t *testing.T) {
	in := [9]bool{true, false, true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = ORTOUDScore(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkORTOUDScore(b *testing.B) {
	b.ReportAllocs()
	in := [9]bool{true, false, true}
	for b.Loop() {
		sinkInt = ORTOUDScore(in)
	}
}
