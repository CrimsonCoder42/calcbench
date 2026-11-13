package scales

import "testing"

// Rimel et al. 1982 abstract: "severe, 3-8; moderate, 9-12; and minor, 13-15". Every edge on both sides; the
// words are written out verbatim here.
func TestGlasgowComaSeverityRimel_edges(t *testing.T) {
	tests := []struct {
		total int
		want  string
	}{
		{3, "severe"},
		{8, "severe"},
		{9, "moderate"},
		{12, "moderate"},
		{13, "minor"},
		{15, "minor"},
	}
	for _, tc := range tests {
		if got, err := GlasgowComaSeverityRimel(tc.total); err != nil || string(got) != tc.want {
			t.Errorf("GlasgowComaSeverityRimel(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
}

func TestGlasgowComaSeverityRimel_errors(t *testing.T) {
	call := func(total int) func() error {
		return func() error { _, err := GlasgowComaSeverityRimel(total); return err }
	}
	runErrCases(t, []errCase{
		{"2", call(2), ErrOutOfRange},
		{"16", call(16), ErrOutOfRange},
		{"0", call(0), ErrOutOfRange},
		{"min int", call(-1 << 63), ErrOutOfRange},
	})
}

func TestGlasgowComaSeverityRimel_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkSeverity, _ = GlasgowComaSeverityRimel(10) }); n != 0 {
		t.Errorf("GlasgowComaSeverityRimel allocates %v times, want 0", n)
	}
}

func BenchmarkGlasgowComaSeverityRimel(b *testing.B) {
	b.ReportAllocs()
	total := 10
	for b.Loop() {
		sinkSeverity, _ = GlasgowComaSeverityRimel(total)
	}
}
