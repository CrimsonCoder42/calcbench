package neuro

import (
	"errors"
	"math"
	"testing"
)

// Kuppermann 2009 abstract: the under-2 and 2-and-older rules; each item alone, and the age edge at 2, HAND-DERIVED.
func TestPECARNHeadInjuryVeryLowRisk_example(t *testing.T) {
	cases := []struct {
		age  float64
		f    PECARNHeadFindings
		want bool
	}{
		{1, PECARNHeadFindings{}, true},
		{1, PECARNHeadFindings{NonFrontalScalpHematoma: true}, false},
		{1, PECARNHeadFindings{LOCFiveSecondsOrMore: true}, false},
		{1, PECARNHeadFindings{PalpableSkullFracture: true}, false},
		{1, PECARNHeadFindings{NotActingNormally: true}, false},
		{1, PECARNHeadFindings{Vomiting: true}, true},          // not an item under 2
		{1.99, PECARNHeadFindings{SevereHeadache: true}, true}, // not an item under 2
		{2, PECARNHeadFindings{SevereHeadache: true}, false},
		{2, PECARNHeadFindings{NonFrontalScalpHematoma: true}, true}, // not an item at 2 and older
		{5, PECARNHeadFindings{AnyLossOfConsciousness: true}, false},
		{5, PECARNHeadFindings{Vomiting: true}, false},
		{5, PECARNHeadFindings{BasilarSkullFractureSign: true}, false},
		{5, PECARNHeadFindings{AbnormalMentalStatus: true}, false},
		{1, PECARNHeadFindings{SevereInjuryMechanism: true}, false},
		{5, PECARNHeadFindings{}, true},
	}
	for _, c := range cases {
		got, err := PECARNHeadInjuryVeryLowRisk(c.age, c.f)
		if err != nil || got != c.want {
			t.Errorf("PECARNHeadInjuryVeryLowRisk(%v, %+v) = (%v, %v), want %v", c.age, c.f, got, err, c.want)
		}
	}
}

func TestPECARNHeadInjuryVeryLowRisk_errors(t *testing.T) {
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := PECARNHeadInjuryVeryLowRisk(a, PECARNHeadFindings{}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: %v, want ErrNotFinite", a, err)
		}
	}
	for _, a := range []float64{-1, 18} {
		if _, err := PECARNHeadInjuryVeryLowRisk(a, PECARNHeadFindings{}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: %v, want ErrOutOfRange", a, err)
		}
	}
}

// parity: ledger pecarn, source cases 1–2. The consumer splits high and intermediate risk (Figure 3 of the
// paper, not in the abstract); a child with no finding in either tier is very low risk on both sides.
func TestPECARNHeadInjuryVeryLowRisk_parity(t *testing.T) {
	if got, _ := PECARNHeadInjuryVeryLowRisk(1, PECARNHeadFindings{}); !got { // parity: ledger pecarn, source case 1
		t.Error("case 1")
	}
	if got, _ := PECARNHeadInjuryVeryLowRisk(8, PECARNHeadFindings{Vomiting: true}); got { // parity: ledger pecarn, source case 2 (intermediate tier)
		t.Error("case 2")
	}
}

func TestPECARNHeadInjuryVeryLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = PECARNHeadInjuryVeryLowRisk(5, PECARNHeadFindings{}) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPECARNHeadInjuryVeryLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = PECARNHeadInjuryVeryLowRisk(5, PECARNHeadFindings{})
	}
}
