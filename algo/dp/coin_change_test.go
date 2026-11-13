package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// coinWaysSlow enumerates the copies of each denomination (small inputs only). It returns the fewest
// coins (−1 when impossible) and the number of combinations.
func coinWaysSlow(coins []int, amount int) (fewest int, ways uint64) {
	if len(coins) == 0 {
		if amount == 0 {
			return 0, 1
		}
		return -1, 0
	}
	fewest = -1
	for copies := 0; copies*coins[0] <= amount; copies++ {
		f, w := coinWaysSlow(coins[1:], amount-copies*coins[0])
		ways += w
		if f >= 0 && (fewest < 0 || f+copies < fewest) {
			fewest = f + copies
		}
	}
	return fewest, ways
}

// Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., §7.1: there are 292 ways to change a dollar
// (100 cents) with pennies, nickels, dimes, quarters and half-dollars.
func TestCoinChangeCount_concreteMathematicsDollar(t *testing.T) {
	got, err := CoinChangeCount([]int{1, 5, 10, 25, 50}, 100)
	if err != nil || got != 292 {
		t.Errorf("CoinChangeCount(US coins, 100) = (%d, %v), want (292, nil)", got, err)
	}
	// The order of the denominations does not matter.
	if got, err := CoinChangeCount([]int{50, 1, 25, 10, 5}, 100); err != nil || got != 292 {
		t.Errorf("CoinChangeCount(shuffled, 100) = (%d, %v), want (292, nil)", got, err)
	}
}

func TestCoinChangeCount_cases(t *testing.T) {
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   uint64
	}{
		{"amount zero, no coins", nil, 0, 1},
		{"amount zero", []int{3}, 0, 1},
		{"no coins", nil, 5, 0},
		{"one coin divides", []int{2}, 6, 1},
		{"one coin does not divide", []int{2}, 5, 0},
		{"two coins", []int{1, 2}, 4, 3}, // 1111, 112, 22
		{"coin above amount", []int{10}, 3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CoinChangeCount(tt.coins, tt.amount); err != nil || got != tt.want {
				t.Errorf("CoinChangeCount(%v, %d) = (%d, %v), want %d", tt.coins, tt.amount, got, err, tt.want)
			}
		})
	}
}

func TestCoinChangeCount_errors(t *testing.T) {
	big := make([]int, 500)
	for i := range big {
		big[i] = i + 1
	}
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   error
	}{
		{"negative amount", []int{1}, -1, ErrNegative},
		{"zero coin", []int{1, 0}, 5, ErrNotPositive},
		{"negative coin", []int{-2}, 5, ErrNotPositive},
		{"duplicate", []int{5, 1, 5}, 5, ErrDuplicate},
		{"too large", []int{1}, 1 << 40, ErrTooLarge},
		{"max int amount", []int{1}, math.MaxInt, ErrTooLarge},
		{"overflow: partitions of 500 exceed 2^64", big, 500, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CoinChangeCount(tt.coins, tt.amount); !errors.Is(err, tt.want) {
				t.Errorf("CoinChangeCount error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCoinChangeCount_doesNotModifyInput(t *testing.T) {
	coins := []int{25, 1, 10, 5}
	if _, err := CoinChangeCount(coins, 30); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(coins, []int{25, 1, 10, 5}) {
		t.Errorf("input modified: %v", coins)
	}
}

// With 1, 6 and 10, greedy pays 12 as 10+1+1 (three coins); the fewest is 6+6 (two). Worked by hand.
func TestCoinChangeMin_greedyFails(t *testing.T) {
	if got, ok, err := CoinChangeMin([]int{1, 6, 10}, 12); err != nil || !ok || got != 2 {
		t.Errorf("CoinChangeMin([1 6 10], 12) = (%d, %v, %v), want (2, true, nil)", got, ok, err)
	}
}

func TestCoinChangeMin_cases(t *testing.T) {
	tests := []struct {
		name   string
		coins  []int
		amount int
		want   int
		ok     bool
	}{
		{"amount zero", nil, 0, 0, true},
		{"no coins", nil, 3, 0, false},
		{"unreachable", []int{2}, 3, 0, false},
		{"US coins, 63 cents", []int{1, 5, 10, 25}, 63, 6, true},
		{"one coin", []int{7}, 7, 1, true},
		{"duplicates allowed", []int{5, 5, 1}, 11, 3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := CoinChangeMin(tt.coins, tt.amount)
			if err != nil || ok != tt.ok || got != tt.want {
				t.Errorf("CoinChangeMin(%v, %d) = (%d, %v, %v), want (%d, %v, nil)", tt.coins, tt.amount, got, ok, err, tt.want, tt.ok)
			}
		})
	}
}

func TestCoinChangeMin_errors(t *testing.T) {
	if _, _, err := CoinChangeMin([]int{1}, -3); !errors.Is(err, ErrNegative) {
		t.Errorf("negative amount: %v", err)
	}
	if _, _, err := CoinChangeMin([]int{0}, 3); !errors.Is(err, ErrNotPositive) {
		t.Errorf("zero coin: %v", err)
	}
	if _, _, err := CoinChangeMin([]int{1}, 1<<40); !errors.Is(err, ErrTooLarge) {
		t.Errorf("huge amount: %v", err)
	}
}

func TestCoinChange_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		k := rng.IntN(4)
		coins := make([]int, 0, k)
		for len(coins) < k {
			c := 1 + rng.IntN(12)
			if !slices.Contains(coins, c) {
				coins = append(coins, c)
			}
		}
		amount := rng.IntN(40)
		wantMin, wantWays := coinWaysSlow(coins, amount)
		gotMin, ok, err := CoinChangeMin(coins, amount)
		if err != nil || ok != (wantMin >= 0) || (ok && gotMin != wantMin) {
			t.Fatalf("trial %d: CoinChangeMin(%v, %d) = (%d, %v, %v), want %d", trial, coins, amount, gotMin, ok, err, wantMin)
		}
		if got, err := CoinChangeCount(coins, amount); err != nil || got != wantWays {
			t.Fatalf("trial %d: CoinChangeCount(%v, %d) = (%d, %v), want %d", trial, coins, amount, got, err, wantWays)
		}
	}
}

var benchCoins = []int{1, 2, 5, 10, 20, 50, 100, 200, 500, 1000}

func BenchmarkCoinChangeMin(b *testing.B) {
	for _, size := range linearSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = CoinChangeMin(benchCoins, size.n)
			}
		})
	}
}

// Ten large primes keep every count up to 10⁵ below 7·10¹⁸ (checked offline), under 2⁶⁴, so the whole
// table is filled.
var benchCountCoins = []int{97, 101, 103, 107, 109, 113, 127, 131, 137, 139}

func BenchmarkCoinChangeCount(b *testing.B) {
	for _, size := range linearSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := CoinChangeCount(benchCountCoins, size.n); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
