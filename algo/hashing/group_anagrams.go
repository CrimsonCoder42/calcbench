package hashing

import "slices"

// GroupAnagrams returns words grouped into anagram classes: two words share a group when they hold the
// same multiset of runes. Each word's signature is its runes in sorted order. Groups are in the order of
// their first word; words keep input order, duplicates included. Bytes that are not valid UTF-8 each read
// as U+FFFD (the Go conversion rule). Empty input gives an empty, non-nil slice. words is read, never
// modified.
//
// ID:      ALG-HASH-010
// Time:    O(Σ Lᵢ log Lᵢ) expected   Lᵢ = runes in words[i]
// Space:   O(Σ Lᵢ) extra for the signatures; O(n) for the result   n = len(words)
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (signatures, §2.8, pp. 18–20)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (grouping by key in a hash table)
// Tags:    hash, anagram, signature, group, string, deterministic
func GroupAnagrams(words []string) [][]string {
	groups := [][]string{}
	index := make(map[string]int)
	var runes []rune
	for _, w := range words {
		runes = runes[:0]
		for _, r := range w {
			runes = append(runes, r)
		}
		slices.Sort(runes)
		key := string(runes)
		g, ok := index[key]
		if !ok {
			g = len(groups)
			index[key] = g
			groups = append(groups, nil)
		}
		groups[g] = append(groups[g], w)
	}
	return groups
}
