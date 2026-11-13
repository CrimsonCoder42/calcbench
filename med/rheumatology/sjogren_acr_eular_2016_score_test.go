package rheumatology

import "testing"

var (
	sinkInt  int
	sinkBool bool
)

// Shiboski 2017 Table 3 weights: biopsy 3, anti-SSA 3, ocular staining 1, Schirmer 1, saliva flow 1; the
// classification threshold "Score ≥ 4" sits between 3 and 4.
func TestSjogrenACREULAR2016Score_example(t *testing.T) {
	tests := []struct {
		b, a, o, s, f bool
		want          int
	}{
		{false, false, false, false, false, 0},
		{true, false, false, false, false, 3},
		{false, true, false, false, false, 3},
		{false, false, true, false, false, 1},
		{false, false, false, true, false, 1},
		{false, false, false, false, true, 1},
		{true, false, false, true, false, 4},
		{false, false, true, true, true, 3},
		{true, true, true, true, true, 9},
	}
	for _, tc := range tests {
		if got := SjogrenACREULAR2016Score(tc.b, tc.a, tc.o, tc.s, tc.f); got != tc.want {
			t.Errorf("SjogrenACREULAR2016Score(%v, %v, %v, %v, %v) = %d, want %d", tc.b, tc.a, tc.o, tc.s, tc.f, got, tc.want)
		}
	}
}

func TestSjogrenACREULAR2016Score_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = SjogrenACREULAR2016Score(true, false, true, false, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkSjogrenACREULAR2016Score(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = SjogrenACREULAR2016Score(x, !x, x, !x, x)
	}
}
