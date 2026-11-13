package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §11.2: $2.00 a year forever at 6 % → $33.33.
func TestZeroGrowthStockPrice_openStax(t *testing.T) {
	got, err := ZeroGrowthStockPrice(200, 0.06)
	if err != nil || got != 3333 {
		t.Errorf("ZeroGrowthStockPrice(200, 0.06) = %d, %v; want 3333", got, err)
	}
}

func TestZeroGrowthStockPrice_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		d    int64
		r    float64
		want int64
	}{
		{"zero dividend", 0, 0.05, 0},
		{"tie to even down", 1, 2, 0}, // 0.5 → 0
		{"tie to even up", 3, 2, 2},   // 1.5 → 2
		{"max at r = 1", maxCents, 1, maxCents},
	} {
		got, err := ZeroGrowthStockPrice(c.d, c.r)
		if err != nil || got != c.want {
			t.Errorf("%s: %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestZeroGrowthStockPrice_errors(t *testing.T) {
	_, err := ZeroGrowthStockPrice(100, 0)
	wantErr(t, "r = 0", err, ErrDomain)
	_, err = ZeroGrowthStockPrice(100, -0.05)
	wantErr(t, "r < 0", err, ErrDomain)
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1} {
		_, err = ZeroGrowthStockPrice(100, bad)
		wantErr(t, "bad r", err, ErrRate)
	}
	_, err = ZeroGrowthStockPrice(-1, 0.05)
	wantErr(t, "dividend < 0", err, ErrNegative)
	_, err = ZeroGrowthStockPrice(maxCents+1, 0.05)
	wantErr(t, "dividend > 2^53", err, ErrOutOfRange)
	_, err = ZeroGrowthStockPrice(maxCents, 0.5)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestZeroGrowthStockPrice_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = ZeroGrowthStockPrice(200, 0.06) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkZeroGrowthStockPrice(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = ZeroGrowthStockPrice(200, 0.06)
	}
}
