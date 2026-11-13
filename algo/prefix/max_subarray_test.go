package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveMaxSubarray is the oracle: every subarray in int, ends ascending and, for each end, starts
// descending (shortest first), keeping only strictly larger sums: the tie rule MaxSubarray documents.
func naiveMaxSubarray(values []int8) (best, start, end int) {
	best = math.MinInt
	for e := 1; e <= len(values); e++ {
		sum := 0
		for s := e - 1; s >= 0; s-- {
			sum += int(values[s])
			if sum > best {
				best, start, end = sum, s, e
			}
		}
	}
	return best, start, end
}

// Bentley, "Programming Pearls: Algorithm Design Techniques", CACM 27(9), 1984 (and Programming Pearls
// 2nd ed., Column 8): the ten-element example vector, maximum 187 from x[2..6] (0-based).
func TestMaxSubarray_bentleyExample(t *testing.T) {
	values := []int{31, -41, 59, 26, -53, 58, 97, -93, -23, 84}
	sum, start, end, err := MaxSubarray(values)
	if err != nil || sum != 187 || start != 2 || end != 7 {
		t.Errorf("MaxSubarray(bentley) = (%d, %d, %d, %v), want (187, 2, 7, nil)", sum, start, end, err)
	}
}

// Laaksonen, Competitive Programmer's Handbook (2018), §2.4, p. 21: −1 2 4 −3 5 2 −5 2 has maximum
// subarray sum 10, from 2 4 −3 5 2.
func TestMaxSubarray_laaksonenExample(t *testing.T) {
	sum, start, end, err := MaxSubarray([]int{-1, 2, 4, -3, 5, 2, -5, 2})
	if err != nil || sum != 10 || start != 1 || end != 6 {
		t.Errorf("MaxSubarray(laaksonen) = (%d, %d, %d, %v), want (10, 1, 6, nil)", sum, start, end, err)
	}
}

func TestMaxSubarray_cases(t *testing.T) {
	tests := []struct {
		name       string
		values     []int64
		sum        int64
		start, end int
	}{
		{"one", []int64{-5}, -5, 0, 1},
		{"two, second larger", []int64{-5, -2}, -2, 1, 2},
		{"all negative", []int64{-3, -1, -2}, -1, 1, 2},
		{"leading zero dropped (shortest)", []int64{0, 1}, 1, 1, 2},
		{"tie: the one that ends first", []int64{1, -1, 1}, 1, 0, 1},
		{"all zeros", []int64{0, 0, 0}, 0, 0, 1},
		{"sorted", []int64{-2, -1, 0, 1, 2}, 3, 3, 5},
		{"reversed", []int64{2, 1, 0, -1, -2}, 3, 0, 2},
		{"duplicates", []int64{4, 4, 4}, 12, 0, 3},
		{"maximum of the type", []int64{math.MaxInt64, -1, 1}, math.MaxInt64, 0, 1},
		{"minimum of the type", []int64{math.MinInt64, math.MinInt64}, math.MinInt64, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum, start, end, err := MaxSubarray(tt.values)
			if err != nil || sum != tt.sum || start != tt.start || end != tt.end {
				t.Errorf("MaxSubarray(%v) = (%d, %d, %d, %v), want (%d, %d, %d, nil)", tt.values, sum, start, end, err, tt.sum, tt.start, tt.end)
			}
		})
	}
	if sum, start, end, err := MaxSubarray([]uint8{1, 0, 2}); err != nil || sum != 3 || start != 0 || end != 3 {
		t.Errorf("MaxSubarray(uint8 1 0 2) = (%d, %d, %d, %v), want (3, 0, 3, nil)", sum, start, end, err)
	}
}

func TestMaxSubarray_errors(t *testing.T) {
	if _, _, _, err := MaxSubarray([]int{}); !errors.Is(err, ErrEmpty) {
		t.Errorf("MaxSubarray(empty) error = %v, want ErrEmpty", err)
	}
	if _, _, _, err := MaxSubarray([]int8{100, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxSubarray(int8 100 100) error = %v, want ErrOverflow", err)
	}
	if _, _, _, err := MaxSubarray([]uint8{200, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("MaxSubarray(uint8 200 100) error = %v, want ErrOverflow", err)
	}
}

func TestMaxSubarray_noAllocation(t *testing.T) {
	values := benchInts(1000)
	allocs := testing.AllocsPerRun(20, func() { _, _, _, _ = MaxSubarray(values) })
	if allocs != 0 {
		t.Errorf("MaxSubarray allocates %v times per call, want 0", allocs)
	}
}

// Oracle: random int8 data against the O(n²) scan; an error exactly when the true maximum exceeds 127.
func TestMaxSubarray_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInt8s(rng, 1+rng.IntN(15))
		before := slices.Clone(values)
		wantSum, wantStart, wantEnd := naiveMaxSubarray(values)
		sum, start, end, err := MaxSubarray(values)
		if wantSum > math.MaxInt8 {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("trial %d: MaxSubarray(%v) = (%d, %d, %d, %v), want ErrOverflow", trial, values, sum, start, end, err)
			}
		} else if err != nil || int(sum) != wantSum || start != wantStart || end != wantEnd {
			t.Fatalf("trial %d: MaxSubarray(%v) = (%d, %d, %d, %v), want (%d, %d, %d)", trial, values, sum, start, end, err, wantSum, wantStart, wantEnd)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkMaxSubarray(b *testing.B) {
	for _, size := range benchSizes {
		values := benchInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, start, _, _ := MaxSubarray(values)
				benchSink = start
			}
		})
	}
}
