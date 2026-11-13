package dp

import (
	"errors"
	"math"
	"testing"
)

// Knuth, TAOCP Vol. 1, 3rd ed., §1.2.8: F(0) … F(12) = 0 1 1 2 3 5 8 13 21 34 55 89 144.
func TestFibonacci_knuthTable(t *testing.T) {
	want := []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	for n, w := range want {
		if got, err := Fibonacci(n); err != nil || got != w {
			t.Errorf("Fibonacci(%d) = (%d, %v), want %d", n, got, err, w)
		}
	}
}

func TestFibonacci_largeValues(t *testing.T) {
	tests := []struct {
		n    int
		want uint64
	}{
		{50, 12586269025},
		{92, 7540113804746346429},  // the largest that fits int64
		{93, 12200160415121876738}, // the largest that fits uint64
	}
	for _, tt := range tests {
		if got, err := Fibonacci(tt.n); err != nil || got != tt.want {
			t.Errorf("Fibonacci(%d) = (%d, %v), want %d", tt.n, got, err, tt.want)
		}
	}
}

func TestFibonacci_errors(t *testing.T) {
	for _, n := range []int{94, 1000, math.MaxInt} {
		if _, err := Fibonacci(n); !errors.Is(err, ErrOverflow) {
			t.Errorf("Fibonacci(%d) error = %v, want ErrOverflow", n, err)
		}
	}
	for _, n := range []int{-1, math.MinInt} {
		if _, err := Fibonacci(n); !errors.Is(err, ErrNegative) {
			t.Errorf("Fibonacci(%d) error = %v, want ErrNegative", n, err)
		}
	}
}

// Erickson, Algorithms, §3.1: the meters (sequences of 1s and 2s) of total length n number F(n+1).
// Oracle: count the sequences directly for n ≤ 20; this is the stair-climbing count.
func TestFibonacci_climbingStairsOracle(t *testing.T) {
	var count func(n int) uint64
	count = func(n int) uint64 {
		if n <= 1 {
			return 1
		}
		return count(n-1) + count(n-2)
	}
	for n := range 21 {
		if got, err := Fibonacci(n + 1); err != nil || got != count(n) {
			t.Errorf("Fibonacci(%d) = (%d, %v), want %d sequences of 1s and 2s summing to %d", n+1, got, err, count(n), n)
		}
	}
}

func TestFibonacci_allocations(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { _, _ = Fibonacci(90) }); n != 0 {
		t.Errorf("Fibonacci allocates %v times, want 0", n)
	}
}

// Scalar function with n ≤ 93: one case at the largest n.
func BenchmarkFibonacci(b *testing.B) {
	b.Run("93", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Fibonacci(93)
		}
	})
}
