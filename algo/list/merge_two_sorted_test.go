package list

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// tagged is a value with its origin, to check stability.
type tagged struct {
	v    int
	from byte
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-011): 1→3→5 and 1→2→4; the equal 1s keep a first.
func TestMergeTwoSortedInPlace_workedExample(t *testing.T) {
	a := FromSlice([]int{1, 3, 5})
	b := FromSlice([]int{1, 2, 4})
	aOne, bOne := a, b
	head, err := MergeTwoSortedInPlace(a, b)
	if got := values(t, head); err != nil || !slices.Equal(got, []int{1, 1, 2, 3, 4, 5}) {
		t.Fatalf("got (%v, %v)", got, err)
	}
	if head != aOne || head.Next != bOne {
		t.Error("not stable: a's 1 must come before b's 1")
	}
}

func TestMergeTwoSortedInPlace_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both empty", nil, nil, nil},
		{"a empty", nil, []int{1, 2}, []int{1, 2}},
		{"b empty", []int{1, 2}, nil, []int{1, 2}},
		{"one each", []int{2}, []int{1}, []int{1, 2}},
		{"a all before b", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"b all before a", []int{3, 4}, []int{1, 2}, []int{1, 2, 3, 4}},
		{"duplicates", []int{2, 2}, []int{2, 2}, []int{2, 2, 2, 2}},
		{"negative", []int{-5, 0}, []int{-3, 7}, []int{-5, -3, 0, 7}},
		{"extremes", []int{math.MinInt, math.MaxInt}, []int{0}, []int{math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, err := MergeTwoSortedInPlace(FromSlice(tt.a), FromSlice(tt.b))
			if got := values(t, head); err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("got (%v, %v), want %v", got, err, tt.want)
			}
		})
	}
}

func TestMergeTwoSortedInPlace_nanFirst(t *testing.T) {
	head, err := MergeTwoSortedInPlace(FromSlice([]float64{1, 2}), FromSlice([]float64{math.NaN(), 1.5}))
	got := values(t, head)
	if err != nil || len(got) != 4 || !math.IsNaN(got[0]) || !slices.Equal(got[1:], []float64{1, 1.5, 2}) {
		t.Errorf("got (%v, %v), want [NaN 1 1.5 2]", got, err)
	}
}

func TestMergeTwoSortedInPlace_errors(t *testing.T) {
	cyc, _ := cyclic([]int{1, 2}, 0)
	ok := FromSlice([]int{1})
	if h, err := MergeTwoSortedInPlace(cyc, ok); !errors.Is(err, ErrCycle) || h != nil {
		t.Errorf("cyclic a: (%v, %v)", h, err)
	}
	if h, err := MergeTwoSortedInPlace(ok, cyc); !errors.Is(err, ErrCycle) || h != nil {
		t.Errorf("cyclic b: (%v, %v)", h, err)
	}
	if got := values(t, ok); !slices.Equal(got, []int{1}) {
		t.Errorf("a changed on error: %v", got)
	}
	shared := FromSlice([]int{1, 2, 3})
	if _, err := MergeTwoSortedInPlace(shared, shared); !errors.Is(err, ErrSharedNodes) {
		t.Errorf("same list twice: err = %v", err)
	}
	if _, err := MergeTwoSortedInPlace(shared, shared.Next); !errors.Is(err, ErrSharedNodes) {
		t.Errorf("tail of a: err = %v", err)
	}
	if got := values(t, shared); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("list changed on error: %v", got)
	}
}

// Oracle: stable merge equals a stable sort of the concatenation with origin tags.
func TestMergeTwoSortedInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		a := randInts(rng, rng.IntN(20), 6)
		b := randInts(rng, rng.IntN(20), 6)
		slices.Sort(a)
		slices.Sort(b)
		var ta, tb []tagged
		for _, v := range a {
			ta = append(ta, tagged{v, 'a'})
		}
		for _, v := range b {
			tb = append(tb, tagged{v, 'b'})
		}
		want := slices.Concat(ta, tb)
		slices.SortStableFunc(want, func(x, y tagged) int { return x.v - y.v })
		// Merge the int lists, keeping a map from node to its tag.
		la, lb := FromSlice(a), FromSlice(b)
		tag := map[*Node[int]]byte{}
		for p := la; p != nil; p = p.Next {
			tag[p] = 'a'
		}
		for p := lb; p != nil; p = p.Next {
			tag[p] = 'b'
		}
		head, err := MergeTwoSortedInPlace(la, lb)
		if err != nil {
			t.Fatal(err)
		}
		var got []tagged
		for p := head; p != nil; p = p.Next {
			got = append(got, tagged{p.Value, tag[p]})
		}
		if !slices.Equal(got, want) {
			t.Fatalf("merge(%v, %v) = %v, want %v", a, b, got, want)
		}
	}
}

func BenchmarkMergeTwoSortedInPlace(b *testing.B) {
	for _, s := range benchSizes {
		evens, odds := make([]int, s.n/2), make([]int, s.n/2)
		for i := range evens {
			evens[i], odds[i] = 2*i, 2*i+1
		}
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				la, lb := FromSlice(evens), FromSlice(odds)
				b.StartTimer()
				_, _ = MergeTwoSortedInPlace(la, lb)
			}
		})
	}
}
