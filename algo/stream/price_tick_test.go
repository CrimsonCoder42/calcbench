package stream

import (
	"errors"
	"math"
	"testing"
)

func TestPriceToTicks_cases(t *testing.T) {
	tests := []struct {
		name            string
		price, tickSize float64
		want            int64
	}{
		{"doc example 101.23 at 0.05 is 2024.6 ticks", 101.23, 0.05, 2025},
		{"decimal below the tick in binary: 1.15/0.05 = 22.999999999999996", 1.15, 0.05, 23},
		{"half tick rounds away from zero", 0.375, 0.25, 2},
		{"negative half tick rounds away from zero", -0.375, 0.25, -2},
		{"half of one tick", 0.125, 0.25, 1},
		{"negative half of one tick", -0.125, 0.25, -1},
		{"zero", 0, 0.01, 0},
		{"whole ticks", 506.25, 0.25, 2025},
		{"int64 minimum", -float64(1 << 63), 1, math.MinInt64},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := PriceToTicks(tc.price, tc.tickSize); err != nil || got != tc.want {
				t.Errorf("PriceToTicks(%v, %v) = (%d, %v), want %d", tc.price, tc.tickSize, got, err, tc.want)
			}
		})
	}
}

func TestPriceToTicks_errors(t *testing.T) {
	for _, ts := range []float64{0, -0.01, math.NaN(), math.Inf(1)} {
		if _, err := PriceToTicks(1, ts); !errors.Is(err, ErrTickSize) {
			t.Errorf("PriceToTicks(1, %v) error = %v, want ErrTickSize", ts, err)
		}
	}
	for _, p := range []float64{math.NaN(), math.Inf(-1)} {
		if _, err := PriceToTicks(p, 0.01); !errors.Is(err, ErrNotFinite) {
			t.Errorf("PriceToTicks(%v, 0.01) error = %v, want ErrNotFinite", p, err)
		}
	}
	for _, c := range [][2]float64{{float64(1 << 63), 1}, {1e300, 1e-10}, {1, 5e-324}, {-1e19, 1}} {
		if _, err := PriceToTicks(c[0], c[1]); !errors.Is(err, ErrOverflow) {
			t.Errorf("PriceToTicks(%v, %v) error = %v, want ErrOverflow", c[0], c[1], err)
		}
	}
}

func TestTicksToPrice(t *testing.T) {
	for _, c := range []struct {
		ticks    int64
		tickSize float64
		want     float64
	}{{2025, 0.25, 506.25}, {-4, 0.5, -2}, {0, 0.01, 0}} {
		if got, err := TicksToPrice(c.ticks, c.tickSize); err != nil || got != c.want {
			t.Errorf("TicksToPrice(%d, %v) = (%v, %v), want %v", c.ticks, c.tickSize, got, err, c.want)
		}
	}
	if _, err := TicksToPrice(1, 0); !errors.Is(err, ErrTickSize) {
		t.Errorf("TicksToPrice(1, 0) error = %v, want ErrTickSize", err)
	}
	if _, err := TicksToPrice(math.MaxInt64, 1e300); !errors.Is(err, ErrOverflow) {
		t.Errorf("TicksToPrice(MaxInt64, 1e300) error = %v, want ErrOverflow", err)
	}
}

func TestRoundToTick(t *testing.T) {
	for _, c := range [][3]float64{{101.23, 0.05, 101.25}, {0.375, 0.25, 0.5}, {-0.375, 0.25, -0.5}, {99.99, 0.25, 100}} {
		if got, err := RoundToTick(c[0], c[1]); err != nil || got != c[2] {
			t.Errorf("RoundToTick(%v, %v) = (%v, %v), want %v", c[0], c[1], got, err, c[2])
		}
	}
	if _, err := RoundToTick(1, -1); !errors.Is(err, ErrTickSize) {
		t.Errorf("RoundToTick(1, -1) error = %v, want ErrTickSize", err)
	}
	if _, err := RoundToTick(math.NaN(), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("RoundToTick(NaN, 1) error = %v, want ErrNotFinite", err)
	}
	if _, err := RoundToTick(1e300, 1e-10); !errors.Is(err, ErrOverflow) {
		t.Errorf("RoundToTick(1e300, 1e-10) error = %v, want ErrOverflow", err)
	}
}

func TestTickHelpers_noAllocation(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		_, _ = PriceToTicks(101.23, 0.05)
		_, _ = TicksToPrice(2025, 0.05)
		_, _ = RoundToTick(101.23, 0.05)
	})
	if allocs != 0 {
		t.Errorf("tick helpers allocate %v times, want 0", allocs)
	}
}

func BenchmarkPriceToTicks(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PriceToTicks(101.23, 0.05)
	}
}

func BenchmarkTicksToPrice(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = TicksToPrice(2025, 0.05)
	}
}

func BenchmarkRoundToTick(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RoundToTick(101.23, 0.05)
	}
}
