package emergency

import "testing"

var sinkInt int

// Barbosa 2013 abstract: swelling 2, hard testicle 2, absent cremasteric reflex 1, nausea/vomiting 1, high
// riding testis 1.
func TestTWISTScore_example(t *testing.T) {
	tests := []struct {
		sw, hard, crem, nv, high bool
		want                     int
	}{
		{false, false, false, false, false, 0}, {true, false, false, false, false, 2}, {false, true, false, false, false, 2},
		{false, false, true, false, false, 1}, {false, false, false, true, false, 1}, {false, false, false, false, true, 1},
		{true, true, false, false, false, 4}, {true, true, true, false, false, 5}, {true, true, true, true, true, 7},
	}
	for _, tc := range tests {
		if got := TWISTScore(tc.sw, tc.hard, tc.crem, tc.nv, tc.high); got != tc.want {
			t.Errorf("TWISTScore(%+v) = %d, want %d", tc, got, tc.want)
		}
	}
}

func TestTWISTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = TWISTScore(true, true, false, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkTWISTScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = TWISTScore(x, x, !x, x, !x)
	}
}
