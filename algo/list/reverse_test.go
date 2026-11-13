package list

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example (Knuth Vol. 1 §2.2.3 Ex. 7, hand-derived): 1→2→3→4 inverted is 4→3→2→1.
func TestReverseInPlace_workedExample(t *testing.T) {
	head, err := ReverseInPlace(FromSlice([]int{1, 2, 3, 4}))
	if err != nil || !slices.Equal(values(t, head), []int{4, 3, 2, 1}) {
		t.Errorf("ReverseInPlace = (%v, %v)", values(t, head), err)
	}
}

func TestReverseInPlace_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
	}{
		{"empty", nil},
		{"one", []int{1}},
		{"two", []int{1, 2}},
		{"duplicates", []int{2, 2, 1, 1}},
		{"negative and zero", []int{-3, 0, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, err := ReverseInPlace(FromSlice(tt.in))
			want := slices.Clone(tt.in)
			slices.Reverse(want)
			if got := values(t, head); err != nil || !slices.Equal(got, want) {
				t.Errorf("ReverseInPlace(%v) = (%v, %v), want %v", tt.in, got, err, want)
			}
		})
	}
}

// Problem-site pattern, in our words: five nodes come back in the opposite order.
func TestReverseInPlace_fiveNodes(t *testing.T) {
	head, _ := ReverseInPlace(FromSlice([]int{1, 2, 3, 4, 5}))
	if got := values(t, head); !slices.Equal(got, []int{5, 4, 3, 2, 1}) {
		t.Errorf("got %v", got)
	}
}

func TestReverseInPlace_cycleLeftUnchanged(t *testing.T) {
	head, entry := cyclic([]int{1, 2, 3, 4}, 1)
	got, err := ReverseInPlace(head)
	if !errors.Is(err, ErrCycle) || got != head {
		t.Fatalf("ReverseInPlace(cycle) = (%p, %v), want (head, ErrCycle)", got, err)
	}
	if start, n := CycleStart(head); start != entry || n != 3 {
		t.Errorf("cycle changed: start %v len %d", start, n)
	}
}

func TestReverseInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randInts(rng, rng.IntN(40), 10)
		want := slices.Clone(in)
		slices.Reverse(want)
		head, err := ReverseInPlace(FromSlice(in))
		if got := values(t, head); err != nil || !slices.Equal(got, want) {
			t.Fatalf("ReverseInPlace(%v) = %v, want %v", in, got, want)
		}
	}
}

func TestReverseInPlace_noAlloc(t *testing.T) {
	head := FromSlice(countUpTo(100))
	if a := testing.AllocsPerRun(100, func() { head, _ = ReverseInPlace(head) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkReverseInPlace(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				head, _ = ReverseInPlace(head)
			}
		})
	}
}
