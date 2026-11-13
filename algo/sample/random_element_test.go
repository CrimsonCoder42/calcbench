package sample

import (
	"errors"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
)

// listNode is a minimal singly linked list, our own, to show RandomElement over a list of unknown length.
type listNode struct {
	val  int
	next *listNode
}

func (n *listNode) all() iter.Seq[int] {
	return func(yield func(int) bool) {
		for p := n; p != nil; p = p.next {
			if !yield(p.val) {
				return
			}
		}
	}
}

// Vitter 1985, Algorithm R with k = 1: a 5-node list, 50 000 draws, each value within 5σ of 10 000.
func TestRandomElement_vitterUniformList(t *testing.T) {
	var head *listNode
	for v := 5; v >= 1; v-- {
		head = &listNode{v, head}
	}
	rng := rand.New(rand.NewPCG(1, 2))
	const trials = 50_000
	counts := map[int]int{}
	for range trials {
		v, err := RandomElement(rng, head.all())
		if err != nil {
			t.Fatal(err)
		}
		counts[v]++
	}
	for v := 1; v <= 5; v++ {
		if !within5Sigma(counts[v], trials, 0.2) {
			t.Errorf("value %d: %d of %d", v, counts[v], trials)
		}
	}
}

func TestRandomElement_cases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	if v, err := RandomElement(rng, slices.Values([]string{"only"})); err != nil || v != "only" {
		t.Errorf("one item: %q, %v", v, err)
	}
	if v, err := RandomElement(rng, slices.Values([]int{4, 4, 4})); err != nil || v != 4 {
		t.Errorf("duplicates: %d, %v", v, err)
	}
}

func TestRandomElement_errors(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	if _, err := RandomElement(rng, slices.Values([]int(nil))); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: err = %v", err)
	}
	if _, err := RandomElement[int](rng, nil); !errors.Is(err, ErrNilSeq) {
		t.Errorf("nil seq: err = %v", err)
	}
	if _, err := RandomElement(nil, slices.Values([]int{1})); !errors.Is(err, ErrNilRand) {
		t.Errorf("nil rng: err = %v", err)
	}
}

// TestRandomElement_oracle: for random lengths n, each position is returned with frequency 1/n.
func TestRandomElement_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for n := 1; n <= 8; n++ {
		items := make([]int, n)
		for i := range items {
			items[i] = i
		}
		const trials = 16_000
		counts := make([]int, n)
		for range trials {
			v, _ := RandomElement(rng, slices.Values(items))
			counts[v]++
		}
		for i, c := range counts {
			if !within5Sigma(c, trials, 1/float64(n)) {
				t.Errorf("n=%d: position %d: %d of %d", n, i, c, trials)
			}
		}
	}
}

func BenchmarkRandomElement(b *testing.B) {
	for _, size := range benchSizes {
		items := make([]int, size.n)
		rng := rand.New(rand.NewPCG(1, 2))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				v, _ := RandomElement(rng, slices.Values(items))
				sink += v
			}
			benchSink = sink
		})
	}
}
