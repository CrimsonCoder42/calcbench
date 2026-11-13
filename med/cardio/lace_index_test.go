package cardio

import "testing"

// van Walraven 2010, Table 3: every L band edge, the C and E caps, maximum 19.
func TestLACEIndex_points(t *testing.T) {
	tests := []struct {
		los   int
		acute bool
		c, e  int
		want  int
	}{
		{0, false, 0, 0, 0},
		{1, false, 0, 0, 1},
		{3, false, 0, 0, 3},
		{4, false, 0, 0, 4},
		{6, false, 0, 0, 4},
		{7, false, 0, 0, 5},
		{13, false, 0, 0, 5},
		{14, false, 0, 0, 7},
		{5, true, 2, 1, 10}, // example
		{0, false, 3, 0, 3},
		{0, false, 4, 0, 5},
		{0, false, 0, 3, 3},
		{0, false, 0, 9, 4},
		{30, true, 10, 10, 19},
	}
	for _, tc := range tests {
		if got, err := LACEIndex(tc.los, tc.acute, tc.c, tc.e); err != nil || got != tc.want {
			t.Errorf("LACEIndex(%d, %v, %d, %d) = (%d, %v), want %d", tc.los, tc.acute, tc.c, tc.e, got, err, tc.want)
		}
	}
	for _, in := range [][3]int{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}} {
		_, err := LACEIndex(in[0], false, in[1], in[2])
		checkErr(t, "negative", err, ErrOutOfRange)
	}
}

func TestLACEIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = LACEIndex(5, true, 2, 1) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkLACEIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = LACEIndex(5, true, 2, 1)
	}
}
