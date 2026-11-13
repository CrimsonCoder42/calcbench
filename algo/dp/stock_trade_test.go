package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// oneTradeSlow tries every pair: O(n²).
func oneTradeSlow(prices []int64) int64 {
	var best int64
	for b := range prices {
		for s := b + 1; s < len(prices); s++ {
			best = max(best, prices[s]-prices[b])
		}
	}
	return best
}

// unlimitedSlow: best(i) = profit from day i on with nothing held = max(best(i+1), max over b ≥ i,
// s > b of p[s] − p[b] + best(s)) — every set of non-overlapping trades, by recursion (n ≤ 10).
func unlimitedSlow(prices []int64, from int) int64 {
	var best int64
	for b := from; b < len(prices); b++ {
		for s := b + 1; s < len(prices); s++ {
			best = max(best, prices[s]-prices[b]+unlimitedSlow(prices, s))
		}
	}
	return best
}

// Bentley, Programming Pearls, §8.4 (scanning): worked by hand on [9 4 6 2 8 5]: buy at 2 (index 3),
// sell at 8 (index 4), profit 6; with unlimited trades 4→6 and 2→8 earn 2 + 6 = 8.
func TestMaxProfit_example(t *testing.T) {
	prices := []int64{9, 4, 6, 2, 8, 5}
	if p, b, s, err := MaxProfitOneTrade(prices); err != nil || p != 6 || b != 3 || s != 4 {
		t.Errorf("MaxProfitOneTrade = (%d, %d, %d, %v), want (6, 3, 4, nil)", p, b, s, err)
	}
	if p, err := MaxProfitUnlimitedTrades(prices); err != nil || p != 8 {
		t.Errorf("MaxProfitUnlimitedTrades = (%d, %v), want (8, nil)", p, err)
	}
}

func TestMaxProfitOneTrade_cases(t *testing.T) {
	tests := []struct {
		name            string
		prices          []int64
		profit          int64
		buy, sell       int
		unlimitedProfit int64
	}{
		{"empty", nil, 0, -1, -1, 0},
		{"one", []int64{5}, 0, -1, -1, 0},
		{"two rising", []int64{3, 7}, 4, 0, 1, 4},
		{"two falling", []int64{7, 3}, 0, -1, -1, 0},
		{"flat", []int64{4, 4, 4}, 0, -1, -1, 0},
		{"sorted", []int64{1, 2, 3, 4}, 3, 0, 3, 3},
		{"reversed", []int64{4, 3, 2, 1}, 0, -1, -1, 0},
		{"tie: smallest sell, earliest low", []int64{1, 5, 1, 5}, 4, 0, 1, 8},
		{"zero price", []int64{0, 9}, 9, 0, 1, 9},
		{"extreme", []int64{0, math.MaxInt64}, math.MaxInt64, 0, 1, math.MaxInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, b, s, err := MaxProfitOneTrade(tt.prices)
			if err != nil || p != tt.profit || b != tt.buy || s != tt.sell {
				t.Errorf("MaxProfitOneTrade(%v) = (%d, %d, %d, %v), want (%d, %d, %d)", tt.prices, p, b, s, err, tt.profit, tt.buy, tt.sell)
			}
			if p, err := MaxProfitUnlimitedTrades(tt.prices); err != nil || p != tt.unlimitedProfit {
				t.Errorf("MaxProfitUnlimitedTrades(%v) = (%d, %v), want %d", tt.prices, p, err, tt.unlimitedProfit)
			}
		})
	}
}

func TestMaxProfit_errors(t *testing.T) {
	if _, _, _, err := MaxProfitOneTrade([]int64{3, -1}); !errors.Is(err, ErrNegative) {
		t.Errorf("MaxProfitOneTrade negative: %v", err)
	}
	if _, err := MaxProfitUnlimitedTrades([]int64{3, -1}); !errors.Is(err, ErrNegative) {
		t.Errorf("MaxProfitUnlimitedTrades negative: %v", err)
	}
	if _, err := MaxProfitUnlimitedTrades([]int64{0, math.MaxInt64, 0, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxProfitUnlimitedTrades overflow: %v", err)
	}
}

func TestMaxProfit_allocations(t *testing.T) {
	prices := []int64{9, 4, 6, 2, 8, 5}
	if n := testing.AllocsPerRun(100, func() {
		_, _, _, _ = MaxProfitOneTrade(prices)
		_, _ = MaxProfitUnlimitedTrades(prices)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func TestMaxProfit_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		prices := benchInt64sFrom(rng, rng.IntN(10), 20)
		p, b, s, err := MaxProfitOneTrade(prices)
		if want := oneTradeSlow(prices); err != nil || p != want {
			t.Fatalf("trial %d: MaxProfitOneTrade(%v) = %d, want %d", trial, prices, p, want)
		}
		if p > 0 && (b >= s || prices[s]-prices[b] != p) {
			t.Fatalf("trial %d: indices (%d, %d) do not give %d", trial, b, s, p)
		}
		if got, want := must(MaxProfitUnlimitedTrades(prices)), unlimitedSlow(prices, 0); got != want {
			t.Fatalf("trial %d: MaxProfitUnlimitedTrades(%v) = %d, want %d", trial, prices, got, want)
		}
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func BenchmarkMaxProfitOneTrade(b *testing.B) {
	for _, size := range linearSizes {
		prices := benchInt64s(size.n, 10_000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _, _ = MaxProfitOneTrade(prices)
			}
		})
	}
}

func BenchmarkMaxProfitUnlimitedTrades(b *testing.B) {
	for _, size := range linearSizes {
		prices := benchInt64s(size.n, 10_000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MaxProfitUnlimitedTrades(prices)
			}
		})
	}
}
