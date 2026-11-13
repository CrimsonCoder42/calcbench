package trauma

import (
	"errors"
	"testing"
)

var sinkPECARNCSpine PECARNCSpineGroup

// Leonard 2024 (PMC11261431): each of the four high-risk factors alone gives the high-risk group, even with
// CART factors present; each of the five CART factors alone gives the CART group; none gives the third group.
func TestPECARNCSpineRule_source(t *testing.T) {
	high := []PECARNCSpineInput{
		{AgeYears: 5, GCS3To8: true},
		{AgeYears: 5, AVPUUnresponsive: true},
		{AgeYears: 5, AbnormalAirwayBreathingCirculation: true},
		{AgeYears: 5, FocalNeurologicDeficit: true},
		{AgeYears: 5, FocalNeurologicDeficit: true, NeckPain: true, MidlineNeckTenderness: true},
	}
	for i, in := range high {
		if got, err := PECARNCSpineRule(in); err != nil || got != PECARNCSpineHighRisk {
			t.Errorf("high %d → %q (%v)", i, got, err)
		}
	}
	cart := []PECARNCSpineInput{
		{AgeYears: 12, NeckPain: true},
		{AgeYears: 12, AlteredMentalStatus: true},
		{AgeYears: 12, SubstantialHeadInjury: true},
		{AgeYears: 12, SubstantialTorsoInjury: true},
		{AgeYears: 12, MidlineNeckTenderness: true},
	}
	for i, in := range cart {
		if got, err := PECARNCSpineRule(in); err != nil || got != PECARNCSpineCART {
			t.Errorf("cart %d → %q (%v)", i, got, err)
		}
	}
	if got, err := PECARNCSpineRule(PECARNCSpineInput{AgeYears: 0}); err != nil || got != PECARNCSpineNone {
		t.Errorf("none → %q (%v)", got, err)
	}
}

// parity: ledger pecarn_cspine, source case 1. No _test.go or JSON example. The entry's input "Neck Pain or
// Midline Neck Tenderness" maps to either CART factor; any one gives the CART group, as the entry's middle tier.
func TestPECARNCSpineRule_parity(t *testing.T) {
	if got, _ := PECARNCSpineRule(PECARNCSpineInput{AgeYears: 10, MidlineNeckTenderness: true}); got != PECARNCSpineCART {
		t.Errorf("got %q, want CART group", got)
	}
}

func TestPECARNCSpineRule_errors(t *testing.T) {
	for _, v := range nonFinite {
		if _, err := PECARNCSpineRule(PECARNCSpineInput{AgeYears: v}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("age %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	for _, v := range []float64{-0.01, 18, 40} {
		if _, err := PECARNCSpineRule(PECARNCSpineInput{AgeYears: v}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: err = %v, want ErrOutOfRange", v, err)
		}
	}
	if _, err := PECARNCSpineRule(PECARNCSpineInput{AgeYears: 17.99}); err != nil {
		t.Errorf("age 17.99: err = %v, want nil", err)
	}
}

func TestPECARNCSpineRule_noAlloc(t *testing.T) {
	in := PECARNCSpineInput{AgeYears: 8, NeckPain: true}
	if n := testing.AllocsPerRun(100, func() { sinkPECARNCSpine, _ = PECARNCSpineRule(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPECARNCSpineRule(b *testing.B) {
	in := PECARNCSpineInput{AgeYears: 8, NeckPain: true}
	b.ReportAllocs()
	for b.Loop() {
		sinkPECARNCSpine, _ = PECARNCSpineRule(in)
	}
}
