package neuro

import (
	"errors"
	"testing"
)

// Lipton 2015: items 0–2, sum 0–8; categories 0, 1–5, 6–7, 8 with each edge on both sides.
func TestMTOQ4Score_example(t *testing.T) {
	cases := []struct {
		in   [4]int
		want int
	}{{[4]int{2, 2, 2, 2}, 8}, {[4]int{}, 0}, {[4]int{1, 1, 2, 2}, 6}}
	for _, c := range cases {
		if got, err := MTOQ4Score(c.in); err != nil || got != c.want {
			t.Errorf("MTOQ4Score(%v) = (%d, %v), want %d", c.in, got, err, c.want)
		}
	}
	for _, bad := range [][4]int{{3, 0, 0, 0}, {0, -1, 0, 0}} {
		if _, err := MTOQ4Score(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("MTOQ4Score(%v): %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestMTOQ4CategoryOf_example(t *testing.T) {
	cases := []struct {
		total int
		want  MTOQ4Category
	}{{0, MTOQ4VeryPoor}, {1, MTOQ4Poor}, {5, MTOQ4Poor}, {6, MTOQ4Moderate}, {7, MTOQ4Moderate}, {8, MTOQ4Maximum}}
	for _, c := range cases {
		if got, err := MTOQ4CategoryOf(c.total); err != nil || got != c.want {
			t.Errorf("MTOQ4CategoryOf(%d) = (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 9} {
		if _, err := MTOQ4CategoryOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("MTOQ4CategoryOf(%d): %v, want ErrOutOfRange", bad, err)
		}
	}
}

// parity: ledger mtoq4, source cases 1–6 (source mtoq4_test.go: zero, poor edge, poor top, moderate edge,
// moderate top, maximum; the consumer's capitalised labels are the same words).
func TestMTOQ4CategoryOf_parity(t *testing.T) {
	cases := []struct {
		in    [4]int
		total int
		want  MTOQ4Category
	}{
		{[4]int{0, 0, 0, 0}, 0, MTOQ4VeryPoor}, // parity: ledger mtoq4, source case 1
		{[4]int{1, 0, 0, 0}, 1, MTOQ4Poor},     // parity: ledger mtoq4, source case 2
		{[4]int{2, 1, 1, 1}, 5, MTOQ4Poor},     // parity: ledger mtoq4, source case 3
		{[4]int{2, 2, 1, 1}, 6, MTOQ4Moderate}, // parity: ledger mtoq4, source case 4
		{[4]int{2, 2, 2, 1}, 7, MTOQ4Moderate}, // parity: ledger mtoq4, source case 5
		{[4]int{2, 2, 2, 2}, 8, MTOQ4Maximum},  // parity: ledger mtoq4, source case 6
	}
	for i, c := range cases {
		s, err := MTOQ4Score(c.in)
		if err != nil || s != c.total {
			t.Errorf("case %d score = (%d, %v), want %d", i+1, s, err, c.total)
		}
		if got, _ := MTOQ4CategoryOf(s); got != c.want {
			t.Errorf("case %d = %q, want %q", i+1, got, c.want)
		}
	}
}

func TestMTOQ4_noAlloc(t *testing.T) {
	in := [4]int{1, 1, 2, 2}
	if n := testing.AllocsPerRun(100, func() {
		s, _ := MTOQ4Score(in)
		c, _ := MTOQ4CategoryOf(s)
		sinkString = string(c)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMTOQ4Score(b *testing.B) {
	b.ReportAllocs()
	in := [4]int{1, 1, 2, 2}
	for b.Loop() {
		sinkInt, _ = MTOQ4Score(in)
	}
}

func BenchmarkMTOQ4CategoryOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		c, _ := MTOQ4CategoryOf(6)
		sinkString = string(c)
	}
}
