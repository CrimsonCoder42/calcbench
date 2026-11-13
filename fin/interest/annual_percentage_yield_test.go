package interest

import (
	"math"
	"testing"
)

// roundAPY rounds a decimal yield to the nearest 0.01 percentage point, as 12 CFR 1030.3(f)(1)
// requires of the disclosed APY, and returns it in percent.
func roundAPY(y float64) float64 { return math.Round(y*10000) / 100 }

func TestAnnualPercentageYield_SourceExample(t *testing.T) {
	// 12 CFR 1030 Appendix A Part I.A: $61.68 interest on $1,000 over 365 days → APY 6.17 %.
	got, err := AnnualPercentageYield(100000, 6168, 365)
	wantFloat(t, "App. A NOW account", got, err, 0.06168, 1e-15)
	if roundAPY(got) != 6.17 {
		t.Errorf("rounded APY %v, want 6.17", roundAPY(got))
	}
	// Part I.A: $30.37 on a $1,000 six-month CD of 182 days → APY 6.18 %.
	got, err = AnnualPercentageYield(100000, 3037, 182)
	if err != nil || roundAPY(got) != 6.18 {
		t.Errorf("App. A CD: got %v (%v %%), %v; want 6.18 %%", got, roundAPY(got), err)
	}
	// Part I.B: $26.68 on $1,000 over 183 days (91 + 92) → APY 5.39 %.
	got, err = AnnualPercentageYield(100000, 2668, 183)
	if err != nil || roundAPY(got) != 5.39 {
		t.Errorf("App. A stepped CD: got %v (%v %%), %v; want 5.39 %%", got, roundAPY(got), err)
	}
}

func TestAnnualPercentageYield_Table(t *testing.T) {
	got, err := AnnualPercentageYield(100000, 0, 90)
	wantFloat(t, "no interest", got, err, 0, 0)
	got, err = AnnualPercentageYield(100000, 5000, 730)
	wantRel(t, "two-year term", got, err, math.Sqrt(1.05)-1, 1e-14)
	got, err = AnnualPercentageYield(100000, -1000, 365)
	wantRel(t, "negative interest (a fee)", got, err, -0.01, 1e-14)
	got, err = AnnualPercentageYield(maxExactCents, 1, 365)
	wantRel(t, "max exact principal", got, err, 1/float64(maxExactCents), 1e-12)
}

func TestAnnualPercentageYield_Errors(t *testing.T) {
	_, err := AnnualPercentageYield(maxExactCents+1, 1, 365)
	wantErr(t, "principal 2^53+1", err, ErrOutOfRange)
	_, err = AnnualPercentageYield(100, -maxExactCents-1, 365)
	wantErr(t, "interest -2^53-1", err, ErrOutOfRange)
	_, err = AnnualPercentageYield(0, 1, 365)
	wantErr(t, "zero principal", err, ErrDomain)
	_, err = AnnualPercentageYield(-100, 1, 365)
	wantErr(t, "negative principal", err, ErrDomain)
	_, err = AnnualPercentageYield(100, -100, 365)
	wantErr(t, "loses everything", err, ErrDomain)
	_, err = AnnualPercentageYield(100, 1, 0)
	wantErr(t, "zero days", err, ErrPeriods)
	_, err = AnnualPercentageYield(1, maxExactCents, 1)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestAnnualPercentageYield_Oracle(t *testing.T) {
	// Identity: for a 365-day term the APY is Interest/Principal (Appendix A Part I.A).
	rng := testRNG()
	for k := 0; k < 300; k++ {
		p := 1 + rng.Int64N(1_000_000_000)
		i := rng.Int64N(p/2 + 1)
		got, err := AnnualPercentageYield(p, i, 365)
		wantRel(t, "365-day identity", got, err, float64(i)/float64(p), 1e-14)
	}
}

func TestAnnualPercentageYield_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnualPercentageYield", func() { sinkFloat, _ = AnnualPercentageYield(100000, 3037, 182) })
}

func BenchmarkAnnualPercentageYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = AnnualPercentageYield(100000, 3037, 182)
	}
}
