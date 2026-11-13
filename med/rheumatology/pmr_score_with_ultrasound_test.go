package rheumatology

import "testing"

// Dasgupta 2012 Table 6, with US: clinical items as without US, plus 1 + 1 for the two US findings; threshold 5.
func TestPMRScoreWithUltrasound_example(t *testing.T) {
	tests := []struct {
		s, h, r, o, us1, us2 bool
		want                 int
	}{
		{false, false, false, false, false, false, 0},
		{false, false, false, false, true, false, 1},
		{false, false, false, false, false, true, 1},
		{true, false, true, false, false, false, 4},
		{true, false, true, false, false, true, 5},
		{true, true, true, true, true, true, 8},
	}
	for _, tc := range tests {
		if got := PMRScoreWithUltrasound(tc.s, tc.h, tc.r, tc.o, tc.us1, tc.us2); got != tc.want {
			t.Errorf("PMRScoreWithUltrasound(%+v) = %d, want %d", tc, got, tc.want)
		}
	}
}

func TestPMRScoreWithUltrasound_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = PMRScoreWithUltrasound(true, false, true, false, true, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPMRScoreWithUltrasound(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = PMRScoreWithUltrasound(x, !x, x, !x, x, x)
	}
}
