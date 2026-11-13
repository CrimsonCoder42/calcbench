package psychiatry

import (
	"errors"
	"testing"
)

// Manual p. 3/p. 6: seven items scored 0–3, summed to 0–21. HAND-DERIVED: all 0 → 0; all 3 → 21.
func TestGAD7Score_example(t *testing.T) {
	for _, tc := range []struct {
		in   [7]int
		want int
	}{{[7]int{}, 0}, {[7]int{3, 3, 3, 3, 3, 3, 3}, 21}, {[7]int{3, 2, 3, 1, 0, 1, 2}, 12}} {
		if got, err := GAD7Score(tc.in); err != nil || got != tc.want {
			t.Errorf("GAD7Score(%v) = (%d, %v), want %d", tc.in, got, err, tc.want)
		}
	}
}

// parity: ledger gad7, source cases 1–7 (the totals; the bands are checked in TestGAD7SeverityOf_parity).
func TestGAD7Score_parity(t *testing.T) {
	tests := []struct {
		in   [7]int
		want int
	}{
		{[7]int{0, 0, 0, 0, 0, 0, 0}, 0},  // source case 1 "zero"
		{[7]int{1, 1, 1, 1, 0, 0, 0}, 4},  // source case 2 "minimal top"
		{[7]int{1, 1, 1, 1, 1, 0, 0}, 5},  // source case 3 "mild edge"
		{[7]int{2, 2, 2, 2, 1, 0, 0}, 9},  // source case 4 "mild top"
		{[7]int{3, 2, 3, 1, 0, 1, 2}, 12}, // source case 5 "worked example"
		{[7]int{3, 2, 2, 2, 2, 2, 2}, 15}, // source case 6 "severe edge"
		{[7]int{3, 3, 3, 3, 3, 3, 3}, 21}, // source case 7 "max"
	}
	for i, tc := range tests {
		if got, err := GAD7Score(tc.in); err != nil || got != tc.want {
			t.Errorf("source case %d: GAD7Score = (%d, %v), want %d", i+1, got, err, tc.want)
		}
	}
}

// The source's refusals (an item of 4) map to ErrOutOfRange; a 1.5 cannot be passed to an int.
func TestGAD7Score_errors(t *testing.T) {
	for i := range 7 {
		for _, v := range []int{-1, 4} {
			var in [7]int
			in[i] = v
			if _, err := GAD7Score(in); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("items[%d] = %d: error = %v, want ErrOutOfRange", i, v, err)
			}
		}
	}
}

func TestGAD7Score_noAlloc(t *testing.T) {
	in := [7]int{3, 2, 3, 1, 0, 1, 2}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GAD7Score(in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGAD7Score(b *testing.B) {
	b.ReportAllocs()
	in := [7]int{3, 2, 3, 1, 0, 1, 2}
	for b.Loop() {
		sinkInt, _ = GAD7Score(in)
	}
}
