package pediatrics

import (
	"errors"
	"testing"
)

// Levine 2015 Table 3: appearance 0/2/4, respirations 0/2, skin pinch 0/2/4, tears 0/1/2.
func TestDHAKADehydrationScore_example(t *testing.T) {
	tests := []struct{ a, r, s, t, want int }{
		{0, 0, 0, 0, 0}, {1, 0, 0, 0, 2}, {2, 0, 0, 0, 4}, {0, 1, 0, 0, 2}, {0, 0, 1, 0, 2}, {0, 0, 2, 0, 4},
		{0, 0, 0, 1, 1}, {0, 0, 0, 2, 2}, {1, 0, 1, 1, 5}, {2, 1, 2, 2, 12},
	}
	for _, tc := range tests {
		if got, err := DHAKADehydrationScore(tc.a, tc.r, tc.s, tc.t); err != nil || got != tc.want {
			t.Errorf("DHAKADehydrationScore(%d, %d, %d, %d) = (%d, %v), want %d", tc.a, tc.r, tc.s, tc.t, got, err, tc.want)
		}
	}
}

func TestDHAKADehydrationScore_errors(t *testing.T) {
	for _, c := range [][4]int{{-1, 0, 0, 0}, {3, 0, 0, 0}, {0, -1, 0, 0}, {0, 2, 0, 0}, {0, 0, -1, 0}, {0, 0, 3, 0}, {0, 0, 0, -1}, {0, 0, 0, 3}} {
		if _, err := DHAKADehydrationScore(c[0], c[1], c[2], c[3]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("DHAKADehydrationScore(%v) error = %v, want ErrOutOfRange", c, err)
		}
	}
}

func TestDHAKADehydrationScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = DHAKADehydrationScore(1, 0, 1, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkDHAKADehydrationScore(b *testing.B) {
	b.ReportAllocs()
	a, r, s, tr := 1, 0, 1, 1
	for b.Loop() {
		sinkInt, _ = DHAKADehydrationScore(a, r, s, tr)
	}
}
