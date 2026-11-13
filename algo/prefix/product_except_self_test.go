package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveProductExceptSelf is the oracle: each product computed directly in int64 (at most 8 int8
// factors, |product| <= 2⁵⁶, so int64 is exact).
func naiveProductExceptSelf(values []int8) []int64 {
	out := make([]int64, len(values))
	for i := range values {
		p := int64(1)
		for j, v := range values {
			if j != i {
				p *= int64(v)
			}
		}
		out[i] = p
	}
	return out
}

// Hand-derived from the prescan definition (Blelloch 1993, §1.1) under × with identity 1:
// exclusive prefixes 1 1 2 6, exclusive suffixes 24 12 4 1, products 24 12 8 6.
func TestProductExceptSelf_worked(t *testing.T) {
	got, err := ProductExceptSelf(nil, []int{1, 2, 3, 4})
	if want := []int{24, 12, 8, 6}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ProductExceptSelf(1 2 3 4) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestProductExceptSelf_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   []int64
	}{
		{"empty", nil, []int64{}},
		{"one: the empty product", []int64{5}, []int64{1}},
		{"two", []int64{3, -7}, []int64{-7, 3}},
		{"one zero", []int64{2, 0, 3}, []int64{0, 6, 0}},
		{"two zeros", []int64{0, 4, 0}, []int64{0, 0, 0}},
		{"negatives", []int64{-1, 2, -3}, []int64{-6, 3, -2}},
		{"duplicates", []int64{2, 2, 2}, []int64{4, 4, 4}},
		{"ones", []int64{1, -1, 1, -1}, []int64{1, -1, 1, -1}},
		{"zero beside huge values", []int64{math.MaxInt64, math.MaxInt64, 0, 0}, []int64{0, 0, 0, 0}},
		{"minimum of the type from 2⁶² × -2", []int64{1 << 62, -2, 1}, []int64{-2, 1 << 62, math.MinInt64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProductExceptSelf(nil, tt.values)
			if err != nil || len(got) != len(tt.want) || !slices.Equal(got, tt.want) {
				t.Errorf("ProductExceptSelf(%v) = (%v, %v), want (%v, nil)", tt.values, got, err, tt.want)
			}
		})
	}
	// int8: the prefix 64 × -2 = -128 is exactly the type's minimum.
	got, err := ProductExceptSelf(nil, []int8{64, -2, 1})
	if want := []int8{-2, 64, -128}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ProductExceptSelf(int8 64 -2 1) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

// The type's minimum reached from a positive prefix beyond the maximum: in 64 2 -1 -1 the prefix before
// index 2 is 128 (beyond MaxInt8), but result 2 is 64·2·-1 = -128, which fits. Hand-derived results:
// 2·-1·-1 = 2, 64·-1·-1 = 64, -128, -128.
func TestProductExceptSelf_minimumFromPositivePrefix(t *testing.T) {
	got, err := ProductExceptSelf(nil, []int8{64, 2, -1, -1})
	if want := []int8{2, 64, -128, -128}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ProductExceptSelf(int8 64 2 -1 -1) = (%v, %v), want (%v, nil)", got, err, want)
	}
	// 64 2 -1: result 2 is 64·2 = 128, which does not fit.
	if got, err := ProductExceptSelf(nil, []int8{64, 2, -1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("ProductExceptSelf(int8 64 2 -1) = (%v, %v), want ErrOverflow", got, err)
	}
}

func TestProductExceptSelf_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
	}{
		{"no zeros, product beyond 2⁶⁴", []int64{math.MaxInt64, math.MaxInt64, 3}},
		{"no zeros, just beyond MaxInt64", []int64{1 << 62, 2, 1}},
		{"one zero, the others overflow", []int64{math.MaxInt64, 2, 0}},
		{"one zero, product beyond 2⁶⁴", []int64{math.MaxInt64, math.MaxInt64, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []int64{9}
			got, err := ProductExceptSelf(dst, tt.values)
			if !errors.Is(err, ErrOverflow) || !slices.Equal(got, []int64{9}) {
				t.Errorf("ProductExceptSelf([9], %v) = (%v, %v), want ([9], ErrOverflow)", tt.values, got, err)
			}
		})
	}
	if _, err := ProductExceptSelf(nil, []uint8{16, 16, 1}); !errors.Is(err, ErrOverflow) {
		t.Errorf("uint8 16 16 1 error = %v, want ErrOverflow", err)
	}
}

func TestProductExceptSelf_unsigned(t *testing.T) {
	got, err := ProductExceptSelf(nil, []uint8{15, 17, 1})
	if want := []uint8{17, 15, 255}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ProductExceptSelf(uint8 15 17 1) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestProductExceptSelf_appendsToDst(t *testing.T) {
	got, err := ProductExceptSelf([]int{-1}, []int{2, 3})
	if want := []int{-1, 3, 2}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ProductExceptSelf([-1], [2 3]) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestProductExceptSelf_noAllocationWithBuffer(t *testing.T) {
	values := benchSigns(1000)
	buf := make([]int64, 0, len(values))
	allocs := testing.AllocsPerRun(20, func() { buf, _ = ProductExceptSelf(buf[:0], values) })
	if allocs != 0 {
		t.Errorf("ProductExceptSelf with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

// Oracle: random int8 data with zeros, ±1 and extremes against direct products; an error exactly when
// some result leaves int8.
func TestProductExceptSelf_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pool := []int8{0, 1, -1, 2, -2, 3, -3, 4, 8, 16, 64, -64, 127, -128}
	for trial := range 500 {
		values := make([]int8, rng.IntN(9))
		for i := range values {
			values[i] = pool[rng.IntN(len(pool))]
		}
		before := slices.Clone(values)
		want := naiveProductExceptSelf(values)
		fits := true
		for _, w := range want {
			if w < math.MinInt8 || w > math.MaxInt8 {
				fits = false
			}
		}
		got, err := ProductExceptSelf(nil, values)
		if !fits {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("trial %d: ProductExceptSelf(%v) = (%v, %v), want ErrOverflow (exact %v)", trial, values, got, err, want)
			}
		} else {
			if err != nil || len(got) != len(want) {
				t.Fatalf("trial %d: ProductExceptSelf(%v) = (%v, %v), want %v", trial, values, got, err, want)
			}
			for i := range want {
				if int64(got[i]) != want[i] {
					t.Fatalf("trial %d: ProductExceptSelf(%v) = %v, want %v", trial, values, got, want)
				}
			}
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

// benchSigns returns n fixed-seed values in {-1, 1} with a few 2s, so products stay in range.
func benchSigns(n int) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]int64, n)
	for i := range s {
		s[i] = int64(2*rng.IntN(2) - 1)
	}
	for i := 0; i < n && i < 40; i++ {
		s[rng.IntN(n)] = 2
	}
	return s
}

func BenchmarkProductExceptSelf(b *testing.B) {
	for _, size := range benchSizes {
		values := benchSigns(size.n)
		buf := make([]int64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = ProductExceptSelf(buf[:0], values)
			}
		})
	}
}
