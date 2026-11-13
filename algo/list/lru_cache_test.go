package list

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// lruModel is the oracle: keys in recency order (index 0 = most recent) and a value map.
type lruModel struct {
	capacity int
	order    []int
	vals     map[int]string
}

func (m *lruModel) touch(k int) {
	i := slices.Index(m.order, k)
	m.order = slices.Insert(slices.Delete(m.order, i, i+1), 0, k)
}

func (m *lruModel) get(k int) (string, bool) {
	v, ok := m.vals[k]
	if ok {
		m.touch(k)
	}
	return v, ok
}

func (m *lruModel) put(k int, v string) (Entry[int, string], bool) {
	if _, ok := m.vals[k]; ok {
		m.vals[k] = v
		m.touch(k)
		return Entry[int, string]{}, false
	}
	var ev Entry[int, string]
	did := false
	if len(m.order) == m.capacity {
		last := m.order[len(m.order)-1]
		ev, did = Entry[int, string]{last, m.vals[last]}, true
		delete(m.vals, last)
		m.order = m.order[:len(m.order)-1]
	}
	m.vals[k] = v
	m.order = slices.Insert(m.order, 0, k)
	return ev, did
}

// Worked example (hand-derived from Mattson et al.'s LRU stack order, catalog/list.md ALG-LIST-015);
// also the problem-site pattern, in our words.
func TestLRUCache_workedExample(t *testing.T) {
	c, err := NewLRUCache[int, string](2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put(1, "a")
	c.Put(2, "b")
	if v, ok := c.Get(1); !ok || v != "a" {
		t.Fatalf("Get(1) = (%q, %v)", v, ok)
	}
	if ev, did, err := c.Put(3, "c"); err != nil || !did || ev != (Entry[int, string]{2, "b"}) {
		t.Fatalf("Put(3) evicted (%v, %v, %v), want key 2", ev, did, err)
	}
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) hit after eviction")
	}
	if ev, did, _ := c.Put(4, "d"); !did || ev.Key != 1 {
		t.Fatalf("Put(4) evicted %v, want key 1", ev)
	}
	if v, ok := c.Get(3); !ok || v != "c" || c.Len() != 2 || c.Cap() != 2 {
		t.Fatalf("Get(3) = (%q, %v), Len %d", v, ok, c.Len())
	}
}

func TestLRUCache_peekDoesNotTouch(t *testing.T) {
	c, _ := NewLRUCache[string, int](2)
	c.Put("x", 1)
	c.Put("y", 2)
	if v, ok := c.Peek("x"); !ok || v != 1 {
		t.Fatal("Peek(x) missed")
	}
	if ev, _, _ := c.Put("z", 3); ev.Key != "x" {
		t.Errorf("evicted %q, want x (Peek must not count as a use)", ev.Key)
	}
	if _, ok := c.Peek("x"); ok {
		t.Error("Peek(x) hit after eviction")
	}
}

func TestLRUCache_putExistingUpdates(t *testing.T) {
	c, _ := NewLRUCache[int, int](2)
	c.Put(1, 10)
	c.Put(2, 20)
	if _, did, _ := c.Put(1, 11); did {
		t.Error("update evicted")
	}
	if ev, _, _ := c.Put(3, 30); ev.Key != 2 {
		t.Errorf("evicted %v, want 2", ev)
	}
	if v, _ := c.Get(1); v != 11 {
		t.Errorf("Get(1) = %d, want 11", v)
	}
}

func TestLRUCache_delete(t *testing.T) {
	c, _ := NewLRUCache[int, int](3)
	c.Put(1, 1)
	c.Put(2, 2)
	if !c.Delete(1) || c.Delete(1) || c.Delete(9) || c.Len() != 1 {
		t.Error("Delete reports wrong")
	}
	c.Put(3, 3)
	c.Put(4, 4)
	if _, did, _ := c.Put(4, 5); did || c.Len() != 3 {
		t.Errorf("Len = %d", c.Len())
	}
}

func TestLRUCache_capacityOne(t *testing.T) {
	c, _ := NewLRUCache[int, int](1)
	c.Put(1, 1)
	if ev, did, _ := c.Put(2, 2); !did || ev.Key != 1 || c.Len() != 1 {
		t.Errorf("(%v, %v), Len %d", ev, did, c.Len())
	}
}

func TestLRUCache_errors(t *testing.T) {
	for _, n := range []int{0, -1} {
		if c, err := NewLRUCache[int, int](n); !errors.Is(err, ErrCapacity) || c != nil {
			t.Errorf("capacity %d: err = %v", n, err)
		}
	}
	c, _ := NewLRUCache[float64, int](2)
	if _, _, err := c.Put(math.NaN(), 1); !errors.Is(err, ErrKeyNotComparable) || c.Len() != 0 {
		t.Errorf("NaN key: err = %v, Len %d", err, c.Len())
	}
	type pt struct{ x, y float64 }
	cs, _ := NewLRUCache[pt, int](2)
	if _, _, err := cs.Put(pt{1, math.NaN()}, 1); !errors.Is(err, ErrKeyNotComparable) {
		t.Errorf("struct NaN key: err = %v", err)
	}
	if _, _, err := c.Put(math.Inf(1), 1); err != nil {
		t.Errorf("Inf key: err = %v", err)
	}
}

func TestLRUCache_zeroAndNil(t *testing.T) {
	var z LRUCache[int, int]
	var n *LRUCache[int, int]
	for _, c := range []*LRUCache[int, int]{&z, n} {
		if c.Len() != 0 || c.Cap() != 0 || c.Delete(1) {
			t.Error("not empty")
		}
		if _, ok := c.Get(1); ok {
			t.Error("Get hit")
		}
		if _, ok := c.Peek(1); ok {
			t.Error("Peek hit")
		}
		if ev, did, err := c.Put(1, 7); err != nil || !did || ev != (Entry[int, int]{1, 7}) || c.Len() != 0 {
			t.Errorf("Put = (%v, %v, %v)", ev, did, err)
		}
	}
}

func TestLRUCache_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		capacity := 1 + rng.IntN(6)
		c, _ := NewLRUCache[int, string](capacity)
		m := &lruModel{capacity: capacity, vals: map[int]string{}}
		for op := range 80 {
			k := rng.IntN(10)
			switch rng.IntN(4) {
			case 0, 1:
				v := string(rune('a' + op%26))
				ev, did, err := c.Put(k, v)
				wev, wdid := m.put(k, v)
				if err != nil || did != wdid || ev != wev {
					t.Fatalf("trial %d op %d Put(%d): (%v, %v, %v), want (%v, %v)", trial, op, k, ev, did, err, wev, wdid)
				}
			case 2:
				v, ok := c.Get(k)
				wv, wok := m.get(k)
				if v != wv || ok != wok {
					t.Fatalf("trial %d op %d Get(%d): (%q, %v), want (%q, %v)", trial, op, k, v, ok, wv, wok)
				}
			default:
				_, had := m.vals[k]
				if had {
					i := slices.Index(m.order, k)
					m.order = slices.Delete(m.order, i, i+1)
					delete(m.vals, k)
				}
				if c.Delete(k) != had {
					t.Fatalf("trial %d op %d Delete(%d) != %v", trial, op, k, had)
				}
			}
			if c.Len() != len(m.order) {
				t.Fatalf("Len = %d, want %d", c.Len(), len(m.order))
			}
		}
	}
}

func TestLRUCache_noAllocSteadyState(t *testing.T) {
	c, _ := NewLRUCache[int, int](64)
	for i := range 64 {
		c.Put(i, i)
	}
	k := 64
	if a := testing.AllocsPerRun(1000, func() { c.Put(k, k); c.Get(k - 10); k++ }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkLRUCache(b *testing.B) {
	for _, s := range benchSizes {
		c, _ := NewLRUCache[int, int](s.n)
		rng := rand.New(rand.NewPCG(1, 2))
		keys := randInts(rng, 1<<16, 2*s.n) // about half the keys miss
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				k := keys[i&(len(keys)-1)]
				if _, ok := c.Get(k); !ok {
					c.Put(k, k)
				}
				i++
			}
		})
	}
}
