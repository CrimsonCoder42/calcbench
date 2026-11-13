package list

import "testing"

func newLRU(t *testing.T, capacity int) *LRUCache[string, int] {
	t.Helper()
	c, err := NewLRUCache[string, int](capacity)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func wantOldest(t *testing.T, c *LRUCache[string, int], key string, value int) {
	t.Helper()
	if k, v, ok := c.Oldest(); !ok || k != key || v != value {
		t.Fatalf("Oldest() = (%q, %d, %v), want (%q, %d, true)", k, v, ok, key, value)
	}
}

func wantNoOldest(t *testing.T, c *LRUCache[string, int]) {
	t.Helper()
	if k, v, ok := c.Oldest(); ok || k != "" || v != 0 {
		t.Fatalf("Oldest() = (%q, %d, %v), want (\"\", 0, false)", k, v, ok)
	}
}

// Worked example (BRIEF B2, from the LRU stack order of Mattson et al.): Put a, b, c; Get a → Oldest = b,
// and asking again does not change it.
func TestLRUCacheOldest_workedExample(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Get("a")
	wantOldest(t, c, "b", 2)
	wantOldest(t, c, "b", 2)
	if c.Len() != 3 {
		t.Fatalf("Len = %d after Oldest, want 3", c.Len())
	}
}

func TestLRUCacheOldest_emptyZeroNil(t *testing.T) {
	wantNoOldest(t, newLRU(t, 2))
	var z LRUCache[string, int]
	wantNoOldest(t, &z)
	var n *LRUCache[string, int]
	wantNoOldest(t, n)
}

func TestLRUCacheOldest_afterGet(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	wantOldest(t, c, "a", 1)
	c.Get("a")
	wantOldest(t, c, "b", 2)
}

func TestLRUCacheOldest_afterPutExisting(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Put("a", 10)
	wantOldest(t, c, "b", 2)
}

func TestLRUCacheOldest_afterPeek(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Peek("a")
	wantOldest(t, c, "a", 1)
}

func TestLRUCacheOldest_afterDeleteOldest(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Delete("a")
	wantOldest(t, c, "b", 2)
	c.Delete("b")
	c.Delete("c")
	wantNoOldest(t, c)
}

// Oldest names exactly the entry the next Put of a new key into a full cache evicts, and does not change it.
func TestLRUCacheOldest_matchesEviction(t *testing.T) {
	c := newLRU(t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Get("a")
	c.Peek("b")
	for i, k := range []string{"d", "e", "f", "g"} {
		wk, ov, ok := c.Oldest()
		if !ok {
			t.Fatalf("step %d: Oldest ok = false on a full cache", i)
		}
		ev, did, err := c.Put(k, 100+i)
		if err != nil || !did || ev.Key != wk || ev.Value != ov {
			t.Fatalf("step %d: Put(%q) evicted (%v, %v, %v), Oldest said (%q, %d)", i, k, ev, did, err, wk, ov)
		}
	}
}

func TestLRUCacheOldest_capacityOne(t *testing.T) {
	c := newLRU(t, 1)
	wantNoOldest(t, c)
	c.Put("a", 1)
	wantOldest(t, c, "a", 1)
	c.Put("b", 2)
	wantOldest(t, c, "b", 2)
	c.Get("b")
	wantOldest(t, c, "b", 2)
}

func TestLRUCacheOldest_noAlloc(t *testing.T) {
	c := newLRU(t, 4)
	c.Put("a", 1)
	c.Put("b", 2)
	if a := testing.AllocsPerRun(1000, func() { c.Oldest() }); a != 0 {
		t.Fatalf("Oldest allocates %v times per call, want 0", a)
	}
}

func BenchmarkLRUCacheOldest(b *testing.B) {
	c, err := NewLRUCache[int, int](1024)
	if err != nil {
		b.Fatal(err)
	}
	for i := range 1024 {
		c.Put(i, i)
	}
	b.ReportAllocs()
	var sink int
	for b.Loop() {
		k, _, _ := c.Oldest()
		sink += k
	}
	_ = sink
}
