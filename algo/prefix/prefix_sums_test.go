package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// laaksonenSums is the array of Laaksonen, Competitive Programmer's Handbook (2018), §9.1 "Sum
// queries", p. 84. Shared by the PrefixSums and RangeSum tests.
var laaksonenSums = []int{1, 3, 4, 8, 6, 1, 4, 2}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// benchSink keeps benchmark results observable.
var benchSink int

// benchInts returns n fixed-seed values in [-1000, 1000).
func benchInts(n int) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]int64, n)
	for i := range s {
		s[i] = int64(rng.IntN(2000) - 1000)
	}
	return s
}

// randomInt8s returns n fixed-seed int8 values; half the time drawn near the extremes so overflow
// is common.
func randomInt8s(rng *rand.Rand, n int) []int8 {
	s := make([]int8, n)
	for i := range s {
		if rng.IntN(2) == 0 {
			s[i] = int8(rng.IntN(256) - 128)
		} else {
			s[i] = int8(rng.IntN(21) - 10)
		}
	}
	return s
}

// Laaksonen §9.1, p. 84 prints the inclusive prefix sums 1 4 8 16 22 23 27 29; ours start with 0.
func TestPrefixSums_laaksonenExample(t *testing.T) {
	want := []int{0, 1, 4, 8, 16, 22, 23, 27, 29}
	got, err := PrefixSums(nil, laaksonenSums)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("PrefixSums(laaksonen) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestPrefixSums_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   []int64
	}{
		{"empty", nil, []int64{0}},
		{"one", []int64{7}, []int64{0, 7}},
		{"two", []int64{7, -2}, []int64{0, 7, 5}},
		{"duplicates", []int64{3, 3, 3}, []int64{0, 3, 6, 9}},
		{"negative and zero", []int64{-1, 0, -2}, []int64{0, -1, -1, -3}},
		{"extremes that cancel", []int64{math.MaxInt64, math.MinInt64, 1}, []int64{0, math.MaxInt64, -1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrefixSums(nil, tt.values)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("PrefixSums(%v) = (%v, %v), want (%v, nil)", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPrefixSums_unsigned(t *testing.T) {
	got, err := PrefixSums(nil, []uint8{200, 55})
	if want := []uint8{0, 200, 255}; err != nil || !slices.Equal(got, want) {
		t.Errorf("PrefixSums(uint8 200 55) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestPrefixSums_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []int8
	}{
		{"above the maximum", []int8{127, 1}},
		{"below the minimum", []int8{-128, -1}},
		{"overflow late", []int8{1, 2, 3, 100, 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []int8{9}
			got, err := PrefixSums(dst, tt.values)
			if !errors.Is(err, ErrOverflow) || !slices.Equal(got, []int8{9}) {
				t.Errorf("PrefixSums([9], %v) = (%v, %v), want ([9], ErrOverflow)", tt.values, got, err)
			}
		})
	}
	if _, err := PrefixSums(nil, []uint8{255, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("PrefixSums(uint8 255 1) error = %v, want ErrOverflow", err)
	}
}

func TestPrefixSums_appendsToDst(t *testing.T) {
	got, err := PrefixSums([]int{-1}, []int{2, 3})
	if want := []int{-1, 0, 2, 5}; err != nil || !slices.Equal(got, want) {
		t.Errorf("PrefixSums([-1], [2 3]) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestPrefixSums_noAllocationWithBuffer(t *testing.T) {
	values := benchInts(1000)
	buf := make([]int64, 0, len(values)+1)
	allocs := testing.AllocsPerRun(20, func() { buf, _ = PrefixSums(buf[:0], values) })
	if allocs != 0 {
		t.Errorf("PrefixSums with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

// Oracle: int8 prefix sums against int sums; an error exactly when some prefix leaves [-128, 127].
func TestPrefixSums_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInt8s(rng, rng.IntN(12))
		before := slices.Clone(values)
		want := []int8{0}
		overflow := false
		sum := 0
		for _, v := range values {
			sum += int(v)
			if sum < math.MinInt8 || sum > math.MaxInt8 {
				overflow = true
				break
			}
			want = append(want, int8(sum))
		}
		got, err := PrefixSums(nil, values)
		if overflow {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("trial %d: PrefixSums(%v) = (%v, %v), want ErrOverflow", trial, values, got, err)
			}
		} else if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: PrefixSums(%v) = (%v, %v), want %v", trial, values, got, err, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: PrefixSums modified its input", trial)
		}
	}
}

func BenchmarkPrefixSums(b *testing.B) {
	for _, size := range benchSizes {
		values := benchInts(size.n)
		buf := make([]int64, 0, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = PrefixSums(buf[:0], values)
			}
		})
	}
}
