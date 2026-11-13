package list

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// reverseKGroupModel is the oracle: reverse each full block of k in a slice copy.
func reverseKGroupModel(in []int, k int) []int {
	out := slices.Clone(in)
	for i := 0; i+k <= len(out); i += k {
		slices.Reverse(out[i : i+k])
	}
	return out
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-007).
func TestReverseKGroupInPlace_workedExample(t *testing.T) {
	for _, tt := range []struct {
		k    int
		want []int
	}{{2, []int{2, 1, 4, 3, 5}}, {3, []int{3, 2, 1, 4, 5}}} {
		head, err := ReverseKGroupInPlace(FromSlice([]int{1, 2, 3, 4, 5}), tt.k)
		if got := values(t, head); err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("k=%d: got (%v, %v), want %v", tt.k, got, err, tt.want)
		}
	}
}

func TestReverseKGroupInPlace_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		k    int
		want []int
	}{
		{"empty", nil, 2, nil},
		{"one", []int{1}, 1, []int{1}},
		{"k = 1 unchanged", []int{1, 2, 3}, 1, []int{1, 2, 3}},
		{"k larger than list", []int{1, 2, 3}, 4, []int{1, 2, 3}},
		{"k equals length", []int{1, 2, 3}, 3, []int{3, 2, 1}},
		{"exact blocks", []int{1, 2, 3, 4, 5, 6}, 3, []int{3, 2, 1, 6, 5, 4}},
		// Problem-site pattern, in our words: swapping neighbours is k = 2.
		{"swap pairs", []int{1, 2, 3, 4}, 2, []int{2, 1, 4, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, err := ReverseKGroupInPlace(FromSlice(tt.in), tt.k)
			if got := values(t, head); err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("got (%v, %v), want %v", got, err, tt.want)
			}
		})
	}
}

func TestReverseKGroupInPlace_errors(t *testing.T) {
	head := FromSlice([]int{1, 2})
	for _, k := range []int{0, -1} {
		if got, err := ReverseKGroupInPlace(head, k); !errors.Is(err, ErrGroupSize) || got != head {
			t.Errorf("k=%d: err = %v", k, err)
		}
	}
	cyc, _ := cyclic([]int{1, 2, 3}, 0)
	if got, err := ReverseKGroupInPlace(cyc, 2); !errors.Is(err, ErrCycle) || got != cyc {
		t.Errorf("cycle: err = %v", err)
	}
	if !DetectCycle(cyc) {
		t.Error("cycle was broken")
	}
}

func TestReverseKGroupInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randInts(rng, rng.IntN(40), 10)
		k := 1 + rng.IntN(8)
		head, err := ReverseKGroupInPlace(FromSlice(in), k)
		if got, want := values(t, head), reverseKGroupModel(in, k); err != nil || !slices.Equal(got, want) {
			t.Fatalf("(%v, k=%d) = %v, want %v", in, k, got, want)
		}
	}
}

func TestReverseKGroupInPlace_noAlloc(t *testing.T) {
	head := FromSlice(countUpTo(100))
	if a := testing.AllocsPerRun(100, func() { head, _ = ReverseKGroupInPlace(head, 3) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkReverseKGroupInPlace(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				head, _ = ReverseKGroupInPlace(head, 7)
			}
		})
	}
}
