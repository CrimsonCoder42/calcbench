package neuro

import (
	"errors"
	"math"
	"testing"
)

// Mower 2005 abstract: low risk = under 65 and none of seven findings; age edge on both sides, HAND-DERIVED.
func TestNEXUSHeadCTLowRisk_example(t *testing.T) {
	cases := []struct {
		age  float64
		f    NEXUSHeadCTFindings
		want bool
	}{
		{64.9, NEXUSHeadCTFindings{}, true},
		{65, NEXUSHeadCTFindings{}, false},
		{0, NEXUSHeadCTFindings{}, true},
		{30, NEXUSHeadCTFindings{ScalpHematoma: true}, false},
		{30, NEXUSHeadCTFindings{Coagulopathy: true}, false},
		{30, NEXUSHeadCTFindings{PersistentVomiting: true}, false},
	}
	for _, c := range cases {
		got, err := NEXUSHeadCTLowRisk(c.age, c.f)
		if err != nil || got != c.want {
			t.Errorf("NEXUSHeadCTLowRisk(%v, %+v) = (%v, %v), want %v", c.age, c.f, got, err, c.want)
		}
	}
}

func TestNEXUSHeadCTLowRisk_errors(t *testing.T) {
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NEXUSHeadCTLowRisk(a, NEXUSHeadCTFindings{}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: %v, want ErrNotFinite", a, err)
		}
	}
	for _, a := range []float64{-0.1, 130.1} {
		if _, err := NEXUSHeadCTLowRisk(a, NEXUSHeadCTFindings{}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: %v, want ErrOutOfRange", a, err)
		}
	}
}

// parity: ledger nexus_head_ct, source cases 1–2 (consumer: any of eight items incl. age ≥ 65 → not low risk).
func TestNEXUSHeadCTLowRisk_parity(t *testing.T) {
	if got, _ := NEXUSHeadCTLowRisk(70, NEXUSHeadCTFindings{}); got { // parity: ledger nexus_head_ct, source case 1
		t.Error("case 1: age 70 must not be low risk")
	}
	if got, _ := NEXUSHeadCTLowRisk(40, NEXUSHeadCTFindings{}); !got { // parity: ledger nexus_head_ct, source case 2
		t.Error("case 2: age 40, none, must be low risk")
	}
}

func TestNEXUSHeadCTLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = NEXUSHeadCTLowRisk(40, NEXUSHeadCTFindings{}) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkNEXUSHeadCTLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = NEXUSHeadCTLowRisk(40, NEXUSHeadCTFindings{})
	}
}
