package neuro

import (
	"errors"
	"testing"
)

// Fischer 2010 Table 1: four components 0–4, maximum 16.
func TestFOURScore_example(t *testing.T) {
	cases := []struct{ e, m, b, r, want int }{{4, 4, 4, 4, 16}, {0, 0, 0, 0, 0}, {3, 2, 4, 1, 10}}
	for _, c := range cases {
		if got, err := FOURScore(c.e, c.m, c.b, c.r); err != nil || got != c.want {
			t.Errorf("FOURScore(%d, %d, %d, %d) = (%d, %v), want %d", c.e, c.m, c.b, c.r, got, err, c.want)
		}
	}
}

func TestFOURScore_errors(t *testing.T) {
	bad := [][4]int{{-1, 0, 0, 0}, {5, 0, 0, 0}, {0, 5, 0, 0}, {0, 0, -1, 0}, {0, 0, 0, 5}}
	for _, b := range bad {
		if _, err := FOURScore(b[0], b[1], b[2], b[3]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("FOURScore%v: %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger four_score_coma, source case 1 (consumer sums E + M + B + R).
func TestFOURScore_parity(t *testing.T) {
	if got, _ := FOURScore(1, 2, 3, 0); got != 6 { // parity: ledger four_score_coma, source case 1
		t.Errorf("case 1 = %d, want 6", got)
	}
}

func TestFOURScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = FOURScore(4, 4, 4, 4) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkFOURScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = FOURScore(4, 4, 4, 4)
	}
}
