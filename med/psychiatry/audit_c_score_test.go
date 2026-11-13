package psychiatry

import (
	"errors"
	"testing"
)

var sinkBool bool

// Flentje 2020: three items 0–4, total 0–12.
func TestAUDITCScore_example(t *testing.T) {
	cases := []struct{ f, q, h, want int }{{0, 0, 0, 0}, {4, 4, 4, 12}, {2, 1, 1, 4}, {1, 1, 1, 3}}
	for _, c := range cases {
		if got, err := AUDITCScore(c.f, c.q, c.h); err != nil || got != c.want {
			t.Errorf("AUDITCScore(%d, %d, %d) = (%d, %v), want %d", c.f, c.q, c.h, got, err, c.want)
		}
	}
	for _, b := range [][3]int{{-1, 0, 0}, {5, 0, 0}, {0, 5, 0}, {0, 0, -1}} {
		if _, err := AUDITCScore(b[0], b[1], b[2]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("AUDITCScore%v: %v, want ErrOutOfRange", b, err)
		}
	}
}

// parity: ledger audit_c, source cases 1–6 (source audit_c_test.go "score"; the sex input only moves the
// consumer's band and is not an input here).
func TestAUDITCScore_parity(t *testing.T) {
	cases := []struct{ f, q, h, want int }{
		{2, 1, 1, 4},  // parity: ledger audit_c, source case 1 (worked example male 4)
		{1, 1, 1, 3},  // parity: ledger audit_c, source case 2 (female 3)
		{1, 1, 1, 3},  // parity: ledger audit_c, source case 3 (male 3)
		{4, 2, 2, 8},  // parity: ledger audit_c, source case 4 (severe 8)
		{0, 0, 0, 0},  // parity: ledger audit_c, source case 5 (zero)
		{4, 4, 4, 12}, // parity: ledger audit_c, source case 6 (max)
	}
	for i, c := range cases {
		if got, err := AUDITCScore(c.f, c.q, c.h); err != nil || got != c.want {
			t.Errorf("case %d = (%d, %v), want %d", i+1, got, err, c.want)
		}
	}
}

func TestAUDITCScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = AUDITCScore(2, 1, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkAUDITCScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = AUDITCScore(2, 1, 1)
	}
}
