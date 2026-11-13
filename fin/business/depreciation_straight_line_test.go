package business

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Principles of Accounting Vol. 1 §11.3, Kenzie Company: cost $58,000, salvage $10,000, 5 years →
// $9,600 a year, accumulated $48,000, final book value $10,000. Exact cents.
func TestDepreciationStraightLine_sourceExample(t *testing.T) {
	got, err := DepreciationStraightLine(5800000, 1000000, 5)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{960000, 960000, 960000, 960000, 960000}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDepreciationStraightLine_table(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		want          []int64
	}{
		{"one period", 1000, 100, 1, []int64{900}},
		{"nothing to depreciate", 500, 500, 3, []int64{0, 0, 0}},
		{"remainder in last period", 1000, 0, 3, []int64{333, 333, 334}},
		{"fewer cents than periods", 2, 0, 3, []int64{0, 0, 2}},
		{"max cost", maxExactCents, 0, 2, []int64{maxExactCents / 2, maxExactCents / 2}},
	}
	for _, c := range cases {
		got, err := DepreciationStraightLine(c.cost, c.salvage, c.life)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDepreciationStraightLine_errors(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		want          error
	}{
		{"cost above 2^53", maxExactCents + 1, 0, 5, ErrOutOfRange},
		{"salvage above 2^53", 0, maxExactCents + 1, 5, ErrOutOfRange},
		{"negative cost", -1, 0, 5, ErrNegative},
		{"negative salvage", 10, -1, 5, ErrNegative},
		{"salvage above cost", 10, 11, 5, ErrDomain},
		{"zero life", 10, 0, 0, ErrPeriods},
		{"negative life", 10, 0, -1, ErrPeriods},
		{"life above limit", 10, 0, maxLifePeriods + 1, ErrPeriods},
	}
	for _, c := range cases {
		got, err := DepreciationStraightLine(c.cost, c.salvage, c.life)
		if !errors.Is(err, c.want) || got != nil {
			t.Errorf("%s: got %v, err = %v, want nil, %v", c.name, got, err, c.want)
		}
	}
}

// Reconciliation: on 300 fixed-seed inputs the schedule sums to cost − salvage to the cent, every period
// but the last is equal, and the last differs from them by less than lifePeriods cents.
func TestDepreciationStraightLine_reconciles(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 300; i++ {
		cost := int64(rng.Uint64N(maxExactCents + 1))
		salvage := int64(rng.Uint64N(uint64(cost) + 1))
		life := 1 + rng.IntN(500)
		got, err := DepreciationStraightLine(cost, salvage, life)
		if err != nil || len(got) != life {
			t.Fatalf("(%d,%d,%d): len %d, err %v", cost, salvage, life, len(got), err)
		}
		var sum int64
		for j, v := range got {
			sum += v
			if j < life-1 && v != got[0] {
				t.Fatalf("(%d,%d,%d): period %d = %d, want %d", cost, salvage, life, j, v, got[0])
			}
		}
		if sum != cost-salvage {
			t.Fatalf("(%d,%d,%d): sum %d, want %d", cost, salvage, life, sum, cost-salvage)
		}
		if extra := got[life-1] - got[0]; extra < 0 || extra >= int64(life) {
			t.Fatalf("(%d,%d,%d): last period exceeds the others by %d", cost, salvage, life, extra)
		}
	}
}

// One allocation: the result slice.
func TestDepreciationStraightLine_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkSlice, _ = DepreciationStraightLine(5800000, 1000000, 5) }); n != 1 {
		t.Fatalf("allocs = %v, want 1", n)
	}
}

func BenchmarkDepreciationStraightLine(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkSlice, _ = DepreciationStraightLine(5800000000, 1000000, n)
			}
		})
	}
}
