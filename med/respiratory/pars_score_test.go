package respiratory

import "testing"

// HAND-DERIVED from Biagini Myers 2019 Table III (PMID 30554722): 2 + 3 = 5; every factor = 14; the
// unattainable totals 1 and 13 ("scores of 1 and 13 unattainable") never occur over all 64 combinations.
func TestPARSScore_tableExample(t *testing.T) {
	if got := PARSScore(PARSFindings{ParentalAsthma: true, EarlyWheezingBeforeAge3: true}); got != 5 {
		t.Errorf("got %d, want 5", got)
	}
	if got := PARSScore(PARSFindings{true, true, true, true, true, true}); got != 14 {
		t.Errorf("all: got %d, want 14", got)
	}
	for m := 0; m < 64; m++ {
		f := PARSFindings{m&1 != 0, m&2 != 0, m&4 != 0, m&8 != 0, m&16 != 0, m&32 != 0}
		s := PARSScore(f)
		if s == 1 || s == 13 || s < 0 || s > 14 {
			t.Errorf("combination %06b: score %d is unattainable", m, s)
		}
	}
}

// Each factor's weight alone.
func TestPARSScore_weights(t *testing.T) {
	cases := []struct {
		f    PARSFindings
		want int
	}{
		{PARSFindings{ParentalAsthma: true}, 2}, {PARSFindings{EczemaBeforeAge3: true}, 2},
		{PARSFindings{WheezingApartFromColds: true}, 3}, {PARSFindings{EarlyWheezingBeforeAge3: true}, 3},
		{PARSFindings{SPTPositive2OrMore: true}, 2}, {PARSFindings{AfricanAmericanRace: true}, 2}, {PARSFindings{}, 0},
	}
	for _, c := range cases {
		if got := PARSScore(c.f); got != c.want {
			t.Errorf("%+v: got %d, want %d", c.f, got, c.want)
		}
	}
}

// parity: ledger pars, source case 1: the source's weights (2, 2, 3, 3, 2, 2) agree with Table III; eczema +
// wheezing apart from colds + SPT → 7.
func TestPARSScore_parity(t *testing.T) {
	if got := PARSScore(PARSFindings{EczemaBeforeAge3: true, WheezingApartFromColds: true, SPTPositive2OrMore: true}); got != 7 {
		t.Errorf("got %d, want 7", got)
	}
}

func TestPARSScore_noAlloc(t *testing.T) {
	f := PARSFindings{ParentalAsthma: true, EarlyWheezingBeforeAge3: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = PARSScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPARSScore(b *testing.B) {
	f := PARSFindings{ParentalAsthma: true, EarlyWheezingBeforeAge3: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt = PARSScore(f)
	}
}
