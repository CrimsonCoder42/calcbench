package invest

import "testing"

// US Treasury, "Price, Yield and Rate Calculations for a Treasury Bill" (2004), p. 1: price 99.937778
// per 100 for 28 days → d = 0.007999971428571 (step 4, 15 decimals; Treasury rounds each step, so the
// comparison is to 1e-15). Face $1,000,000 makes the price whole cents.
func TestBankDiscountYield_treasury(t *testing.T) {
	got, err := BankDiscountYield(100000000, 99937778, 28)
	if err != nil || !near(got, 0.007999971428571, 1e-15) {
		t.Errorf("BankDiscountYield = %v, %v; want 0.007999971428571", got, err)
	}
	if !near(got, 0.00800, 5e-6) { // displayed 0.800 %
		t.Errorf("rounded rate %v, want 0.800 %%", got)
	}
}

// TreasuryDirect "Understanding Pricing": $1,000 26-week bill at 0.145 % sells for $999.27; the
// rate recovered from the rounded price is 0.145 % to ±2e-5 (the price is rounded to the cent).
func TestBankDiscountYield_treasuryDirectPricing(t *testing.T) {
	got, err := BankDiscountYield(100000, 99927, 182)
	if err != nil || !near(got, 0.00145, 2e-5) {
		t.Errorf("BankDiscountYield = %v, %v; want ≈ 0.00145", got, err)
	}
}

func TestBankDiscountYield_cases(t *testing.T) {
	cases := []struct {
		name        string
		face, price int64
		days        int
		want        float64
	}{
		{"at face", 100000, 100000, 91, 0},
		{"premium", 100000, 100100, 360, -0.001},
		{"one day", 36000, 35999, 1, 0.01},
		{"max", maxCents, maxCents, 1, 0},
	}
	for _, c := range cases {
		got, err := BankDiscountYield(c.face, c.price, c.days)
		if err != nil || !near(got, c.want, 1e-15) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBankDiscountYield_errors(t *testing.T) {
	_, err := BankDiscountYield(0, 100, 10)
	wantErr(t, "face 0", err, ErrZeroDenominator)
	_, err = BankDiscountYield(100, 0, 10)
	wantErr(t, "price 0", err, ErrDomain)
	_, err = BankDiscountYield(-100, 100, 10)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = BankDiscountYield(100, -1, 10)
	wantErr(t, "price < 0", err, ErrNegative)
	_, err = BankDiscountYield(maxCents+1, 100, 10)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = BankDiscountYield(100, 99, 0)
	wantErr(t, "days 0", err, ErrPeriods)
}

func TestBankDiscountYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = BankDiscountYield(100000000, 99937778, 28) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkBankDiscountYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = BankDiscountYield(100000000, 99937778, 28)
	}
}
