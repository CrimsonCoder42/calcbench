package sorting

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// CLRS 3rd ed. §8.3 (Figure 8.3) illustrates LSD radix sort on seven 3-digit numbers. Our own seven
// 3-digit numbers; the expected output is hand-derived.
func TestRadixSortUint64_clrsShapedExample(t *testing.T) {
	got := RadixSortUint64(nil, []uint64{329, 457, 657, 839, 436, 720, 355})
	want := []uint64{329, 355, 436, 457, 657, 720, 839}
	if !slices.Equal(got, want) {
		t.Errorf("RadixSortUint64 = %v, want %v", got, want)
	}
}

func TestRadixSortUint64_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []uint64
		want   []uint64
	}{
		{"empty", []uint64{}, []uint64{}},
		{"one element", []uint64{7}, []uint64{7}},
		{"two elements reversed", []uint64{2, 1}, []uint64{1, 2}},
		{"all equal (every pass skipped)", []uint64{5, 5, 5}, []uint64{5, 5, 5}},
		{"differ only in the top byte", []uint64{3 << 56, 1 << 56, 2 << 56}, []uint64{1 << 56, 2 << 56, 3 << 56}},
		{"zero and the maximum", []uint64{math.MaxUint64, 0, 1}, []uint64{0, 1, math.MaxUint64}},
		{"duplicates", []uint64{9, 1, 9, 1}, []uint64{1, 1, 9, 9}},
		{"reversed", []uint64{1 << 40, 1 << 20, 1 << 10, 1}, []uint64{1, 1 << 10, 1 << 20, 1 << 40}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RadixSortUint64(nil, tt.values); !slices.Equal(got, tt.want) {
				t.Errorf("RadixSortUint64(%v) = %v, want %v", tt.values, got, tt.want)
			}
		})
	}
}

func TestRadixSortUint64_inPlaceAndInputKept(t *testing.T) {
	values := []uint64{300, 2, 70000, 1}
	RadixSortUint64(nil, values)
	if !slices.Equal(values, []uint64{300, 2, 70000, 1}) {
		t.Errorf("RadixSortUint64 modified its input: %v", values)
	}
	RadixSortUint64(values, values)
	if !slices.Equal(values, []uint64{1, 2, 300, 70000}) {
		t.Errorf("in place: %v", values)
	}
}

// Odd and even numbers of non-skipped passes leave the data in either buffer; both must end in dst.
func TestRadixSortUint64_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(300)
		bits := 1 + rng.IntN(64)
		values := make([]uint64, n)
		for i := range values {
			values[i] = rng.Uint64() >> (64 - bits)
		}
		want := slices.Sorted(slices.Values(values))
		if got := RadixSortUint64(nil, values); !slices.Equal(got, want) {
			t.Fatalf("trial %d (bits %d): RadixSortUint64 differs from slices.Sort", trial, bits)
		}
	}
}

func BenchmarkRadixSortUint64(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		values := make([]uint64, size.n)
		for i := range values {
			values[i] = rng.Uint64()
		}
		dst := make([]uint64, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = int(RadixSortUint64(dst, values)[0] & 1)
			}
		})
	}
}
