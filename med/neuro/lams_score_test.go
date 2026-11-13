package neuro

import (
	"errors"
	"testing"
)

// Arthur 2023 / Thu 2025: facial 0–1, arm 0–2, grip 0–2, total 0–5.
func TestLAMSScore_example(t *testing.T) {
	cases := []struct{ f, a, g, want int }{{0, 0, 0, 0}, {1, 2, 2, 5}, {1, 1, 0, 2}, {0, 2, 1, 3}}
	for _, c := range cases {
		if got, err := LAMSScore(c.f, c.a, c.g); err != nil || got != c.want {
			t.Errorf("LAMSScore(%d, %d, %d) = (%d, %v), want %d", c.f, c.a, c.g, got, err, c.want)
		}
	}
}

func TestLAMSScore_errors(t *testing.T) {
	bad := [][3]int{{-1, 0, 0}, {2, 0, 0}, {0, -1, 0}, {0, 3, 0}, {0, 0, -1}, {0, 0, 3}}
	for _, b := range bad {
		if _, err := LAMSScore(b[0], b[1], b[2]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("LAMSScore%v: %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger lams, source case 1 (consumer sums the three items).
func TestLAMSScore_parity(t *testing.T) {
	if got, _ := LAMSScore(1, 2, 1); got != 4 { // parity: ledger lams, source case 1
		t.Errorf("case 1 = %d, want 4", got)
	}
}

func TestLAMSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = LAMSScore(1, 1, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkLAMSScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = LAMSScore(1, 1, 1)
	}
}
