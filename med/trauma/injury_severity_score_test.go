package trauma

import "testing"

// Clark 2018 / Deng 2016: sum of squares of the three largest region maxima; AIS 6 → 75. HAND-DERIVED.
func TestInjurySeverityScore_rule(t *testing.T) {
	tests := []struct {
		r    [6]int
		want int
	}{
		{[6]int{0, 0, 0, 0, 0, 0}, 0},
		{[6]int{1, 0, 0, 0, 0, 0}, 1},
		{[6]int{4, 0, 3, 0, 2, 0}, 29},
		{[6]int{2, 2, 2, 2, 2, 2}, 12},
		{[6]int{5, 5, 5, 5, 5, 5}, 75},
		{[6]int{1, 2, 3, 4, 5, 1}, 50},
		{[6]int{5, 1, 4, 3, 2, 1}, 50},
		{[6]int{0, 0, 0, 0, 0, 6}, 75},
		{[6]int{6, 1, 1, 1, 1, 1}, 75},
		{[6]int{3, 3, 1, 1, 0, 0}, 19},
	}
	for _, tc := range tests {
		r := tc.r
		if got, err := InjurySeverityScore(r[0], r[1], r[2], r[3], r[4], r[5]); err != nil || got != tc.want {
			t.Errorf("%v: (%d, %v), want %d", r, got, err, tc.want)
		}
	}
}

// Oracle: every combination of six AIS values 0–5 (46 656 cases) against sorting the six values.
func TestInjurySeverityScore_oracle(t *testing.T) {
	var r [6]int
	for n := range 46656 {
		m := n
		for i := range r {
			r[i] = m % 6
			m /= 6
		}
		s := r
		for i := 1; i < 6; i++ {
			for j := i; j > 0 && s[j] > s[j-1]; j-- {
				s[j], s[j-1] = s[j-1], s[j]
			}
		}
		want := s[0]*s[0] + s[1]*s[1] + s[2]*s[2]
		if got, _ := InjurySeverityScore(r[0], r[1], r[2], r[3], r[4], r[5]); got != want {
			t.Fatalf("%v: %d, want %d", r, got, want)
		}
	}
}

// parity: ledger injury_severity_score, source example: head 4, chest 3, extremity 2 → 29; any 6 → 75. The app's
// bands (< 9, 9–15, 16–24, ≥ 25) are not in either secondary and are not ported.
func TestInjurySeverityScore_parity(t *testing.T) {
	if got, _ := InjurySeverityScore(4, 0, 3, 0, 2, 0); got != 29 {
		t.Errorf("got %d, want 29", got)
	}
	if got, _ := InjurySeverityScore(1, 0, 6, 0, 0, 0); got != 75 {
		t.Errorf("got %d, want 75", got)
	}
}

func TestInjurySeverityScore_errors(t *testing.T) {
	for i := range 6 {
		for _, bad := range []int{-1, 7} {
			var r [6]int
			r[i] = bad
			_, err := InjurySeverityScore(r[0], r[1], r[2], r[3], r[4], r[5])
			checkErr(t, "region", err, ErrOutOfRange)
		}
	}
}

func TestInjurySeverityScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = InjurySeverityScore(4, 0, 3, 0, 2, 0) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkInjurySeverityScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = InjurySeverityScore(4, 0, 3, 0, 2, 0)
	}
}
