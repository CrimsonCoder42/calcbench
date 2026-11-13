package list

// Oldest returns the least recently used entry, the one the next Put of a new key into a full cache would
// evict. The entry is not removed and its recency is unchanged: Oldest is not a use (unlike Get). ok is
// false on an empty cache, the zero value or a nil *LRUCache, with key and value their zero values.
//
// ID:      ALG-LIST-018
// Time:    O(1)   one read of the recency list's back element
// Space:   O(1), no allocation
// Source:  Mattson, Gecsei, Slutz, Traiger, "Evaluation techniques for storage hierarchies", IBM Systems Journal, 1970, 9(2):78–117, doi:10.1147/sj.92.0078 (LRU replacement as a stack algorithm)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (doubly linked list, O(1) delete) and ch. 11 (hash tables)
// Tags:    cache, lru, eviction, oldest, linked-list, generic
// Example: Put a, b, c; Get a → Oldest() = (b, its value, true); Oldest() again = b (unchanged)
func (c *LRUCache[K, V]) Oldest() (key K, value V, ok bool) {
	if c == nil {
		return key, value, false
	}
	e := c.order.Back()
	if e == nil {
		return key, value, false
	}
	return e.Value.Key, e.Value.Value, true
}
