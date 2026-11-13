package invest

import "testing"

// OpenStax Principles of Finance §15.1: 3M $5.91 dividends on $175.00 → 3.38 %; Target $2.68 on
// $128.74 → $0.0208 per dollar (both rounded in the source, ±5e-5).
func TestDividendYield_openStax(t *testing.T) {
	got, err := DividendYield(591, 17500)
	if err != nil || !near(got, 0.0338, 5e-5) || got != 591.0/17500.0 {
		t.Errorf("3M: %v, %v; want 591/17500 ≈ 0.0338", got, err)
	}
	got, err = DividendYield(268, 12874)
	if err != nil || !near(got, 0.0208, 5e-5) {
		t.Errorf("Target: %v, %v; want ≈ 0.0208", got, err)
	}
}

func TestDividendYield_cases(t *testing.T) {
	for _, c := range []struct {
		name     string
		d, price int64
		want     float64
	}{{"no dividend", 0, 100, 0}, {"max", maxCents, maxCents, 1}, {"small", 1, 100, 0.01}} {
		got, err := DividendYield(c.d, c.price)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDividendYield_errors(t *testing.T) {
	_, err := DividendYield(10, 0)
	wantErr(t, "price 0", err, ErrZeroDenominator)
	_, err = DividendYield(-1, 10)
	wantErr(t, "dividend < 0", err, ErrNegative)
	_, err = DividendYield(1, -10)
	wantErr(t, "price < 0", err, ErrNegative)
	_, err = DividendYield(maxCents+1, 10)
	wantErr(t, "dividend > 2^53", err, ErrOutOfRange)
}

func TestDividendYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = DividendYield(591, 17500) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkDividendYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = DividendYield(591, 17500)
	}
}
