package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Laaksonen §9.1, p. 84: sum(3, 6) = 8 + 6 + 1 + 4 = 19 = 27 − 8. Closed [3, 6] is our [3, 7).
func TestRangeSum_laaksonenExample(t *testing.T) {
	prefix, err := PrefixSums(nil, laaksonenSums)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := RangeSum(prefix, 3, 7); err != nil || got != 19 {
		t.Errorf("RangeSum(laaksonen, 3, 7) = (%d, %v), want (19, nil)", got, err)
	}
}

func TestRangeSum_cases(t *testing.T) {
	prefix, _ := PrefixSums(nil, []int{5, -2, 0, 7})
	tests := []struct {
		name   string
		lo, hi int
		want   int
	}{
		{"empty at the start", 0, 0, 0},
		{"empty at the end", 4, 4, 0},
		{"one", 1, 2, -2},
		{"zero element", 2, 3, 0},
		{"all", 0, 4, 10},
		{"middle", 1, 4, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RangeSum(prefix, tt.lo, tt.hi); err != nil || got != tt.want {
				t.Errorf("RangeSum(%d, %d) = (%d, %v), want (%d, nil)", tt.lo, tt.hi, got, err, tt.want)
			}
		})
	}
}

func TestRangeSum_errors(t *testing.T) {
	prefix := []int{0, 1, 3}
	for _, r := range [][2]int{{-1, 1}, {2, 1}, {0, 3}, {math.MinInt, 0}, {0, math.MaxInt}} {
		if _, err := RangeSum(prefix, r[0], r[1]); !errors.Is(err, ErrRange) {
			t.Errorf("RangeSum(%d, %d) error = %v, want ErrRange", r[0], r[1], err)
		}
	}
	if _, err := RangeSum([]int{}, 0, 0); !errors.Is(err, ErrRange) {
		t.Errorf("RangeSum(empty prefix, 0, 0) error = %v, want ErrRange", err)
	}
	// Every prefix fits int8 (0, -128, -1, 126) but values[1:3] sums to 254.
	p8, err := PrefixSums(nil, []int8{-128, 127, 127})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := RangeSum(p8, 1, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("RangeSum(int8, 1, 3) error = %v, want ErrOverflow", err)
	}
}

func TestRangeSum_noAllocation(t *testing.T) {
	prefix, _ := PrefixSums(nil, benchInts(1000))
	allocs := testing.AllocsPerRun(20, func() { _, _ = RangeSum(prefix, 10, 900) })
	if allocs != 0 {
		t.Errorf("RangeSum allocates %v times per call, want 0", allocs)
	}
}

// Oracle: every range of random int8 data against a direct int sum.
func TestRangeSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := make([]int8, rng.IntN(10))
		for i := range values {
			values[i] = int8(rng.IntN(41) - 20)
		}
		prefix, err := PrefixSums(nil, values)
		if err != nil {
			t.Fatal(err)
		}
		before := slices.Clone(prefix)
		for lo := 0; lo <= len(values); lo++ {
			for hi := lo; hi <= len(values); hi++ {
				want := 0
				for _, v := range values[lo:hi] {
					want += int(v)
				}
				got, err := RangeSum(prefix, lo, hi)
				if err != nil || int(got) != want {
					t.Fatalf("trial %d: RangeSum(%v, %d, %d) = (%d, %v), want %d", trial, values, lo, hi, got, err, want)
				}
			}
		}
		if !slices.Equal(prefix, before) {
			t.Fatalf("trial %d: RangeSum modified prefix", trial)
		}
	}
}

func BenchmarkRangeSum(b *testing.B) {
	for _, size := range benchSizes {
		prefix, _ := PrefixSums(nil, benchInts(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var sum int64
			for b.Loop() {
				for lo := 0; lo < size.n; lo += 7 {
					s, _ := RangeSum(prefix, lo, size.n)
					sum += s
				}
			}
			benchSink = int(sum)
		})
	}
}
