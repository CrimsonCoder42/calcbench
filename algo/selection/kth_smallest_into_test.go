package selection

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example on the 16 keys of Knuth, TAOCP Vol. 3, 2nd ed., ch. 5: the 8th smallest is 509 and the
// first 8 of the work buffer are the 8 smallest keys (read off the sorted list by hand).
func TestKthSmallestInto_knuthExample(t *testing.T) {
	buf := make([]int, 0, 16)
	got, work, err := KthSmallestInto(buf, knuth16, 8)
	if err != nil || got != 509 {
		t.Fatalf("KthSmallestInto(knuth16, 8) = (%d, %v), want (509, nil)", got, err)
	}
	if len(work) != 16 || &work[0] != &buf[:1][0] {
		t.Errorf("work does not reuse dst: len %d", len(work))
	}
	smallest := slices.Sorted(slices.Values(work[:8]))
	want := []int{61, 87, 154, 170, 275, 426, 503, 509}
	if !slices.Equal(smallest, want) {
		t.Errorf("work[:8] sorted = %v, want %v", smallest, want)
	}
}

func TestKthSmallestInto_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   int
	}{
		{"one element", []int{42}, 1, 42},
		{"two elements", []int{9, 3}, 2, 9},
		{"duplicates", []int{4, 1, 4, 1, 4}, 3, 4},
		{"negative", []int{-3, 0, -7, 5, -1}, 1, -7},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, 3, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, work, err := KthSmallestInto(nil, tt.values, tt.k)
			if err != nil || got != tt.want || len(work) != len(tt.values) {
				t.Errorf("KthSmallestInto(nil, %v, %d) = (%d, len %d, %v), want %d", tt.values, tt.k, got, len(work), err, tt.want)
			}
		})
	}
}

func TestKthSmallestInto_inPlace(t *testing.T) {
	values := slices.Clone(knuth16)
	got, work, err := KthSmallestInto(values, values, 1)
	if err != nil || got != 61 || &work[0] != &values[0] {
		t.Errorf("in place: (%d, %v), want 61 in the same buffer", got, err)
	}
}

func TestKthSmallestInto_errors(t *testing.T) {
	buf := []int{9, 9, 9}
	tests := []struct {
		name    string
		values  []int
		k       int
		wantErr error
	}{
		{"empty input", nil, 1, ErrEmpty},
		{"k zero", []int{1, 2, 3}, 0, ErrRank},
		{"k above n", []int{1, 2, 3}, 4, ErrRank},
		{"k at the minimum of int", []int{1, 2, 3}, math.MinInt, ErrRank},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, work, err := KthSmallestInto(buf, tt.values, tt.k)
			if !errors.Is(err, tt.wantErr) || got != 0 || !slices.Equal(work, []int{9, 9, 9}) {
				t.Errorf("KthSmallestInto(%v, %d) = (%d, %v, %v), want (0, dst unwritten, %v)", tt.values, tt.k, got, work, err, tt.wantErr)
			}
		})
	}
}

func TestKthSmallestInto_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, _, err := KthSmallestInto(nil, values, 5); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, knuth16) {
		t.Errorf("KthSmallestInto modified its input: %v", values)
	}
}

func TestKthSmallestInto_noAllocation(t *testing.T) {
	buf := make([]int, 16)
	allocs := testing.AllocsPerRun(100, func() {
		v, _, _ := KthSmallestInto(buf, knuth16, 8)
		benchSink = v
	})
	if allocs != 0 {
		t.Errorf("KthSmallestInto with a large enough dst allocates %v times per call, want 0", allocs)
	}
}

func TestKthSmallestInto_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	var buf []int
	for trial := range 400 {
		n := 1 + rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(2*n))
		sorted := slices.Sorted(slices.Values(values))
		k := 1 + rng.IntN(n)
		got, work, err := KthSmallestInto(buf, values, k)
		if err != nil || got != sorted[k-1] {
			t.Fatalf("trial %d: KthSmallestInto(n=%d, k=%d) = (%d, %v), want %d", trial, n, k, got, err, sorted[k-1])
		}
		prefix := slices.Sorted(slices.Values(work[:k]))
		if !slices.Equal(prefix, sorted[:k]) {
			t.Fatalf("trial %d: work[:k] is not the k smallest", trial)
		}
		buf = work
	}
}

func BenchmarkKthSmallestInto(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		buf := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				v, _, _ := KthSmallestInto(buf, inputs[i&7], size.n/2)
				benchSink = v
			}
		})
	}
}
