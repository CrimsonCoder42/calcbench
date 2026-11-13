package interest

import "testing"

func TestGrowthRatePerPeriod_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 §7.1, Eq. 14, Example 17: sales
	// 10,503.0 → 14,146.4 over 4 years → g = 0.077291; Example 18: 8.52 → 7.35 million over 5 → −0.029111.
	got, err := GrowthRatePerPeriod(1050300, 1414640, 4)
	wantFloat(t, "QIA Example 17", got, err, 0.077291, 5e-7)
	got, err = GrowthRatePerPeriod(852, 735, 5)
	wantFloat(t, "QIA Example 18", got, err, -0.029111, 5e-7)
	// Principles of Finance (OpenStax 2022) §7.3: PV $1,000, FV $1,125.51, N 4 → 3 % (to the cent's rounding).
	got, err = GrowthRatePerPeriod(100000, 112551, 4)
	wantFloat(t, "§7.3", got, err, 0.03, 1e-6)
}

func TestGrowthRatePerPeriod_Table(t *testing.T) {
	got, err := GrowthRatePerPeriod(5000, 5000, 7)
	wantFloat(t, "no growth", got, err, 0, 0)
	got, err = GrowthRatePerPeriod(5000, 10000, 1)
	wantFloat(t, "doubling in one period", got, err, 1, 1e-15)
	got, err = GrowthRatePerPeriod(-5000, -10000, 1)
	wantFloat(t, "both negative", got, err, 1, 1e-15)
	got, err = GrowthRatePerPeriod(1, maxExactCents, 53)
	wantFloat(t, "1 cent to 2^53 in 53 periods", got, err, 1, 1e-13)
	got, err = GrowthRatePerPeriod(maxExactCents, maxExactCents-1, 1)
	wantRel(t, "tiny decline keeps precision", got, err, -1.0/float64(maxExactCents), 1e-12)
}

func TestGrowthRatePerPeriod_Errors(t *testing.T) {
	_, err := GrowthRatePerPeriod(maxExactCents+1, 1, 1)
	wantErr(t, "pv 2^53+1", err, ErrOutOfRange)
	_, err = GrowthRatePerPeriod(1, -maxExactCents-1, 1)
	wantErr(t, "fv -2^53-1", err, ErrOutOfRange)
	_, err = GrowthRatePerPeriod(0, 100, 1)
	wantErr(t, "zero pv", err, ErrDomain)
	_, err = GrowthRatePerPeriod(100, 0, 1)
	wantErr(t, "zero fv", err, ErrDomain)
	_, err = GrowthRatePerPeriod(-100, 100, 1)
	wantErr(t, "mixed signs", err, ErrDomain)
	_, err = GrowthRatePerPeriod(100, 200, 0)
	wantErr(t, "zero periods", err, ErrPeriods)
}

func TestGrowthRatePerPeriod_Oracle(t *testing.T) {
	// Round trip with FutureValue: the rate that grows PV to FV(PV, r, n) is r, to the cent's rounding.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		pv := 1_000_000_000 + rng.Int64N(1_000_000_000)
		r := rng.Float64()*0.2 - 0.05
		n := 1 + rng.IntN(40)
		fv, err := FutureValue(pv, r, n)
		if err != nil {
			t.Fatal(err)
		}
		got, err := GrowthRatePerPeriod(pv, fv, n)
		wantFloat(t, "round trip", got, err, r, 1e-9)
	}
}

func TestGrowthRatePerPeriod_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "GrowthRatePerPeriod", func() { sinkFloat, _ = GrowthRatePerPeriod(1050300, 1414640, 4) })
}

func BenchmarkGrowthRatePerPeriod(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GrowthRatePerPeriod(1050300, 1414640, 4)
	}
}
