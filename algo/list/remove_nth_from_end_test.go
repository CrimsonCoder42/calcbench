package list

import (
	"errors"
	"slices"
	"testing"
)

// Worked example (hand-derived, catalog/list.md ALG-LIST-012): remove the 2nd from the end of 1…5.
// It is also the problem-site pattern, in our words.
func TestRemoveNthFromEndInPlace_workedExample(t *testing.T) {
	head := FromSlice([]int{1, 2, 3, 4, 5})
	removed := nodeAt(head, 3)
	got, err := RemoveNthFromEndInPlace(head, 2)
	if err != nil || !slices.Equal(values(t, got), []int{1, 2, 3, 5}) {
		t.Fatalf("got (%v, %v)", values(t, got), err)
	}
	if removed.Next != nil {
		t.Error("removed node still linked")
	}
}

// Every position of every length up to 12, against a slice model.
func TestRemoveNthFromEndInPlace_allPositions(t *testing.T) {
	for n := 1; n <= 12; n++ {
		for k := 1; k <= n; k++ {
			in := countUpTo(n)
			want := slices.Delete(slices.Clone(in), n-k, n-k+1)
			head, err := RemoveNthFromEndInPlace(FromSlice(in), k)
			if got := values(t, head); err != nil || !slices.Equal(got, want) {
				t.Fatalf("n=%d k=%d: (%v, %v), want %v", n, k, got, err, want)
			}
		}
	}
}

func TestRemoveNthFromEndInPlace_onlyNode(t *testing.T) {
	if head, err := RemoveNthFromEndInPlace(FromSlice([]int{9}), 1); head != nil || err != nil {
		t.Errorf("(%v, %v), want (nil, nil)", head, err)
	}
}

func TestRemoveNthFromEndInPlace_errors(t *testing.T) {
	head := FromSlice([]int{1, 2, 3})
	for _, k := range []int{0, -1, 4} {
		if got, err := RemoveNthFromEndInPlace(head, k); !errors.Is(err, ErrPosition) || got != head {
			t.Errorf("k=%d: err = %v", k, err)
		}
	}
	if _, err := RemoveNthFromEndInPlace[int](nil, 1); !errors.Is(err, ErrPosition) {
		t.Errorf("nil: err = %v", err)
	}
	if got := values(t, head); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("list changed on error: %v", got)
	}
	cyc, _ := cyclic([]int{1, 2, 3}, 1)
	if _, err := RemoveNthFromEndInPlace(cyc, 1); !errors.Is(err, ErrCycle) {
		t.Errorf("cycle: err = %v", err)
	}
}

func BenchmarkRemoveNthFromEndInPlace(b *testing.B) {
	for _, s := range benchSizes {
		in := countUpTo(s.n)
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				head := FromSlice(in)
				b.StartTimer()
				_, _ = RemoveNthFromEndInPlace(head, s.n/3)
			}
		})
	}
}
