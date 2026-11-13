package scales

import "testing"

// Apgar 1953, p. 267: "groups 8, 9, and 10 indicate infants in good condition, 0, 1, and 2, poor condition,
// and the remaining scores, fair condition." Every edge on both sides; the words are written out verbatim here.
func TestApgarCondition1953_edges(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{0, "poor condition"},
		{2, "poor condition"},
		{3, "fair condition"},
		{7, "fair condition"},
		{8, "good condition"},
		{10, "good condition"},
	}
	for _, tc := range tests {
		if got, err := ApgarCondition1953(tc.score); err != nil || string(got) != tc.want {
			t.Errorf("ApgarCondition1953(%d) = (%q, %v), want %q", tc.score, got, err, tc.want)
		}
	}
}

func TestApgarCondition1953_errors(t *testing.T) {
	call := func(s int) func() error {
		return func() error { _, err := ApgarCondition1953(s); return err }
	}
	runErrCases(t, []errCase{
		{"-1", call(-1), ErrOutOfRange},
		{"11", call(11), ErrOutOfRange},
		{"max int", call(1<<63 - 1), ErrOutOfRange},
	})
}

func TestApgarCondition1953_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkCondition, _ = ApgarCondition1953(9) }); n != 0 {
		t.Errorf("ApgarCondition1953 allocates %v times, want 0", n)
	}
}

func BenchmarkApgarCondition1953(b *testing.B) {
	b.ReportAllocs()
	s := 9
	for b.Loop() {
		sinkCondition, _ = ApgarCondition1953(s)
	}
}
