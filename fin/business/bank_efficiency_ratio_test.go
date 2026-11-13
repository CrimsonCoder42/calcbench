package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// The UBPR guide (item 32, UBPRE088) states the formula without numbers. Hand-checked: $60,000 ÷ ($80,000 +
// $20,000) = 0.6, the UBPR's 60 (percent). Exact: 0.6 is the correctly rounded 6/10.
func TestBankEfficiencyRatio_sourceExample(t *testing.T) {
	got, err := BankEfficiencyRatio(6000000, 8000000, 2000000)
	if err != nil || got != 0.6 {
		t.Fatalf("got %v, %v; want 0.6", got, err)
	}
}

func TestBankEfficiencyRatio_table(t *testing.T) {
	cases := []struct {
		name                string
		expense, nii, nonii int64
		want                float64
	}{
		{"zero expense", 0, 10, 10, 0},
		{"negative noninterest income", 50, 120, -20, 0.5},
		{"negative net interest income", 50, -20, 120, 0.5},
		{"negative sum", 50, -100, 0, -0.5},
		{"max", maxExactCents, maxExactCents, 0, 1},
	}
	for _, c := range cases {
		got, err := BankEfficiencyRatio(c.expense, c.nii, c.nonii)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBankEfficiencyRatio_errors(t *testing.T) {
	cases := []struct {
		name                string
		expense, nii, nonii int64
		want                error
	}{
		{"zero income", 1, 0, 0, ErrZeroDenominator},
		{"incomes cancel", 1, 5, -5, ErrZeroDenominator},
		{"expense above 2^53", maxExactCents + 1, 1, 1, ErrOutOfRange},
		{"nii above 2^53", 1, maxExactCents + 1, 1, ErrOutOfRange},
		{"nonii below -2^53", 1, 1, -maxExactCents - 1, ErrOutOfRange},
		{"negative expense", -1, 1, 1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := BankEfficiencyRatio(c.expense, c.nii, c.nonii); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle against math/big: exact when |nii + nonii| ≤ 2⁵³, else within 2⁻⁵² relative (two roundings).
func TestBankEfficiencyRatio_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		e := int64(rng.Uint64N(maxExactCents + 1))
		nii := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		nonii := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		sum := nii + nonii
		if sum == 0 {
			continue
		}
		got, err := BankEfficiencyRatio(e, nii, nonii)
		if err != nil {
			t.Fatal(err)
		}
		want := ratOracle(e, sum)
		if sum >= -maxExactCents && sum <= maxExactCents {
			if got != want {
				t.Fatalf("(%d,%d,%d) = %v, want %v", e, nii, nonii, got, want)
			}
		} else if math.Abs(got-want) > math.Abs(want)*0x1p-52 {
			t.Fatalf("(%d,%d,%d) = %v, want %v within 2^-52 relative", e, nii, nonii, got, want)
		}
	}
}

func TestBankEfficiencyRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = BankEfficiencyRatio(6000000, 8000000, 2000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkBankEfficiencyRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BankEfficiencyRatio(6000000, 8000000, 2000000)
	}
}
