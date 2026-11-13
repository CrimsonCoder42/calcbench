package scales

import "testing"

// Bands from the official ESS site, About the ESS (read 2026-11-12): 0–5, 6–10, 11–12, 13–15, 16–24. Every
// edge on both sides. Labels are written out verbatim here, so a typo in a constant fails the test.
func TestEpworthCategory_edges(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{0, "Lower Normal Daytime Sleepiness"},
		{5, "Lower Normal Daytime Sleepiness"},
		{6, "Higher Normal Daytime Sleepiness"},
		{10, "Higher Normal Daytime Sleepiness"},
		{11, "Mild Excessive Daytime Sleepiness"},
		{12, "Mild Excessive Daytime Sleepiness"},
		{13, "Moderate Excessive Daytime Sleepiness"},
		{15, "Moderate Excessive Daytime Sleepiness"},
		{16, "Severe Excessive Daytime Sleepiness"},
		{24, "Severe Excessive Daytime Sleepiness"},
	}
	for _, tc := range tests {
		if got, err := EpworthCategory(tc.score); err != nil || string(got) != tc.want {
			t.Errorf("EpworthCategory(%d) = (%q, %v), want %q", tc.score, got, err, tc.want)
		}
	}
}

func TestEpworthCategory_errors(t *testing.T) {
	call := func(s int) func() error {
		return func() error { _, err := EpworthCategory(s); return err }
	}
	runErrCases(t, []errCase{
		{"-1", call(-1), ErrOutOfRange},
		{"25", call(25), ErrOutOfRange},
		{"min int", call(-1 << 63), ErrOutOfRange},
		{"max int", call(1<<63 - 1), ErrOutOfRange},
	})
}

func TestEpworthCategory_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkEpworth, _ = EpworthCategory(12) }); n != 0 {
		t.Errorf("EpworthCategory allocates %v times, want 0", n)
	}
}

func BenchmarkEpworthCategory(b *testing.B) {
	b.ReportAllocs()
	s := 12
	for b.Loop() {
		sinkEpworth, _ = EpworthCategory(s)
	}
}
