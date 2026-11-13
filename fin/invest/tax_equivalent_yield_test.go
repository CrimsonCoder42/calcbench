package invest

import (
	"math"
	"testing"
)

// MSRB "Understanding Taxable Municipal Bonds": 3.00 % tax-exempt in the 37 % bracket →
// 3.00 %/(1 − 0.37) = 4.76 % (two decimals: ±5e-5).
func TestTaxEquivalentYield_msrb(t *testing.T) {
	got, err := TaxEquivalentYield(0.03, 0.37)
	if err != nil || !near(got, 0.0476, 5e-5) || !near(got, 0.03/0.63, 1e-15) {
		t.Errorf("TaxEquivalentYield = %v, %v; want 0.047619", got, err)
	}
}

func TestTaxEquivalentYield_cases(t *testing.T) {
	for _, c := range []struct{ y, t, want float64 }{{0.03, 0, 0.03}, {0, 0.37, 0}, {0.05, 0.5, 0.1}, {-0.01, 0.5, -0.02}} {
		got, err := TaxEquivalentYield(c.y, c.t)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("TaxEquivalentYield(%v, %v) = %v, %v; want %v", c.y, c.t, got, err, c.want)
		}
	}
}

func TestTaxEquivalentYield_errors(t *testing.T) {
	_, err := TaxEquivalentYield(0.03, 1)
	wantErr(t, "tax 100 %", err, ErrDomain)
	_, err = TaxEquivalentYield(0.03, -0.1)
	wantErr(t, "tax < 0", err, ErrDomain)
	_, err = TaxEquivalentYield(0.03, math.Inf(1))
	wantErr(t, "tax Inf", err, ErrNotFinite)
	_, err = TaxEquivalentYield(1e308, 0.999999) // 1e314 is beyond float64
	wantErr(t, "overflow", err, ErrOverflow)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1} {
		_, err = TaxEquivalentYield(bad, 0.3)
		wantErr(t, "bad yield", err, ErrRate)
	}
}

func TestTaxEquivalentYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = TaxEquivalentYield(0.03, 0.37) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkTaxEquivalentYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = TaxEquivalentYield(0.03, 0.37)
	}
}
