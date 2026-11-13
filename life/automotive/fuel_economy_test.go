package automotive

import (
	"errors"
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from the 40 CFR 600.002 definition (distance per volume of fuel):
// 300 mi / 12 gal = 25 mpg; 600 km / 45 L = 13.333… km/L; 45 L over 600 km → 100 × 45 / 600 = 7.5 L/100 km;
// 100 km on 8 L → 8 L/100 km. Cost: 300 mi / 25 mpg = 12 gal × 350 c = 4200 c; 250 km / 12.5 km/L = 20 L × 189 c = 3780 c;
// 100 mi / 30 mpg × 349 c = 1163.33… → 1163 c; 1 mi / 2 mpg × 3 c = 1.5 → 2 (half to even); 1 mi / 2 mpg × 5 c = 2.5 → 2.
// Tolerance 1e-15 relative for the float results (one divide, one multiply).
func TestFuelEconomy_example(t *testing.T) {
	tests := []struct{ d, f, want float64 }{{300, 12, 25}, {600, 45, 600.0 / 45}}
	for _, tc := range tests {
		if got, err := FuelEconomy(tc.d, tc.f); err != nil || math.Abs(got-tc.want) > 1e-15*tc.want {
			t.Errorf("FuelEconomy(%v, %v) = (%v, %v), want %v", tc.d, tc.f, got, err, tc.want)
		}
	}
}

func TestFuelConsumptionPer100_example(t *testing.T) {
	tests := []struct{ f, d, want float64 }{{45, 600, 7.5}, {8, 100, 8}}
	for _, tc := range tests {
		if got, err := FuelConsumptionPer100(tc.f, tc.d); err != nil || math.Abs(got-tc.want) > 1e-15*tc.want {
			t.Errorf("FuelConsumptionPer100(%v, %v) = (%v, %v), want %v", tc.f, tc.d, got, err, tc.want)
		}
	}
}

// Reciprocal identity: FuelEconomy × FuelConsumptionPer100 = 100 for the same trip. Tolerance 1e-14 relative.
func TestFuelEconomy_reciprocal(t *testing.T) {
	for _, d := range []float64{1, 37.5, 300, 1234.5} {
		for _, f := range []float64{0.5, 4, 12.25} {
			e, _ := FuelEconomy(d, f)
			c, _ := FuelConsumptionPer100(f, d)
			if math.Abs(e*c-100) > 1e-12 {
				t.Errorf("d=%v f=%v: %v × %v = %v, want 100", d, f, e, c, e*c)
			}
		}
	}
}

func TestFuelCostCents_example(t *testing.T) {
	tests := []struct {
		d, e  float64
		price int64
		want  int64
	}{
		{300, 25, 350, 4200},
		{250, 12.5, 189, 3780},
		{100, 30, 349, 1163},
		{1, 2, 3, 2},
		{1, 2, 5, 2},
		{0, 25, 350, 0},
		{300, 25, 0, 0},
	}
	for _, tc := range tests {
		if got, err := FuelCostCents(tc.d, tc.e, tc.price); err != nil || got != tc.want {
			t.Errorf("FuelCostCents(%v, %v, %d) = (%d, %v), want %d", tc.d, tc.e, tc.price, got, err, tc.want)
		}
	}
}

func TestFuelEconomy_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	type call func() error
	cases := []struct {
		name string
		f    call
		want error
	}{
		{"economy NaN", func() error { _, err := FuelEconomy(nan, 1); return err }, ErrNotFinite},
		{"economy Inf fuel", func() error { _, err := FuelEconomy(1, inf); return err }, ErrNotFinite},
		{"economy zero fuel", func() error { _, err := FuelEconomy(1, 0); return err }, ErrOutOfRange},
		{"economy negative distance", func() error { _, err := FuelEconomy(-1, 1); return err }, ErrOutOfRange},
		{"economy overflow", func() error { _, err := FuelEconomy(math.MaxFloat64, 0.5); return err }, ErrDomain},
		{"economy underflow", func() error { _, err := FuelEconomy(5e-324, 1e300); return err }, ErrDomain},
		{"consumption zero distance", func() error { _, err := FuelConsumptionPer100(1, 0); return err }, ErrOutOfRange},
		{"consumption -Inf", func() error { _, err := FuelConsumptionPer100(math.Inf(-1), 1); return err }, ErrNotFinite},
		{"consumption overflow", func() error { _, err := FuelConsumptionPer100(math.MaxFloat64, 1); return err }, ErrDomain},
		{"cost NaN distance", func() error { _, err := FuelCostCents(nan, 25, 350); return err }, ErrNotFinite},
		{"cost negative distance", func() error { _, err := FuelCostCents(-1, 25, 350); return err }, ErrOutOfRange},
		{"cost zero economy", func() error { _, err := FuelCostCents(1, 0, 350); return err }, ErrOutOfRange},
		{"cost negative price", func() error { _, err := FuelCostCents(1, 25, -1); return err }, ErrOutOfRange},
		{"cost overflow", func() error { _, err := FuelCostCents(math.MaxFloat64, 1e-300, 350); return err }, ErrDomain},
		{"cost beyond int64", func() error { _, err := FuelCostCents(1e18, 1, 100); return err }, ErrDomain},
	}
	for _, tc := range cases {
		if err := tc.f(); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestFuelEconomy_noAllocation(t *testing.T) {
	var s float64
	var c int64
	if n := testing.AllocsPerRun(100, func() {
		s, _ = FuelEconomy(300, 12)
		s, _ = FuelConsumptionPer100(45, 600)
		c, _ = FuelCostCents(300, 25, 350)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	_, _ = s, c
}

func BenchmarkFuelEconomy(b *testing.B) {
	b.ReportAllocs()
	var s float64
	for b.Loop() {
		s, _ = FuelEconomy(300, 12)
	}
	_ = s
}

func BenchmarkFuelConsumptionPer100(b *testing.B) {
	b.ReportAllocs()
	var s float64
	for b.Loop() {
		s, _ = FuelConsumptionPer100(45, 600)
	}
	_ = s
}

func BenchmarkFuelCostCents(b *testing.B) {
	b.ReportAllocs()
	var c int64
	for b.Loop() {
		c, _ = FuelCostCents(300, 25, 350)
	}
	_ = c
}
