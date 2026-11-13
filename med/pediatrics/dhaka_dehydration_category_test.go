package pediatrics

import (
	"errors"
	"testing"
)

// Levine 2015 Table 3: 0–1 no, 2–3 some, ≥ 4 severe; each edge on both sides.
func TestDHAKADehydrationCategory_example(t *testing.T) {
	tests := []struct {
		s    int
		want DHAKACategory
	}{
		{0, DHAKANoDehydration}, {1, DHAKANoDehydration}, {2, DHAKASomeDehydration}, {3, DHAKASomeDehydration},
		{4, DHAKASevereDehydration}, {12, DHAKASevereDehydration},
	}
	for _, tc := range tests {
		if got, err := DHAKADehydrationCategory(tc.s); err != nil || got != tc.want {
			t.Errorf("DHAKADehydrationCategory(%d) = (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 13} {
		if _, err := DHAKADehydrationCategory(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("DHAKADehydrationCategory(%d) error = %v, want ErrOutOfRange", s, err)
		}
	}
}

func TestDHAKADehydrationCategory_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { c, _ := DHAKADehydrationCategory(3); sinkString = string(c) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkDHAKADehydrationCategory(b *testing.B) {
	b.ReportAllocs()
	s := 3
	for b.Loop() {
		c, _ := DHAKADehydrationCategory(s)
		sinkString = string(c)
	}
}
