package psychiatry

import (
	"errors"
	"testing"
)

// Wesson & Ling 2003 Appendix 1: printed item values; sum of 11; bands 5–12, 13–24, 25–36, > 36.
func TestCOWSScore_example(t *testing.T) {
	maxIt := COWSItems{4, 4, 5, 5, 4, 4, 5, 4, 4, 4, 5}
	if got, err := COWSScore(maxIt); err != nil || got != 48 {
		t.Errorf("max = (%d, %v), want 48", got, err)
	}
	if got, err := COWSScore(COWSItems{}); err != nil || got != 0 {
		t.Errorf("zero = (%d, %v), want 0", got, err)
	}
	ex := COWSItems{2, 1, 1, 1, 1, 2, 2, 1, 1, 1, 0}
	if got, _ := COWSScore(ex); got != 13 {
		t.Errorf("example = %d, want 13", got)
	}
}

func TestCOWSScore_errors(t *testing.T) {
	bad := []COWSItems{
		{RestingPulse: 3}, {Sweating: 5}, {Restlessness: 2}, {PupilSize: 3}, {BoneJointAches: 3},
		{RunnyNoseTearing: 3}, {GIUpset: 4}, {Tremor: 3}, {Yawning: 3}, {AnxietyIrritability: 3},
		{GoosefleshSkin: 1}, {GoosefleshSkin: 6}, {RestingPulse: -1},
	}
	for _, b := range bad {
		if _, err := COWSScore(b); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("COWSScore(%+v): %v, want ErrOutOfRange", b, err)
		}
	}
}

func TestCOWSSeverityOf_example(t *testing.T) {
	cases := []struct {
		total int
		want  COWSSeverity
	}{
		{0, COWSNoBand}, {4, COWSNoBand}, {5, COWSMild}, {12, COWSMild}, {13, COWSModerate}, {24, COWSModerate},
		{25, COWSModeratelySevere}, {36, COWSModeratelySevere}, {37, COWSSevere}, {48, COWSSevere},
	}
	for _, c := range cases {
		if got, err := COWSSeverityOf(c.total); err != nil || got != c.want {
			t.Errorf("COWSSeverityOf(%d) = (%q, %v), want %q", c.total, got, err, c.want)
		}
	}
	for _, bad := range []int{-1, 49} {
		if _, err := COWSSeverityOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("COWSSeverityOf(%d): %v, want ErrOutOfRange", bad, err)
		}
	}
}

// parity: ledger cows_score, source cases 1–2 (consumer maps option index to the printed points, then sums;
// its bands use > 36 for severe).
func TestCOWSScore_parity(t *testing.T) {
	it := COWSItems{4, 3, 3, 2, 2, 2, 3, 2, 2, 2, 3} // consumer option indices 3,3,2,2,2,2,3,2,2,2,1
	if got, _ := COWSScore(it); got != 28 {          // parity: ledger cows_score, source case 1
		t.Errorf("case 1 = %d, want 28", got)
	}
	if got, _ := COWSSeverityOf(37); got != COWSSevere { // parity: ledger cows_score, source case 2 (> 36)
		t.Errorf("case 2 = %q", got)
	}
}

func TestCOWS_noAlloc(t *testing.T) {
	it := COWSItems{2, 1, 1, 1, 1, 2, 2, 1, 1, 1, 0}
	if n := testing.AllocsPerRun(100, func() {
		s, _ := COWSScore(it)
		c, _ := COWSSeverityOf(s)
		sinkString = string(c)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCOWSScore(b *testing.B) {
	b.ReportAllocs()
	it := COWSItems{2, 1, 1, 1, 1, 2, 2, 1, 1, 1, 0}
	for b.Loop() {
		sinkInt, _ = COWSScore(it)
	}
}

func BenchmarkCOWSSeverityOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		c, _ := COWSSeverityOf(13)
		sinkString = string(c)
	}
}
