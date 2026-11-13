package rheumatology

import "testing"

// Dasgupta 2012 Table 6, without US: stiffness 2, hip 1, RF/ACPA absent 2, other joints absent 1; threshold 4.
func TestPMRScoreWithoutUltrasound_example(t *testing.T) {
	tests := []struct {
		s, h, r, o bool
		want       int
	}{
		{false, false, false, false, 0}, {true, false, false, false, 2}, {false, true, false, false, 1},
		{false, false, true, false, 2}, {false, false, false, true, 1}, {true, true, false, false, 3},
		{true, false, true, false, 4}, {true, true, true, true, 6},
	}
	for _, tc := range tests {
		if got := PMRScoreWithoutUltrasound(tc.s, tc.h, tc.r, tc.o); got != tc.want {
			t.Errorf("PMRScoreWithoutUltrasound(%v, %v, %v, %v) = %d, want %d", tc.s, tc.h, tc.r, tc.o, got, tc.want)
		}
	}
}

func TestPMRScoreWithoutUltrasound_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = PMRScoreWithoutUltrasound(true, false, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPMRScoreWithoutUltrasound(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = PMRScoreWithoutUltrasound(x, !x, x, !x)
	}
}
