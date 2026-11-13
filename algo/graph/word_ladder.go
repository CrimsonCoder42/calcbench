package graph

// WordLadder returns a shortest word ladder from start to target: a list of words, start first and target
// last, where each word differs from the one before it in exactly one letter position, and every word
// after start is in dictionary. It returns nil when no ladder exists (target not in dictionary, a length
// mismatch, or no connection) and [start] when start == target. Among equally short ladders it returns
// the first one breadth-first search finds (deterministic for the same input).
// Letters are bytes: a change to one multi-byte UTF-8 character is not a one-letter change. Dictionary
// words of another length than start are ignored; duplicates count once. dictionary is read, never
// modified.
//
// ID:      ALG-GRAPH-021
// Time:    O(N·L²) expected   N = len(dictionary), L = len(start) (N·L hashed keys of L-1 bytes; each key's bucket is scanned once)
// Space:   O(N·L²) extra (the keys), O(path length) for the result
// Source:  Knuth, The Stanford GraphBase: A Platform for Combinatorial Computing, ACM Press / Addison-Wesley, 1993, GB_WORDS ("The output phase": words adjacent when equal after suppressing one letter, one hash table per position) and LADDERS (shortest word ladders)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS shortest paths)
// Tags:    graph, word-ladder, bfs, implicit-graph, strings, shortest-path
func WordLadder(start, target string, dictionary []string) []string {
	if start == target {
		return []string{start}
	}
	L := len(start)
	if len(target) != L {
		return nil
	}
	// words[0] is start; words[1:] are the distinct dictionary words of length L, in dictionary order.
	words := make([]string, 1, len(dictionary)+1)
	words[0] = start
	index := make(map[string]int, len(dictionary)+1)
	index[start] = 0
	targetID := -1
	for _, w := range dictionary {
		if len(w) != L {
			continue
		}
		if _, dup := index[w]; dup {
			continue
		}
		index[w] = len(words)
		if w == target {
			targetID = len(words)
		}
		words = append(words, w)
	}
	if targetID == -1 {
		return nil
	}
	// A bucket holds every word that reads `rest` once the letter at position pos is suppressed.
	type key struct {
		pos  int
		rest string
	}
	buckets := make(map[key][]int, len(words)*L)
	for id, w := range words {
		for pos := range L {
			k := key{pos, w[:pos] + w[pos+1:]}
			buckets[k] = append(buckets[k], id)
		}
	}
	parent := make([]int, len(words))
	for i := range parent {
		parent[i] = -1
	}
	parent[0] = 0
	queue := []int{0}
	for head := 0; head < len(queue) && parent[targetID] == -1; head++ {
		v := queue[head]
		w := words[v]
		for pos := range L {
			k := key{pos, w[:pos] + w[pos+1:]}
			for _, u := range buckets[k] {
				if parent[u] == -1 {
					parent[u] = v
					queue = append(queue, u)
				}
			}
			delete(buckets, k) // every word in it is now discovered: never scan it again
		}
	}
	if parent[targetID] == -1 {
		return nil
	}
	length := 1
	for v := targetID; v != 0; v = parent[v] {
		length++
	}
	ladder := make([]string, length)
	for i, v := length-1, targetID; i >= 0; i-- {
		ladder[i] = words[v]
		v = parent[v]
	}
	return ladder
}
