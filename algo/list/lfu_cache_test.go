package list

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// lfuModel is the oracle: per key its value, use count and the logical time of its last use; the victim
// is the smallest (count, last use). O(n) per eviction, obviously right.
type lfuModel struct {
	capacity int
	now      int
	vals     map[int]int
	count    map[int]uint64
	last     map[int]int
}

func newLFUModel(capacity int) *lfuModel {
	return &lfuModel{capacity: capacity, vals: map[int]int{}, count: map[int]uint64{}, last: map[int]int{}}
}

func (m *lfuModel) use(k int) { m.now++; m.count[k]++; m.last[k] = m.now }

func (m *lfuModel) get(k int) (int, bool) {
	v, ok := m.vals[k]
	if ok {
		m.use(k)
	}
	return v, ok
}

func (m *lfuModel) put(k, v int) (Entry[int, int], bool) {
	if _, ok := m.vals[k]; ok {
		m.vals[k] = v
		m.use(k)
		return Entry[int, int]{}, false
	}
	var ev Entry[int, int]
	did := false
	if len(m.vals) == m.capacity {
		victim, first := 0, true
		for key := range m.vals {
			if first || m.count[key] < m.count[victim] || (m.count[key] == m.count[victim] && m.last[key] < m.last[victim]) {
				victim, first = key, false
			}
		}
		ev, did = Entry[int, int]{victim, m.vals[victim]}, true
		m.del(victim)
	}
	m.vals[k] = v
	m.use(k)
	return ev, did
}

func (m *lfuModel) del(k int) bool {
	_, ok := m.vals[k]
	delete(m.vals, k)
	delete(m.count, k)
	delete(m.last, k)
	return ok
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-016); also the problem-site pattern, in our words.
func TestLFUCache_workedExample(t *testing.T) {
	c, err := NewLFUCache[int, string](2)
	if err != nil {
		t.Fatal(err)
	}
	c.Put(1, "a")
	c.Put(2, "b")
	if v, ok := c.Get(1); !ok || v != "a" {
		t.Fatalf("Get(1) = (%q, %v)", v, ok)
	}
	if ev, did, _ := c.Put(3, "c"); !did || ev.Key != 2 {
		t.Fatalf("Put(3) evicted %v, want 2 (count 1 < 2)", ev)
	}
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) hit")
	}
	if v, ok := c.Get(3); !ok || v != "c" {
		t.Fatal("Get(3) missed")
	}
	if ev, did, _ := c.Put(4, "d"); !did || ev.Key != 1 {
		t.Fatalf("Put(4) evicted %v, want 1 (tie at 2, 1 used less recently)", ev)
	}
	if c.Len() != 2 || c.Cap() != 2 {
		t.Errorf("Len %d Cap %d", c.Len(), c.Cap())
	}
}

func TestLFUCache_putCountsAsUse(t *testing.T) {
	c, _ := NewLFUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 10) // count(1) = 2
	if ev, _, _ := c.Put(3, 3); ev.Key != 2 {
		t.Errorf("evicted %v, want 2", ev)
	}
	if v, _ := c.Get(1); v != 10 {
		t.Errorf("Get(1) = %d", v)
	}
}

func TestLFUCache_delete(t *testing.T) {
	c, _ := NewLFUCache[int, int](2)
	c.Put(1, 1)
	c.Get(1)
	c.Put(2, 2)
	if !c.Delete(1) || c.Delete(1) || c.Delete(7) || c.Len() != 1 {
		t.Error("Delete reports wrong")
	}
	c.Put(1, 1) // back with count 1, more recent than 2
	if ev, _, _ := c.Put(3, 3); ev.Key != 2 {
		t.Errorf("evicted %v, want 2 (Delete forgets the count)", ev)
	}
}

func TestLFUCache_capacityOne(t *testing.T) {
	c, _ := NewLFUCache[int, int](1)
	c.Put(1, 1)
	c.Get(1)
	c.Get(1)
	if ev, did, _ := c.Put(2, 2); !did || ev.Key != 1 || c.Len() != 1 {
		t.Errorf("(%v, %v) Len %d", ev, did, c.Len())
	}
	if ev, _, _ := c.Put(3, 3); ev.Key != 2 {
		t.Errorf("evicted %v, want 2", ev)
	}
}

func TestLFUCache_saturatedCount(t *testing.T) {
	c, _ := NewLFUCache[int, int](2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1) // 1 alone in its bucket (count 2)
	c.items[1].Value.bucket.Value.count = math.MaxUint64
	c.Get(1) // must stay at the maximum, not wrap to 0
	if got := c.items[1].Value.bucket.Value.count; got != math.MaxUint64 {
		t.Fatalf("count = %d, want MaxUint64", got)
	}
	if ev, _, _ := c.Put(3, 3); ev.Key != 2 {
		t.Errorf("evicted %v, want 2", ev)
	}
}

func TestLFUCache_errors(t *testing.T) {
	for _, n := range []int{0, -3} {
		if c, err := NewLFUCache[int, int](n); !errors.Is(err, ErrCapacity) || c != nil {
			t.Errorf("capacity %d: err = %v", n, err)
		}
	}
	c, _ := NewLFUCache[float64, int](2)
	if _, _, err := c.Put(math.NaN(), 1); !errors.Is(err, ErrKeyNotComparable) || c.Len() != 0 {
		t.Errorf("NaN key: err = %v", err)
	}
}

func TestLFUCache_zeroAndNil(t *testing.T) {
	var z LFUCache[int, int]
	var n *LFUCache[int, int]
	for _, c := range []*LFUCache[int, int]{&z, n} {
		if c.Len() != 0 || c.Cap() != 0 || c.Delete(1) {
			t.Error("not empty")
		}
		if _, ok := c.Get(1); ok {
			t.Error("Get hit")
		}
		if ev, did, err := c.Put(1, 7); err != nil || !did || ev != (Entry[int, int]{1, 7}) || c.Len() != 0 {
			t.Errorf("Put = (%v, %v, %v)", ev, did, err)
		}
	}
}

// checkBuckets verifies the frequency list: ascending counts, no empty bucket, sizes sum to Len.
func checkBuckets(t *testing.T, c *LFUCache[int, int]) {
	t.Helper()
	total := 0
	var prev uint64
	for b := c.buckets.Front(); b != nil; b = b.Next() {
		if b.Value.items.Len() == 0 || b.Value.count <= prev {
			t.Fatalf("bucket count %d (prev %d) size %d", b.Value.count, prev, b.Value.items.Len())
		}
		for e := b.Value.items.Front(); e != nil; e = e.Next() {
			if e.Value.bucket != b {
				t.Fatal("entry points to the wrong bucket")
			}
		}
		prev = b.Value.count
		total += b.Value.items.Len()
	}
	if total != c.Len() || len(c.items) != c.Len() {
		t.Fatalf("bucket sizes %d, map %d, Len %d", total, len(c.items), c.Len())
	}
}

func TestLFUCache_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		capacity := 1 + rng.IntN(6)
		c, _ := NewLFUCache[int, int](capacity)
		m := newLFUModel(capacity)
		for op := range 100 {
			k := rng.IntN(10)
			switch rng.IntN(5) {
			case 0, 1:
				ev, did, err := c.Put(k, op)
				wev, wdid := m.put(k, op)
				if err != nil || did != wdid || ev != wev {
					t.Fatalf("trial %d op %d Put(%d): (%v, %v, %v), want (%v, %v)", trial, op, k, ev, did, err, wev, wdid)
				}
			case 2, 3:
				v, ok := c.Get(k)
				wv, wok := m.get(k)
				if v != wv || ok != wok {
					t.Fatalf("trial %d op %d Get(%d): (%d, %v), want (%d, %v)", trial, op, k, v, ok, wv, wok)
				}
			default:
				if c.Delete(k) != m.del(k) {
					t.Fatalf("trial %d op %d Delete(%d) disagrees", trial, op, k)
				}
			}
			checkBuckets(t, c)
		}
	}
}

func TestLFUCache_noAllocRepeatedGet(t *testing.T) {
	c, _ := NewLFUCache[int, int](8)
	c.Put(1, 1)
	// 1 is alone in its bucket, so every Get renumbers the bucket in place.
	if a := testing.AllocsPerRun(1000, func() { c.Get(1) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkLFUCache(b *testing.B) {
	for _, s := range benchSizes {
		c, _ := NewLFUCache[int, int](s.n)
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
