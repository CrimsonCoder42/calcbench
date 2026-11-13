package neuro

import (
	"errors"
	"math"
	"testing"
)

// Palchak 2003 abstract: five predictors; scalp hematoma only at ≤ 2 years (edge on both sides), HAND-DERIVED.
func TestUCDavisHeadTraumaLowRisk_example(t *testing.T) {
	cases := []struct {
		age  float64
		f    UCDavisHeadTraumaFindings
		want bool
	}{
		{5, UCDavisHeadTraumaFindings{}, true},
		{2, UCDavisHeadTraumaFindings{ScalpHematoma: true}, false},
		{2.01, UCDavisHeadTraumaFindings{ScalpHematoma: true}, true},
		{10, UCDavisHeadTraumaFindings{Headache: true}, false},
		{10, UCDavisHeadTraumaFindings{HistoryOfVomiting: true}, false},
		{10, UCDavisHeadTraumaFindings{AbnormalMentalStatus: true}, false},
		{10, UCDavisHeadTraumaFindings{SkullFractureSigns: true}, false},
		{17.99, UCDavisHeadTraumaFindings{}, true},
	}
	for _, c := range cases {
		got, err := UCDavisHeadTraumaLowRisk(c.age, c.f)
		if err != nil || got != c.want {
			t.Errorf("UCDavisHeadTraumaLowRisk(%v, %+v) = (%v, %v), want %v", c.age, c.f, got, err, c.want)
		}
	}
}

func TestUCDavisHeadTraumaLowRisk_errors(t *testing.T) {
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := UCDavisHeadTraumaLowRisk(a, UCDavisHeadTraumaFindings{}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: %v, want ErrNotFinite", a, err)
		}
	}
	for _, a := range []float64{-0.1, 18, 18.5} {
		if _, err := UCDavisHeadTraumaLowRisk(a, UCDavisHeadTraumaFindings{}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: %v, want ErrOutOfRange", a, err)
		}
	}
}

// parity: ledger palchak_pediatric_head_trauma, source cases 1–3 (consumer counts scalp hematoma only when age ≤ 2).
func TestUCDavisHeadTraumaLowRisk_parity(t *testing.T) {
	if got, _ := UCDavisHeadTraumaLowRisk(1, UCDavisHeadTraumaFindings{ScalpHematoma: true}); got { // parity: ledger palchak_pediatric_head_trauma, source case 1
		t.Error("case 1")
	}
	if got, _ := UCDavisHeadTraumaLowRisk(6, UCDavisHeadTraumaFindings{ScalpHematoma: true}); !got { // parity: ledger palchak_pediatric_head_trauma, source case 2
		t.Error("case 2")
	}
	if got, _ := UCDavisHeadTraumaLowRisk(6, UCDavisHeadTraumaFindings{SkullFractureSigns: true}); got { // parity: ledger palchak_pediatric_head_trauma, source case 3
		t.Error("case 3")
	}
}

func TestUCDavisHeadTraumaLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = UCDavisHeadTraumaLowRisk(5, UCDavisHeadTraumaFindings{}) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkUCDavisHeadTraumaLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = UCDavisHeadTraumaLowRisk(5, UCDavisHeadTraumaFindings{})
	}
}
