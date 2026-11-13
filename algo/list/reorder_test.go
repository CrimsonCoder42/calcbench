package list

import (
	"errors"
	"slices"
	"testing"
)

// reorderModel is the oracle: take from the front and the back alternately.
func reorderModel(in []int) []int {
	out := make([]int, 0, len(in))
	for i, j := 0, len(in)-1; i <= j; i, j = i+1, j-1 {
		out = append(out, in[i])
		if i != j {
			out = append(out, in[j])
		}
	}
	return out
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-013); the two shapes are also the
// problem-site pattern, in our words.
func TestReorderInPlace_workedExample(t *testing.T) {
	for _, tt := range []struct{ in, want []int }{
		{[]int{1, 2, 3, 4}, []int{1, 4, 2, 3}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 5, 2, 4, 3}},
	} {
		head, err := ReorderInPlace(FromSlice(tt.in))
		if got := values(t, head); err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("ReorderInPlace(%v) = (%v, %v), want %v", tt.in, got, err, tt.want)
		}
	}
}

func TestReorderInPlace_everyLength(t *testing.T) {
	for n := range 30 {
		in := countUpTo(n)
		head, err := ReorderInPlace(FromSlice(in))
		if got, want := values(t, head), reorderModel(in); err != nil || !slices.Equal(got, want) {
			t.Fatalf("n=%d: (%v, %v), want %v", n, got, err, want)
		}
	}
}

func TestReorderInPlace_cycle(t *testing.T) {
	cyc, entry := cyclic([]int{1, 2, 3, 4}, 2)
	if got, err := ReorderInPlace(cyc); !errors.Is(err, ErrCycle) || got != cyc {
		t.Errorf("err = %v", err)
	}
	if s, l := CycleStart(cyc); s != entry || l != 2 {
		t.Error("cycle changed")
	}
}

func TestReorderInPlace_noAlloc(t *testing.T) {
	head := FromSlice(countUpTo(100))
	if a := testing.AllocsPerRun(100, func() { head, _ = ReorderInPlace(head) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkReorderInPlace(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(countUpTo(s.n))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				head, _ = ReorderInPlace(head)
			}
		})
	}
}
