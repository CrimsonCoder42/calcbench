package hashing

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func groupsEqual[K comparable, T comparable](a, b []Group[K, T]) bool {
	return slices.EqualFunc(a, b, func(x, y Group[K, T]) bool { return x.Key == y.Key && slices.Equal(x.Items, y.Items) })
}

// Hand-derived: 1…7 by parity; the odd group comes first because 1 is first.
func TestGroupBy_parity(t *testing.T) {
	got, err := GroupBy([]int{1, 2, 3, 4, 5, 6, 7}, func(v int) bool { return v%2 == 0 })
	want := []Group[bool, int]{{false, []int{1, 3, 5, 7}}, {true, []int{2, 4, 6}}}
	if err != nil || !groupsEqual(got, want) {
		t.Errorf("GroupBy(parity) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestGroupBy_cases(t *testing.T) {
	byLen := func(s string) int { return len(s) }
	tests := []struct {
		name   string
		values []string
		want   []Group[int, string]
	}{
		{"empty", nil, []Group[int, string]{}},
		{"one", []string{"ab"}, []Group[int, string]{{2, []string{"ab"}}}},
		{"first appearance order", []string{"ccc", "a", "bb", "dd", "e"},
			[]Group[int, string]{{3, []string{"ccc"}}, {1, []string{"a", "e"}}, {2, []string{"bb", "dd"}}}},
		{"duplicates kept", []string{"x", "x"}, []Group[int, string]{{1, []string{"x", "x"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GroupBy(tt.values, byLen)
			if err != nil || got == nil || !groupsEqual(got, tt.want) {
				t.Errorf("GroupBy(%v) = (%#v, %v), want (%v, nil)", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestGroupBy_nanKeys(t *testing.T) {
	got, err := GroupBy([]float64{1, 2, 3}, func(float64) float64 { return math.NaN() })
	if err != nil || len(got) != 3 {
		t.Errorf("GroupBy(NaN keys) = (%v, %v), want 3 singleton groups", got, err)
	}
}

func TestGroupBy_errors(t *testing.T) {
	if got, err := GroupBy[int, int]([]int{1}, nil); !errors.Is(err, ErrNilKeyFunc) || got != nil {
		t.Errorf("GroupBy(nil key) = (%v, %v), want (nil, ErrNilKeyFunc)", got, err)
	}
}

func TestGroupBy_keyCalledOncePerValue(t *testing.T) {
	calls := 0
	if _, err := GroupBy([]int{5, 6, 7}, func(v int) int { calls++; return v % 2 }); err != nil || calls != 3 {
		t.Errorf("key called %d times (err %v), want 3", calls, err)
	}
}

// Oracle: for each distinct key in first-appearance order, filter the input.
func TestGroupBy_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	mod := func(v int) int { return v % 7 }
	for trial := range 300 {
		values := randomInts(rng, rng.IntN(200), 1000)
		var want []Group[int, int]
		for _, v := range values {
			k := mod(v)
			if slices.ContainsFunc(want, func(g Group[int, int]) bool { return g.Key == k }) {
				continue
			}
			g := Group[int, int]{Key: k}
			for _, w := range values {
				if mod(w) == k {
					g.Items = append(g.Items, w)
				}
			}
			want = append(want, g)
		}
		got, err := GroupBy(values, mod)
		if err != nil || !groupsEqual(got, want) {
			t.Fatalf("trial %d: GroupBy mismatch", trial)
		}
	}
}

func BenchmarkGroupBy(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n)
		key := func(v int) int { return v % 100 }
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = GroupBy(values, key)
			}
		})
	}
}
