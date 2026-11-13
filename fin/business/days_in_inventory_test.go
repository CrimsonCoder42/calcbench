package business

import (
	"errors"
	"math"
	"testing"
)

// Principles of Finance §6.2, Clear Lake Sporting Goods: $40,000 ÷ $60,000 × 365 = 243.33, printed
// "243 days". Checked to 1e-12 relative against 365 × 4,000,000 ÷ 6,000,000 (two roundings) and to 0.5
// against the printed whole days.
func TestDaysInInventory_sourceExample(t *testing.T) {
	got, err := DaysInInventory(4000000, 6000000, 365)
	if err != nil {
		t.Fatal(err)
	}
	want := 365.0 * 4000000 / 6000000
	if math.Abs(got-want) > 1e-12*want || math.Abs(got-243) > 0.5 {
		t.Fatalf("got %v, want %v (printed 243)", got, want)
	}
}

func TestDaysInInventory_table(t *testing.T) {
	cases := []struct {
		name      string
		inv, cogs int64
		days      int
		want      float64
	}{
		{"zero inventory", 0, 100, 365, 0},
		{"turnover 10", 10, 100, 365, 36.5},
		{"one day", 1, 1, 1, 1},
		{"max inventory", maxExactCents, maxExactCents, 365, 365},
	}
	for _, c := range cases {
		got, err := DaysInInventory(c.inv, c.cogs, c.days)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDaysInInventory_errors(t *testing.T) {
	cases := []struct {
		name      string
		inv, cogs int64
		days      int
		want      error
	}{
		{"zero cogs", 1, 0, 365, ErrZeroDenominator},
		{"zero days", 1, 1, 0, ErrPeriods},
		{"inventory above 2^53", maxExactCents + 1, 1, 365, ErrOutOfRange},
		{"cogs above 2^53", 1, maxExactCents + 1, 365, ErrOutOfRange},
		{"negative inventory", -1, 1, 365, ErrNegative},
		{"negative cogs", 1, -1, 365, ErrNegative},
	}
	for _, c := range cases {
		if _, err := DaysInInventory(c.inv, c.cogs, c.days); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDaysInInventory_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DaysInInventory(4000000, 6000000, 365) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDaysInInventory(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DaysInInventory(4000000, 6000000, 365)
	}
}
