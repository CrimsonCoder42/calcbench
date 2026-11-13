package scales

import "testing"

// HAND-DERIVED from Apgar 1953, pp. 261–262: five signs, each rated 0, 1 or 2, summed; "A score of ten
// indicated a baby in the best possible condition." All 2 → 10; all 0 → 0; 2+1+1+1+0 → 5; 2+2+2+2+1 → 9.
func TestApgarScore_example(t *testing.T) {
	tests := []struct{ hr, resp, reflex, tone, color, want int }{
		{2, 2, 2, 2, 2, 10},
		{0, 0, 0, 0, 0, 0},
		{2, 1, 1, 1, 0, 5},
		{2, 2, 2, 2, 1, 9},
	}
	for _, tc := range tests {
		if got, err := ApgarScore(tc.hr, tc.resp, tc.reflex, tc.tone, tc.color); err != nil || got != tc.want {
			t.Errorf("ApgarScore(%d, %d, %d, %d, %d) = (%v, %v), want %d", tc.hr, tc.resp, tc.reflex, tc.tone, tc.color, got, err, tc.want)
		}
	}
}

// Each sign on both sides of 0–2.
func TestApgarScore_errors(t *testing.T) {
	call := func(hr, resp, reflex, tone, color int) func() error {
		return func() error { _, err := ApgarScore(hr, resp, reflex, tone, color); return err }
	}
	runErrCases(t, []errCase{
		{"heartRate -1", call(-1, 1, 1, 1, 1), ErrOutOfRange},
		{"heartRate 3", call(3, 1, 1, 1, 1), ErrOutOfRange},
		{"respiratoryEffort -1", call(1, -1, 1, 1, 1), ErrOutOfRange},
		{"respiratoryEffort 3", call(1, 3, 1, 1, 1), ErrOutOfRange},
		{"reflexIrritability -1", call(1, 1, -1, 1, 1), ErrOutOfRange},
		{"reflexIrritability 3", call(1, 1, 3, 1, 1), ErrOutOfRange},
		{"muscleTone -1", call(1, 1, 1, -1, 1), ErrOutOfRange},
		{"muscleTone 3", call(1, 1, 1, 3, 1), ErrOutOfRange},
		{"color -1", call(1, 1, 1, 1, -1), ErrOutOfRange},
		{"color 3", call(1, 1, 1, 1, 3), ErrOutOfRange},
	})
}

func TestApgarScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ApgarScore(2, 1, 1, 1, 0) }); n != 0 {
		t.Errorf("ApgarScore allocates %v times, want 0", n)
	}
}

func BenchmarkApgarScore(b *testing.B) {
	b.ReportAllocs()
	hr, resp, reflex, tone, color := 2, 1, 1, 1, 0
	for b.Loop() {
		sinkInt, _ = ApgarScore(hr, resp, reflex, tone, color)
	}
}
