package invest

import "testing"

// US Treasury, "Price, Yield and Rate Calculations for a Treasury Bill" (2004), p. 2: price 99.937778,
// 28 days, 366-day year → i = 0.008138368141143 (0.814 %). Treasury rounds each step to 15 decimals,
// so the comparison is to 1e-14.
func TestBondEquivalentYield_treasury(t *testing.T) {
	got, err := BondEquivalentYield(100000000, 99937778, 28, 366)
	if err != nil || !near(got, 0.008138368141143, 1e-14) {
		t.Errorf("BondEquivalentYield = %v, %v; want 0.008138368141143", got, err)
	}
	if !near(got, 0.00814, 5e-6) {
		t.Errorf("rounded %v, want 0.814 %%", got)
	}
	bdy, _ := BankDiscountYield(100000000, 99937778, 28)
	if got <= bdy {
		t.Errorf("investment rate %v not above discount rate %v", got, bdy)
	}
}

func TestBondEquivalentYield_cases(t *testing.T) {
	cases := []struct {
		name         string
		face, price  int64
		days, inYear int
		want         float64
	}{
		{"at face", 100000, 100000, 91, 365, 0},
		{"half year 365", 100000, 98000, 182, 365, 2000.0 / 98000 * 365 / 182},
		{"half year 366", 100000, 98000, 183, 366, 2000.0 / 98000 * 2},
		{"one day", 36500, 36499, 1, 365, 365.0 / 36499},
	}
	for _, c := range cases {
		got, err := BondEquivalentYield(c.face, c.price, c.days, c.inYear)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBondEquivalentYield_errors(t *testing.T) {
	_, err := BondEquivalentYield(100, 0, 10, 365)
	wantErr(t, "price 0", err, ErrZeroDenominator)
	_, err = BondEquivalentYield(0, 100, 10, 365)
	wantErr(t, "face 0", err, ErrDomain)
	_, err = BondEquivalentYield(-1, 100, 10, 365)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = BondEquivalentYield(100, maxCents+1, 10, 365)
	wantErr(t, "price > 2^53", err, ErrOutOfRange)
	_, err = BondEquivalentYield(100, 99, 0, 365)
	wantErr(t, "days 0", err, ErrPeriods)
	_, err = BondEquivalentYield(100, 99, 10, 360)
	wantErr(t, "360-day year", err, ErrDomain)
	_, err = BondEquivalentYield(100, 99, 183, 365)
	wantErr(t, "over half a year", err, ErrDomain)
}

func TestBondEquivalentYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = BondEquivalentYield(100000000, 99937778, 28, 366) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBondEquivalentYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = BondEquivalentYield(100000000, 99937778, 28, 366)
	}
}
