package respiratory

import (
	"errors"
	"math"
	"testing"
)

func pompeBase() POMPECFindings {
	return POMPECFindings{WeightLb: 150, HighestRespRatePerMin: 20, SaO2Percent: 95}
}

// pompeRef evaluates the printed equation directly, in the paper's own form 100·(1 − 1/(1 + e^L)).
func pompeRef(f POMPECFindings) float64 {
	b := func(x bool) float64 {
		if x {
			return 1
		}
		return 0
	}
	l := 3.718 + b(f.DoNotResuscitate)*1.551 + b(f.RespiratoryDistress)*0.800 + b(f.UnilateralLimbSwelling)*0.734 +
		b(f.AlteredMentalStatus)*1.473 + b(f.HeartRateAbove99)*1.027 + f.HighestRespRatePerMin*0.044 +
		f.SaO2Percent*-0.063 + f.WeightLb*-0.012
	return 100 * (1 - 1/(1+math.Exp(l)))
}

// HAND-DERIVED from Kline 2012's equation (PMID 22475313): 150 lb, rate 20, SaO₂ 95 → L = −3.187,
// P = 100/(1 + e^3.187) = 3.966 %. Tolerance 1e-3 percentage points (four printed figures).
func TestPOMPECMortalityPercent_equationExample(t *testing.T) {
	got, err := POMPECMortalityPercent(pompeBase())
	if err != nil || math.Abs(got-3.966) > 1e-3 {
		t.Errorf("got (%v, %v), want 3.966", got, err)
	}
}

// Every item against the paper's own form of the equation; tolerance 1e-9 relative (the two forms differ only
// by rounding).
func TestPOMPECMortalityPercent_matchesPrintedForm(t *testing.T) {
	cases := []POMPECFindings{
		pompeBase(),
		{WeightLb: 120, HighestRespRatePerMin: 32, SaO2Percent: 86, HeartRateAbove99: true, AlteredMentalStatus: true, RespiratoryDistress: true, DoNotResuscitate: true, UnilateralLimbSwelling: true},
		{WeightLb: 250, HighestRespRatePerMin: 12, SaO2Percent: 99},
		{WeightLb: 180, HighestRespRatePerMin: 24, SaO2Percent: 91, DoNotResuscitate: true},
		{WeightLb: 1, HighestRespRatePerMin: 0, SaO2Percent: 0},
	}
	for _, f := range cases {
		got, err := POMPECMortalityPercent(f)
		want := pompeRef(f)
		if err != nil || math.Abs(got-want) > 1e-9*math.Max(1, want) {
			t.Errorf("%+v: got (%v, %v), want %v", f, got, err, want)
		}
	}
}

// parity: ledger pompe_c, source case 1: the source uses respiratory distress 0.799 and SaO₂ −0.0636 where the
// paper's equation prints 0.800 and −0.063. Asserting the PRIMARY: 150 lb, rate 20, SaO₂ 95, respiratory
// distress → L = −2.387, 8.42 % (the source's coefficients give 7.98 %). Tolerance 0.01 percentage points.
func TestPOMPECMortalityPercent_parity(t *testing.T) {
	f := pompeBase()
	f.RespiratoryDistress = true
	got, err := POMPECMortalityPercent(f)
	want := 100 / (1 + math.Exp(2.387))
	if err != nil || math.Abs(got-want) > 1e-2 {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
}

// The kg twin agrees with the lb function to 1e-12 relative (one exact conversion).
func TestPOMPECMortalityPercentKg_agrees(t *testing.T) {
	f := pompeBase()
	lb, _ := POMPECMortalityPercent(f)
	kg, err := POMPECMortalityPercentKg(POMPECFindings{HighestRespRatePerMin: 20, SaO2Percent: 95}, 150*kgPerPound)
	if err != nil || math.Abs(kg-lb) > 1e-12*lb {
		t.Errorf("kg twin: got (%v, %v), want %v", kg, err, lb)
	}
}

func TestPOMPECMortalityPercent_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*POMPECFindings)
		want error
	}{
		{"weight 0", func(f *POMPECFindings) { f.WeightLb = 0 }, ErrOutOfRange},
		{"rate -1", func(f *POMPECFindings) { f.HighestRespRatePerMin = -1 }, ErrOutOfRange},
		{"SaO2 101", func(f *POMPECFindings) { f.SaO2Percent = 101 }, ErrOutOfRange},
		{"SaO2 -1", func(f *POMPECFindings) { f.SaO2Percent = -1 }, ErrOutOfRange},
		{"weight NaN", func(f *POMPECFindings) { f.WeightLb = math.NaN() }, ErrNotFinite},
		{"rate +Inf", func(f *POMPECFindings) { f.HighestRespRatePerMin = math.Inf(1) }, ErrNotFinite},
		{"SaO2 -Inf", func(f *POMPECFindings) { f.SaO2Percent = math.Inf(-1) }, ErrNotFinite},
	}
	for _, c := range cases {
		f := pompeBase()
		c.mod(&f)
		if _, err := POMPECMortalityPercent(f); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
	for _, bad := range []struct {
		kg   float64
		want error
	}{{0, ErrOutOfRange}, {-1, ErrOutOfRange}, {math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}} {
		if _, err := POMPECMortalityPercentKg(pompeBase(), bad.kg); !errors.Is(err, bad.want) {
			t.Errorf("kg %v: err = %v, want %v", bad.kg, err, bad.want)
		}
	}
}

func TestPOMPECMortalityPercent_noAlloc(t *testing.T) {
	f := pompeBase()
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = POMPECMortalityPercent(f) }); n != 0 {
		t.Errorf("lb allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = POMPECMortalityPercentKg(f, 68) }); n != 0 {
		t.Errorf("kg allocs = %v, want 0", n)
	}
}

func BenchmarkPOMPECMortalityPercent(b *testing.B) {
	f := pompeBase()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = POMPECMortalityPercent(f)
	}
}

func BenchmarkPOMPECMortalityPercentKg(b *testing.B) {
	f := pompeBase()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = POMPECMortalityPercentKg(f, 68)
	}
}
